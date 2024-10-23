package bootstrap

import (
	"github.com/CeoFred/gin-boilerplate/internal/repository"
	"github.com/CeoFred/gin-boilerplate/internal/service"
	"github.com/CeoFred/gin-boilerplate/internal/service/streaming"

	"gorm.io/gorm"
)

type AppDependencies struct {
	EmailService    service.EmailServicer
	UserRepo        repository.UserRepositoryInterface
	EventProducer   streaming.EventProducer
	DatabaseService *gorm.DB
}

func InitializeDependencies(db *gorm.DB) *AppDependencies {
	return &AppDependencies{
		UserRepo:        repository.NewUserRepository(db),
		EmailService:    service.NewEmailService(),
		DatabaseService: db,
	}
}
