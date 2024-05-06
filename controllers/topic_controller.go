package controllers

import (
	log "github.com/sirupsen/logrus"

	"github.com/Raihanki/go-notes/config"
	"github.com/Raihanki/go-notes/helpers"
	"github.com/Raihanki/go-notes/models"
	"github.com/Raihanki/go-notes/resources"
	"github.com/gofiber/fiber/v2"
)

type TopicController struct {
}

func NewTopicController() *TopicController {
	return &TopicController{}
}

func (controller *TopicController) Index(ctx *fiber.Ctx) error {
	var topics []models.Topic

	err := config.DB.Find(&topics).Error
	if err != nil {
		log.Error("Error fetching notes : " + err.Error())
		helpers.Response(ctx, fiber.StatusInternalServerError, "Failed to fetch notes", nil)
		return err
	}

	var topicResources []resources.TopicResource
	for _, topic := range topics {
		topicResources = append(topicResources, resources.TopicResource{
			ID:   topic.ID,
			Name: topic.Name,
		})
	}
	return helpers.Response(ctx, fiber.StatusOK, "Successfully get all topics", topicResources)
}
