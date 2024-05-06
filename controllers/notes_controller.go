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
)

type NoteController struct {
}

func NewNoteController() *NoteController {
	return &NoteController{}
}

func (controller *NoteController) Index(ctx *fiber.Ctx) error {
	var notes []models.Note
	errNote := config.DB.Preload("User").Preload("Topic").Order("created_at desc").Find(&notes).Error
	if errNote != nil {
		log.Error("Error fetching notes : " + errNote.Error())
		return helpers.Response(ctx, fiber.StatusInternalServerError, "Failed to fetch notes", nil)
	}

	var noteResources []resources.NoteResource
	for _, note := range notes {
		noteResources = append(noteResources, resources.NoteResource{
			ID: note.ID,
			User: resources.UserPublicResource{
				ID:   note.User.ID,
				Name: note.User.Name,
			},
			Title:   note.Title,
			Content: note.Content,
			Topic: resources.TopicResource{
				ID:   note.Topic.ID,
				Name: note.Topic.Name,
			},
			CreatedAt: note.CreatedAt,
		})
	}

	return helpers.Response(ctx, fiber.StatusOK, "Successfully get all notes", noteResources)
}

func (controller *NoteController) Show(ctx *fiber.Ctx) error {
	idNote := ctx.Params("id")

	note := models.Note{}
	errNote := config.DB.Preload("User").Preload("Topic").Where("id = ?", idNote).Take(&note).Error
	if errors.Is(errNote, gorm.ErrRecordNotFound) {
		return helpers.Response(ctx, fiber.StatusNotFound, "Note not found", nil)
	}
	if errNote != nil {
		log.Error("Error fetching note : " + errNote.Error())
		return helpers.Response(ctx, fiber.StatusInternalServerError, "Failed to fetch note", nil)
	}

	noteResource := resources.NoteResource{
		ID: note.ID,
		User: resources.UserPublicResource{
			ID:   note.User.ID,
			Name: note.User.Name,
		},
		Title:   note.Title,
		Content: note.Content,
		Topic: resources.TopicResource{
			ID:   note.Topic.ID,
			Name: note.Topic.Name,
		},
		CreatedAt: note.CreatedAt,
	}

	return helpers.Response(ctx, fiber.StatusOK, "Successfully get note", noteResource)
}

func (controller *NoteController) Store(ctx *fiber.Ctx) error {
	noteRequest := new(request.NoteRequest)
	err := ctx.BodyParser(noteRequest)
	if err != nil {
		return helpers.Response(ctx, fiber.StatusUnprocessableEntity, "Unprocessable Entity", nil)
	}

	// todo user from auth
	user := models.User{}
	errUser := config.DB.Where("id = ?", 6).First(&user).Error
	if errors.Is(errUser, gorm.ErrRecordNotFound) {
		return helpers.Response(ctx, fiber.StatusNotFound, "User not found", nil)
	}
	if errUser != nil {
		log.Error("Error fetching user : " + errUser.Error())
		return helpers.Response(ctx, fiber.StatusInternalServerError, "Failed to fetch user", nil)
	}

	topic := models.Topic{}
	errTopic := config.DB.Where("id = ?", noteRequest.TopicID).First(&topic).Error
	if errors.Is(errTopic, gorm.ErrRecordNotFound) {
		return helpers.Response(ctx, fiber.StatusNotFound, "Topic not found", nil)
	}
	if errTopic != nil {
		log.Error("Error fetching topic : " + errTopic.Error())
		return helpers.Response(ctx, fiber.StatusInternalServerError, "Failed to fetch topic", nil)
	}

	note := models.Note{
		UserID:  user.ID,
		Title:   noteRequest.Title,
		Content: noteRequest.Content,
		TopicID: topic.ID,
	}
	errNote := config.DB.Create(&note).Error
	if errNote != nil {
		log.Error("Error creating note : " + errNote.Error())
		return helpers.Response(ctx, fiber.StatusInternalServerError, "Failed to create note", nil)
	}

	noteResource := resources.NoteResource{
		ID: note.ID,
		User: resources.UserPublicResource{
			ID:   user.ID,
			Name: user.Name,
		},
		Title:   note.Title,
		Content: note.Content,
		Topic: resources.TopicResource{
			ID:   topic.ID,
			Name: topic.Name,
		},
		CreatedAt: note.CreatedAt,
	}

	return helpers.Response(ctx, fiber.StatusCreated, "Successfully create note", noteResource)
}

func (controller *NoteController) Update(ctx *fiber.Ctx) error {
	idNote := ctx.Params("id")

	noteReqeust := new(request.NoteRequest)
	err := ctx.BodyParser(noteReqeust)
	if err != nil {
		return helpers.Response(ctx, fiber.StatusUnprocessableEntity, "Unprocessable Entity", nil)
	}

	note := models.Note{}
	errNote := config.DB.Where("id = ?", idNote).Take(&note).Error
	if errors.Is(errNote, gorm.ErrRecordNotFound) {
		return helpers.Response(ctx, fiber.StatusNotFound, "Note not found", nil)
	}
	if errNote != nil {
		log.Error("Error fetching note : " + errNote.Error())
		return helpers.Response(ctx, fiber.StatusInternalServerError, "Failed to fetch note", nil)
	}

	//todo athorization
	user := models.User{}
	errUser := config.DB.Where("id = ?", 6).First(&user).Error
	if errors.Is(errUser, gorm.ErrRecordNotFound) {
		return helpers.Response(ctx, fiber.StatusNotFound, "User not found", nil)
	}
	if errUser != nil {
		log.Error("Error fetching user : " + errUser.Error())
		return helpers.Response(ctx, fiber.StatusInternalServerError, "Failed to fetch user", nil)
	}

	topic := models.Topic{}
	errTopic := config.DB.Where("id = ?", noteReqeust.TopicID).First(&topic).Error
	if errors.Is(errTopic, gorm.ErrRecordNotFound) {
		return helpers.Response(ctx, fiber.StatusNotFound, "Topic not found", nil)
	}
	if errTopic != nil {
		log.Error("Error fetching topic : " + errTopic.Error())
		return helpers.Response(ctx, fiber.StatusInternalServerError, "Failed to fetch topic", nil)
	}

	errUpdateNote := config.DB.Where("id = ?", note.ID).Updates(&models.Note{
		Title:   noteReqeust.Title,
		Content: noteReqeust.Content,
		TopicID: topic.ID,
	}).Error
	if errUpdateNote != nil {
		log.Error("Error updating note : " + errUpdateNote.Error())
		return helpers.Response(ctx, fiber.StatusInternalServerError, "Failed to update note", nil)
	}

	return helpers.Response(ctx, fiber.StatusOK, "Successfully update note", nil)
}

func (controller *NoteController) Destroy(ctx *fiber.Ctx) error {
	idNote := ctx.Params("id")

	note := models.Note{}
	result := config.DB.Where("id = ?", idNote).Delete(&note)
	if result.RowsAffected == 0 {
		return helpers.Response(ctx, fiber.StatusNotFound, "Note not found", nil)
	}
	if result.Error != nil {
		log.Error("Error deleting note : " + result.Error.Error())
		return helpers.Response(ctx, fiber.StatusInternalServerError, "Failed to delete note", nil)
	}

	return helpers.Response(ctx, fiber.StatusOK, "Successfully delete note", nil)
}
