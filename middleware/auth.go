package middleware

import (
	"errors"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Raihanki/go-notes/config"
	"github.com/Raihanki/go-notes/helpers"
	"github.com/Raihanki/go-notes/models"
	"github.com/Raihanki/go-notes/resources"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func Auth(ctx *fiber.Ctx) error {
	// check authorization header
	headerAuth := ctx.Get("Authorization")
	if headerAuth == "" {
		return helpers.Response(ctx, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	// check token
	token := strings.Split(headerAuth, " ")[1]
	if token == "" {
		return helpers.Response(ctx, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	// validate token
	validatedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.ENV.JWT_SECRET), nil
	})
	if err != nil {
		return helpers.Response(ctx, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	userId := validatedToken.Claims.(jwt.MapClaims)["id"]
	user := new(models.User)
	errUser := config.DB.Where("id = ?", userId).Take(&user).Error
	if errors.Is(errUser, gorm.ErrRecordNotFound) {
		return helpers.Response(ctx, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	if errUser != nil {
		log.Error("Error fetching user : " + errUser.Error())
		return helpers.Response(ctx, fiber.StatusInternalServerError, "Failed to fetch user", nil)
	}

	ctx.Locals("user", resources.UserResource{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	})
	return ctx.Next()
}
