package main

import (
	"github.com/jeagerism/agnos-backend-assignment/config"
	"github.com/jeagerism/agnos-backend-assignment/database"
	"github.com/jeagerism/agnos-backend-assignment/internal/handler"
	"github.com/jeagerism/agnos-backend-assignment/internal/repository"
	"github.com/jeagerism/agnos-backend-assignment/internal/service"
	"github.com/jeagerism/agnos-backend-assignment/internal/router"
)

func main() {
	cfg := config.LoadConfig()
	database.ConnectDB(cfg)

	staffRepo := repository.NewStaffRepository(database.DB)
	staffSvc := service.NewStaffService(staffRepo, cfg.JWTSecret)
	staffHandler := handler.NewStaffHandler(staffSvc)

	r := router.SetupRouter(staffHandler)
	r.Run(":" + cfg.AppPort)
}
