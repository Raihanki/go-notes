package config

import (
	"fmt"

	"github.com/Raihanki/go-notes/models"
	log "github.com/sirupsen/logrus"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=true&loc=Local",
		ENV.DB_USERNAME,
		ENV.DB_PASSWORD,
		ENV.DB_HOST,
		ENV.DB_PORT,
		ENV.DB_NAME,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Err DB : " + err.Error())
	}

	// migrate database
	db.AutoMigrate(&models.Topic{}, &models.User{}, &models.Note{})

	DB = db
	log.Info("Database successfully connected")
}
