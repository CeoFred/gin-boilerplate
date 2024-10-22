package repository

import (
	"errors"
	"fmt"

	"github.com/CeoFred/gin-boilerplate/internal/models"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	Create(user *models.User) error
	Find(id uuid.UUID) (*models.User, error)
	Exists(id uuid.UUID) (bool, error)
	Where(condition, value string) ([]*models.User, error)
	FindByCondition(condition, value string) (*models.User, bool, error)
	FindByAccountType(value string) ([]*models.User, bool, error)
	RawSmartSelect(q string, res interface{}, args ...interface{}) error
	Save(user *models.User) (*models.User, error)
}

type UserRepository struct {
	database *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &UserRepository{
		database: db,
	}
}

func (a *UserRepository) Find(id uuid.UUID) (*models.User, error) {
	var user *models.User
	err := a.database.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (a *UserRepository) Exists(id uuid.UUID) (bool, error) {
	var user *models.User
	err := a.database.First(&user, "id = ?", id).Error
	if err != nil {
		return false, err
	}

	if user == nil {
		return false, nil
	}
	return true, nil
}

func (a *UserRepository) Where(condition, value string) ([]*models.User, error) {
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

func (a *UserRepository) Create(user *models.User) error {
	return a.database.Model(&user).Create(user).Error
}

func (a *UserRepository) Save(user *models.User) (*models.User, error) {

	txn := a.database.Model(user).Where("id = ?", user.ID).Updates(&user).First(user)

	if txn.RowsAffected == 0 {
		return nil, errors.New("no record updated")
	}

	if txn.Error != nil {
		return nil, txn.Error
	}

	return user, nil
}

func (a *UserRepository) RawCount(q string, count *int64) error {
	return a.database.Model(&models.User{}).Raw(q).Count(count).Error
}

func (a *UserRepository) QueryWithArgs(q string, args ...interface{}) ([]*models.User, error) {
	var n []*models.User
	err := a.database.Raw(q, args...).Find(&n).Error

	if err != nil {
		return nil, err
	}

	if n != nil {
		return n, nil
	}

	return nil, nil
}

func (a *UserRepository) RawSmartSelect(q string, res interface{}, args ...interface{}) error {
	return a.database.Model(&models.User{}).Raw(q, args...).Find(res).Error
}
