package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jeagerism/agnos-backend-assignment/internal/model"
	"github.com/jeagerism/agnos-backend-assignment/internal/repository"
	"github.com/jeagerism/agnos-backend-assignment/internal/request"
	"github.com/jeagerism/agnos-backend-assignment/internal/response"
	"github.com/jeagerism/agnos-backend-assignment/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockPatientRepo struct {
	searchFn      func(ctx context.Context, filter repository.PatientSearchFilter) ([]model.Patient, int64, error)
	getByIDFn     func(ctx context.Context, id string) (model.Patient, error)
}

func (m *mockPatientRepo) Search(ctx context.Context, filter repository.PatientSearchFilter) ([]model.Patient, int64, error) {
	return m.searchFn(ctx, filter)
}

func (m *mockPatientRepo) GetPatientByID(ctx context.Context, id string) (model.Patient, error) {
	if m.getByIDFn == nil {
		return model.Patient{}, nil
	}
	return m.getByIDFn(ctx, id)
}

func TestPatientServiceSearch(t *testing.T) {
	nationalID := "1111111111111"
	passportID := "P12345678"
	firstName := "John"
	middleName := "Michael"
	lastName := "Dylan"
	dateOfBirth := "1990-05-15"
	phone := "0812345678"
	email := "john@example.com"
	patientID := uuid.MustParse("26fd07af-0d69-4ff6-a594-3e2f7e56fc2c")
	patientDOB := time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		req           request.SearchPatientRequest
		repoResult    []model.Patient
		repoTotal     int64
		repoErr       error
		assertFilter  func(t *testing.T, f repository.PatientSearchFilter)
		assertResult  func(t *testing.T, gotErr error, got []response.PatientResponse, total int64)
	}{
		{
			name: "success maps request to filter and model to response",
			req: request.SearchPatientRequest{
				Hospital:    "Hospital A",
				NationalID:  &nationalID,
				FirstName:   firstName,
				DateOfBirth: &dateOfBirth,
				PhoneNumber: &phone,
				Email:       &email,
				Page:        10,
				Limit:       20,
			},
			repoResult: []model.Patient{
				{
					ID:          patientID,
					Hospital:    "Hospital A",
					PatientHN:   "HN001",
					NationalID:  &nationalID,
					FirstNameTH: "จอห์น",
					LastNameTH:  "ดีแลน",
					DateOfBirth: patientDOB,
					PhoneNumber: &phone,
					Email:       &email,
					Gender:      "M",
				},
			},
			repoTotal: 42,
			assertFilter: func(t *testing.T, f repository.PatientSearchFilter) {
				assert.Equal(t, "Hospital A", f.Hospital)
				require.NotNil(t, f.NationalID)
				assert.Equal(t, nationalID, *f.NationalID)
				assert.Equal(t, firstName, f.FirstName)
				require.NotNil(t, f.DateOfBirth)
				assert.Equal(t, dateOfBirth, *f.DateOfBirth)
				require.NotNil(t, f.PhoneNumber)
				assert.Equal(t, phone, *f.PhoneNumber)
				require.NotNil(t, f.Email)
				assert.Equal(t, email, *f.Email)
				assert.Equal(t, 10, f.Page)
				assert.Equal(t, 20, f.Limit)
			},
			assertResult: func(t *testing.T, gotErr error, got []response.PatientResponse, total int64) {
				assert.NoError(t, gotErr)
				assert.Equal(t, int64(42), total)
				require.Len(t, got, 1)
				assert.Equal(t, patientID, got[0].ID)
				assert.Equal(t, "HN001", got[0].PatientHN)
				assert.Equal(t, "Hospital A", got[0].Hospital)
				assert.Equal(t, "M", got[0].Gender)
				require.NotNil(t, got[0].NationalID)
				assert.Equal(t, nationalID, *got[0].NationalID)
			},
		},
		{
			name: "returns repository error as-is",
			req: request.SearchPatientRequest{
				Hospital: "Hospital A",
			},
			repoErr: errors.New("db timeout"),
			assertFilter: func(t *testing.T, f repository.PatientSearchFilter) {
				assert.Equal(t, "Hospital A", f.Hospital)
			},
			assertResult: func(t *testing.T, gotErr error, got []response.PatientResponse, total int64) {
				assert.Error(t, gotErr)
				assert.EqualError(t, gotErr, "db timeout")
				assert.Equal(t, int64(0), total)
				assert.Nil(t, got)
			},
		},
		{
			name: "maps multiple query params including passport and names",
			req: request.SearchPatientRequest{
				Hospital:   "Hospital B",
				PassportID: &passportID,
				FirstName:  "jo",
				MiddleName: &middleName,
				LastName:   &lastName,
				Page:      5,
				Limit:      15,
			},
			assertFilter: func(t *testing.T, f repository.PatientSearchFilter) {
				assert.Equal(t, "Hospital B", f.Hospital)
				require.NotNil(t, f.PassportID)
				assert.Equal(t, passportID, *f.PassportID)
				assert.Equal(t, "jo", f.FirstName)
				require.NotNil(t, f.MiddleName)
				assert.Equal(t, middleName, *f.MiddleName)
				require.NotNil(t, f.LastName)
				assert.Equal(t, lastName, *f.LastName)
				assert.Equal(t, 5, f.Page)
				assert.Equal(t, 15, f.Limit)
			},
			assertResult: func(t *testing.T, gotErr error, got []response.PatientResponse, total int64) {
				assert.NoError(t, gotErr)
				assert.Equal(t, int64(0), total)
				assert.Empty(t, got)
			},
		},
		{
			name: "keeps zero pagination when request does not send values",
			req: request.SearchPatientRequest{
				Hospital: "Hospital A",
			},
			assertFilter: func(t *testing.T, f repository.PatientSearchFilter) {
				assert.Equal(t, "Hospital A", f.Hospital)
				assert.Nil(t, f.PassportID)
				assert.Nil(t, f.NationalID)
				assert.Equal(t, "", f.FirstName)
				assert.Nil(t, f.MiddleName)
				assert.Nil(t, f.LastName)
				assert.Nil(t, f.DateOfBirth)
				assert.Nil(t, f.PhoneNumber)
				assert.Nil(t, f.Email)
				assert.Equal(t, 0, f.Page)
				assert.Equal(t, 0, f.Limit)
			},
			assertResult: func(t *testing.T, gotErr error, got []response.PatientResponse, total int64) {
				assert.NoError(t, gotErr)
				assert.Equal(t, int64(0), total)
				assert.Empty(t, got)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockPatientRepo{
				searchFn: func(ctx context.Context, filter repository.PatientSearchFilter) ([]model.Patient, int64, error) {
					tt.assertFilter(t, filter)
					return tt.repoResult, tt.repoTotal, tt.repoErr
				},
			}

			svc := service.NewPatientService(repo)
			got, total, err := svc.Search(context.Background(), tt.req)
			tt.assertResult(t, err, got, total)
		})
	}
}

func TestPatientServiceGetPatientByID(t *testing.T) {
	nationalID := "1111111111111"
	patientID := uuid.MustParse("26fd07af-0d69-4ff6-a594-3e2f7e56fc2c")
	patientDOB := time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name         string
		inputID      string
		repoResult   model.Patient
		repoErr      error
		assertResult func(t *testing.T, got response.PatientResponse, err error)
	}{
		{
			name:    "success maps repository patient to response",
			inputID: nationalID,
			repoResult: model.Patient{
				ID:          patientID,
				Hospital:    "Hospital A",
				PatientHN:   "HN001",
				NationalID:  &nationalID,
				FirstNameTH: "สมชาย",
				LastNameTH:  "ใจดี",
				DateOfBirth: patientDOB,
				Gender:      "M",
			},
			assertResult: func(t *testing.T, got response.PatientResponse, err error) {
				assert.NoError(t, err)
				assert.Equal(t, patientID, got.ID)
				assert.Equal(t, "Hospital A", got.Hospital)
				assert.Equal(t, "HN001", got.PatientHN)
				require.NotNil(t, got.NationalID)
				assert.Equal(t, nationalID, *got.NationalID)
			},
		},
		{
			name:       "returns repository error as-is",
			inputID:    "P00000001",
			repoErr:    errors.New("record not found"),
			assertResult: func(t *testing.T, got response.PatientResponse, err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, "record not found")
				assert.Equal(t, response.PatientResponse{}, got)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockPatientRepo{
				searchFn: func(ctx context.Context, filter repository.PatientSearchFilter) ([]model.Patient, int64, error) {
					return nil, 0, nil
				},
				getByIDFn: func(ctx context.Context, id string) (model.Patient, error) {
					assert.Equal(t, tt.inputID, id)
					return tt.repoResult, tt.repoErr
				},
			}

			svc := service.NewPatientService(repo)
			got, err := svc.GetPatientByID(context.Background(), tt.inputID)
			tt.assertResult(t, got, err)
		})
	}
}

