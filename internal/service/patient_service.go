package service

import (
	"context"

	"github.com/jeagerism/agnos-backend-assignment/internal/model"
	"github.com/jeagerism/agnos-backend-assignment/internal/response"
	"github.com/jeagerism/agnos-backend-assignment/internal/repository"
	"github.com/jeagerism/agnos-backend-assignment/internal/request"
)

type PatientService interface {
	Search(ctx context.Context, req request.SearchPatientRequest) ([]response.PatientResponse, int64, error)
}

type patientService struct {
	repo repository.PatientRepository
}

func NewPatientService(repo repository.PatientRepository) PatientService {
	return &patientService{repo: repo}
}

func (s *patientService) Search(ctx context.Context, req request.SearchPatientRequest) ([]response.PatientResponse, int64, error) {
	filter := repository.PatientSearchFilter{
		Hospital:    req.Hospital,
		PassportID:  req.PassportID,
		NationalID:  req.NationalID,
		FirstName:   req.FirstName,
		MiddleName:  req.MiddleName,
		LastName:    req.LastName,
		DateOfBirth: req.DateOfBirth,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		Limit:       req.Limit,
		Page:        req.Page,
	}

	patients, total, err := s.repo.Search(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return toPatientResponses(patients), total, nil
}

func toPatientResponses(patients []model.Patient) []response.PatientResponse {
	result := make([]response.PatientResponse, 0, len(patients))
	for _, p := range patients {
		result = append(result, response.PatientResponse{
			ID:           p.ID,
			Hospital:     p.Hospital,
			PatientHN:    p.PatientHN,
			NationalID:   p.NationalID,
			PassportID:   p.PassportID,
			FirstNameTH:  p.FirstNameTH,
			MiddleNameTH: p.MiddleNameTH,
			LastNameTH:   p.LastNameTH,
			FirstNameEN:  p.FirstNameEN,
			MiddleNameEN: p.MiddleNameEN,
			LastNameEN:   p.LastNameEN,
			DateOfBirth:  p.DateOfBirth,
			PhoneNumber:  p.PhoneNumber,
			Email:        p.Email,
			Gender:       p.Gender,
		})
	}
	return result
}

