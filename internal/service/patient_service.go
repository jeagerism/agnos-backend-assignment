package service

import "github.com/jeagerism/agnos-backend-assignment/internal/repository"

type PatientService interface {
	// TODO: define patient service methods
}

type patientService struct {
	repo repository.PatientRepository
}

func NewPatientService(repo repository.PatientRepository) PatientService {
	return &patientService{repo: repo}
}
