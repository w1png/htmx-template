package main

import (
	"log"

	"github.com/w1png/htmx-template/config"
	"github.com/w1png/htmx-template/models"
	"github.com/w1png/htmx-template/storage"
	"gorm.io/gorm"
)

func createDefaultAdmin() error {
	if err := storage.GormStorageInstance.DB.First(&models.User{}, 1).Error; err == nil {
		return nil
	} else if gorm.ErrRecordNotFound != err {
		return err
	}

	admin, err := models.NewUser("admin", "admin", true)
	if err != nil {
		return err
	}
	return storage.GormStorageInstance.DB.Create(&admin).Error
}

func main() {
	if err := config.InitConfig(); err != nil {
		log.Fatal(err)
	}

	if err := storage.InitStorage(); err != nil {
		log.Fatal(err)
	}

	if err := createDefaultAdmin(); err != nil {
		log.Fatal(err)
	}

	server := NewHTTPServer()

	log.Fatal(server.Run())
}
