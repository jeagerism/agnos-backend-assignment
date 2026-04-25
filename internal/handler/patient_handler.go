package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jeagerism/agnos-backend-assignment/internal/service"
)

type PatientHandler struct {
	service service.PatientService
}

func NewPatientHandler(svc service.PatientService) *PatientHandler {
	return &PatientHandler{service: svc}
}

// TODO: implement patient handlers
var _ = (*gin.Context)(nil)
