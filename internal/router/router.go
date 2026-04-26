package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jeagerism/agnos-backend-assignment/internal/handler"
	"github.com/jeagerism/agnos-backend-assignment/internal/middleware"
)

func SetupRouter(staffHandler *handler.StaffHandler, patientHandler *handler.PatientHandler, jwtSecret string) *gin.Engine {
	r := gin.Default()

	staffGroup := r.Group("/staff")
	{
		staffGroup.POST("/create", staffHandler.Create)
		staffGroup.POST("/login", staffHandler.Login)
	}

	patientGroup := r.Group("/patient")
	patientGroup.Use(middleware.Auth(jwtSecret))
	patientGroup.GET("/search", patientHandler.Search)

	return r
}
