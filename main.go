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

	patientRepo := repository.NewPatientRepository(database.DB)
	patientSvc := service.NewPatientService(patientRepo)
	patientHandler := handler.NewPatientHandler(patientSvc)

	r := router.SetupRouter(staffHandler, patientHandler, cfg.JWTSecret)

	r.Run(":" + cfg.AppPort)	
}
