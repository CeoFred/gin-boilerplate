package repository

import (
	"errors"
	"fmt"

	"github.com/CeoFred/gin-boilerplate/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	database *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		database: db,
	}
}

func (a *UserRepository) FindRecordsByCondition(condition, value string) ([]*models.User, error) {
	var users []*models.User
	err := a.database.Raw(fmt.Sprintf(`SELECT * FROM users WHERE %s = ?`, condition), value).Scan(&users).Error
	if err != nil {
		return nil, err
	}
	if users != nil {
		return users, nil
	}
	return nil, nil
}

func (a *UserRepository) FindByCondition(condition, value string) (*models.User, bool, error) {
	var user *models.User
	err := a.database.Where(condition, value).Find(&user).Error
	if err != nil {
		return nil, false, err
	}
	if user.Email != "" {
		return user, true, nil
	}
	return nil, false, nil
}

func (a *UserRepository) FindUserByCondition(condition, value string) (*models.User, bool, error) {
	var user *models.User
	err := a.database.Raw(fmt.Sprintf(`SELECT * FROM users WHERE %s = ?`, condition), value).Scan(&user).Error
	if err != nil {
		return nil, false, err
	}
	if user != nil {
		return user, true, nil
	}
	return nil, false, nil
}

func (a *UserRepository) FindByAccountType(value string) ([]*models.User, bool, error) {
	var user []*models.User
	err := a.database.Raw(`SELECT * FROM users WHERE account_type = ?`, value).Scan(&user).Error
	if err != nil {
		return nil, false, err
	}
	if user != nil {
		return user, true, nil
	}
	return nil, false, nil
}

func (a *UserRepository) CreateUser(user *models.User) error {
	return a.database.Model(&models.User{}).Create(user).Error
}

func (a *UserRepository) UpdateUserByCondition(condition, value string, update *models.User) (*models.User, error) {
	user := &models.User{}
	rows := a.database.Model(user).Where(fmt.Sprintf(`%s = ?`, condition), value).Updates(&update).First(user)
	if rows.RowsAffected == 0 {
		return nil, errors.New("no record updated")
	}
	return user, nil
}
