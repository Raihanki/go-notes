package controllers

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/Raihanki/go-notes/config"
	"github.com/Raihanki/go-notes/helpers"
	"github.com/Raihanki/go-notes/models"
	"github.com/Raihanki/go-notes/request"
	"github.com/Raihanki/go-notes/resources"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
}

func NewUserController() *UserController {
	return &UserController{}
}

func (controller *UserController) Register(ctx *fiber.Ctx) error {
	registerRequest := new(request.RegisterRequest)
	err := ctx.BodyParser(registerRequest)
	if err != nil {
		return helpers.Response(ctx, fiber.StatusUnprocessableEntity, "Unprocessable Entity", nil)
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

	userResource := resources.UserResource{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
	return helpers.Response(ctx, fiber.StatusCreated, "Successfully registered", userResource)
}

func (controller *UserController) Login(ctx *fiber.Ctx) error {
	loginRequest := new(request.LoginRequest)
	err := ctx.BodyParser(loginRequest)
	if err != nil {
		return helpers.Response(ctx, fiber.StatusUnprocessableEntity, "Unprocessable Entity", nil)
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

	response := map[string]interface{}{
		"token": "abc123xyz",
	}
	return helpers.Response(ctx, fiber.StatusOK, "Successfully login", response)
}
