package helpers

import (
	"github.com/CeoFred/gin-boilerplate/internal/models"
	"github.com/gofrs/uuid"
)

type UserRepositoryInterface interface {
	Where(condition, value string) ([]*models.User, error)
	FindByCondition(condition, value string) (*models.User, bool, error)
	FindByAccountType(value string) ([]*models.User, bool, error)
	Find(id uuid.UUID) (*models.User, error)
	Exists(id uuid.UUID) (bool, error)
	Create(item *models.User) error
	Save(u *models.User) (*models.User, error)
	RawCount(q string, count *int64) error
	QueryWithArgs(q string, args ...interface{}) ([]*models.User, error)
	RawSmartSelect(q string, res interface{}, args ...interface{}) error
}
