package repository

import (
	"github.com/jeagerism/agnos-backend-assignment/internal/model"
	"gorm.io/gorm"
)

type PatientRepository interface {
	// TODO: define patient repository methods
}

type patientRepository struct {
	db *gorm.DB
}

func NewPatientRepository(db *gorm.DB) PatientRepository {
	return &patientRepository{db: db}
}

var _ = (*model.Patient)(nil)
