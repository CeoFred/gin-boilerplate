package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/CeoFred/gin-boilerplate/internal/bootstrap"
)

func Routes(r *gin.RouterGroup, d *bootstrap.AppDependencies) {

	RegisterUserRoutes(r, d)
	RegisterAuthRoutes(r, d)

}
