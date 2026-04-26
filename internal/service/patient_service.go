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
	GetPatientByID(ctx context.Context, hn string) (response.PatientResponse, error)
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

func (s *patientService) GetPatientByID(ctx context.Context, id string) (response response.PatientResponse, err error) {
	patient, err := s.repo.GetPatientByID(ctx, id)
	if err != nil {
		return response, err
	}
	return toPatientResponse(patient), nil
}

func toPatientResponse(patient model.Patient) response.PatientResponse {
	return response.PatientResponse{
		ID:           patient.ID,
		Hospital:     patient.Hospital,
		PatientHN:    patient.PatientHN,
		NationalID:   patient.NationalID,
		PassportID:   patient.PassportID,
		FirstNameTH:  patient.FirstNameTH,
		MiddleNameTH: patient.MiddleNameTH,
		LastNameTH:   patient.LastNameTH,
		FirstNameEN:  patient.FirstNameEN,
		MiddleNameEN: patient.MiddleNameEN,
		LastNameEN:   patient.LastNameEN,
		DateOfBirth:  patient.DateOfBirth,
		PhoneNumber:  patient.PhoneNumber,
		Email:        patient.Email,
		Gender:       patient.Gender,
	}
}