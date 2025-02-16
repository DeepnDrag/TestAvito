package storage

import (
	"TestAvito/internal/models"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (s *UserRepo) CreateUser(username, password string) (*models.User, error) {
	user := models.User{Username: username, Password: password}
	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserRepo) GetUserByUsername(username string) (*models.User, error) {
	var user models.User

	err := s.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserRepo) UpdateUser(updatedUser *models.User) (*models.User, error) {
	var user models.User

	err := s.db.First(&user, updatedUser.ID).Error
	if err != nil {
		return nil, err
	}

	err = s.db.Model(&user).Updates(updatedUser).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserRepo) UpdateTwoUsers(updatedUser1 *models.User, updatedUser2 *models.User) (*models.User, *models.User, error) {
	tx := s.db.Begin()

	var user1 models.User
	err := tx.Model(&user1).Where("id = ?", updatedUser1.ID).Updates(updatedUser1).Error
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	var user2 models.User
	err = tx.Model(&user2).Where("id = ?", updatedUser2.ID).Updates(updatedUser2).Error
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, nil, err
	}

	return &user1, &user2, nil
}
