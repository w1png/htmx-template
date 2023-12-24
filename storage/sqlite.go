package storage

import (
	"fmt"

	"github.com/w1png/htmx-template/config"
	"github.com/w1png/htmx-template/errors"
	"github.com/w1png/htmx-template/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SqliteStorage struct {
	DB *gorm.DB
}

func NewSQLiteStorage() (*SqliteStorage, error) {
	storage := &SqliteStorage{}

	var err error
	if storage.DB, err = gorm.Open(sqlite.Open(config.ConfigInstance.SqlitePath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}); err != nil {
		return nil, errors.NewDatabaseConnectionError(err.Error())
	}

	if err := storage.autoMigrate(); err != nil {
		return nil, errors.NewDatabaseMigrationError(err.Error())
	}

	return storage, nil
}

func (s *SqliteStorage) autoMigrate() error {
	return s.DB.AutoMigrate(&models.User{})
}

func (s *SqliteStorage) CreateUser(user *models.User) error {
	return s.DB.Create(user).Error
}

func (s *SqliteStorage) GetUserById(id uint) (*models.User, error) {
	var user models.User
	if err := s.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewObjectNotFoundError(fmt.Sprintf("User with id: %d", id))
		}
		return nil, err
	}

	return &user, nil
}

func (s *SqliteStorage) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := s.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewObjectNotFoundError(fmt.Sprintf("User with username: %s", username))
		}
		return nil, err
	}

	return &user, nil
}

func (s *SqliteStorage) GetUsers() ([]*models.User, error) {
	var users []*models.User
	if err := s.DB.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (s *SqliteStorage) UpdateUser(user *models.User) error {
	if _, err := s.GetUserById(user.ID); err != nil {
		return err
	}

	return s.DB.Save(user).Error
}

func (s *SqliteStorage) DeleteUserById(id uint) error {
	if id == 1 {
		return errors.NewMainAdminDeletionError("Cannot delete main admin")
	}

	if _, err := s.GetUserById(id); err != nil {
		return err
	}

	return s.DB.Delete(&models.User{}, id).Error
}
