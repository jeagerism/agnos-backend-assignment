package response

import (
	"time"

	"github.com/google/uuid"
)

type SearchPatientResponse struct {
	Message  string             `json:"message"`
	Data     []PatientResponse  `json:"data"`
	Paginate PaginateResponse   `json:"paginate"`
}

type PaginateResponse struct {
	Page    int   `json:"page"`
	PerPage int   `json:"perPage"`
	Total   int64 `json:"total"`
}

type PatientResponse struct {
	ID           uuid.UUID `json:"id"`
	Hospital     string    `json:"hospital"`
	PatientHN    string    `json:"patient_hn"`
	NationalID   *string   `json:"national_id"`
	PassportID   *string   `json:"passport_id"`
	FirstNameTH  string    `json:"first_name_th"`
	MiddleNameTH *string   `json:"middle_name_th"`
	LastNameTH   string    `json:"last_name_th"`
	FirstNameEN  *string   `json:"first_name_en"`
	MiddleNameEN *string   `json:"middle_name_en"`
	LastNameEN   *string   `json:"last_name_en"`
	DateOfBirth  time.Time `json:"date_of_birth"`
	PhoneNumber  *string   `json:"phone_number"`
	Email        *string   `json:"email"`
	Gender       string    `json:"gender"`
}


type GetPatientByIDResponse struct {
	Message string         `json:"message"`
	Data    PatientResponse `json:"data"`
}