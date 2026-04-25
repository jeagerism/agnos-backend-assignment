package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jeagerism/agnos-backend-assignment/internal/handler"
)

func SetupRouter(staffHandler *handler.StaffHandler) *gin.Engine {
	r := gin.Default()

	staffGroup := r.Group("/staff")
	{
		staffGroup.POST("/create", staffHandler.Create)
		staffGroup.POST("/login", staffHandler.Login)
	}

	return r
}
