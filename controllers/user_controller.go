package controllers

import (
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/Raihanki/go-notes/config"
	"github.com/Raihanki/go-notes/helpers"
	"github.com/Raihanki/go-notes/models"
	"github.com/Raihanki/go-notes/request"
	"github.com/Raihanki/go-notes/resources"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
}

func NewUserController() *UserController {
	return &UserController{}
}

func (controller *UserController) Show(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(resources.UserResource)
	return helpers.Response(ctx, fiber.StatusOK, "Successfully get user", user)
}

func (controller *UserController) Register(ctx *fiber.Ctx) error {
	registerRequest := new(request.RegisterRequest)

	err := ctx.BodyParser(registerRequest)
	if err != nil {
		return helpers.Response(ctx, fiber.StatusUnprocessableEntity, "Unprocessable Entity", nil)
	}

	errValidate := validator.New().Struct(registerRequest)
	if errValidate != nil {
		var errorResponse []helpers.ErrorResponse
		for _, err := range errValidate.(validator.ValidationErrors) {
			errorResponse = append(errorResponse, helpers.ErrorResponse{
				FailedField: err.Field(),
				Tag:         err.Tag(),
				Message:     err.Error(),
			})
		}
		return helpers.Response(ctx, fiber.StatusUnprocessableEntity, "Unprocessable Entity", errorResponse)
	}

	if registerRequest.Password != registerRequest.PasswordConfirmation {
		return helpers.Response(ctx, fiber.StatusUnprocessableEntity, "Password not match", nil)
	}

	//hash password
	hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if errHash != nil {
		log.Error("Error hashing password : " + errHash.Error())
		return helpers.Response(ctx, fiber.StatusInternalServerError, "Failed to register user", nil)
	}

	user := models.User{
		Name:     registerRequest.Name,
		Email:    registerRequest.Email,
		Password: string(hashedPassword),
	}
	errCreateUser := config.DB.Create(&user).Error
	if errCreateUser != nil {
		log.Error("Error creating user : " + errCreateUser.Error())
		return helpers.Response(ctx, fiber.StatusInternalServerError, "Failed to register user", nil)
	}

	//generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})
	signedToken, errToken := token.SignedString([]byte(config.ENV.JWT_SECRET))
	if errToken != nil {
		log.Error("Error generating token : " + errToken.Error())
		return helpers.Response(ctx, fiber.StatusInternalServerError, "Failed to login", nil)
	}

	return helpers.Response(ctx, fiber.StatusCreated, "Successfully registered", map[string]interface{}{
		"token": signedToken,
	})
}

func (controller *UserController) Login(ctx *fiber.Ctx) error {
	loginRequest := new(request.LoginRequest)

	err := ctx.BodyParser(loginRequest)
	if err != nil {
		return helpers.Response(ctx, fiber.StatusUnprocessableEntity, "Unprocessable Entity", nil)
	}

	errValidate := validator.New().Struct(loginRequest)
	if errValidate != nil {
		var errorResponse []helpers.ErrorResponse
		for _, err := range errValidate.(validator.ValidationErrors) {
			errorResponse = append(errorResponse, helpers.ErrorResponse{
				FailedField: err.Field(),
				Tag:         err.Tag(),
				Message:     err.Error(),
			})
		}
		return helpers.Response(ctx, fiber.StatusUnprocessableEntity, "Unprocessable Entity", errorResponse)
	}

	user := models.User{}
	errUser := config.DB.Where("email = ?", loginRequest.Email).Take(&user).Error
	if errors.Is(errUser, gorm.ErrRecordNotFound) {
		return helpers.Response(ctx, fiber.StatusNotFound, "Invalid email or password", nil)
	}
	if errUser != nil {
		log.Error("Error fetching user : " + errUser.Error())
		return helpers.Response(ctx, fiber.StatusInternalServerError, "Failed to login", nil)
	}

	//verify password
	errPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if errPassword != nil {
		return helpers.Response(ctx, fiber.StatusUnauthorized, "Invalid email or password", nil)
	}

	//generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})
	signedToken, errToken := token.SignedString([]byte(config.ENV.JWT_SECRET))
	if errToken != nil {
		log.Error("Error generating token : " + errToken.Error())
		return helpers.Response(ctx, fiber.StatusInternalServerError, "Failed to login", nil)
	}

	response := map[string]interface{}{
		"token": signedToken,
	}
	return helpers.Response(ctx, fiber.StatusOK, "Successfully login", response)
}
