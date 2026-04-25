package model

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Patient struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Hospital     string    `gorm:"index;not null;uniqueIndex:idx_hosp_hn" json:"hospital"`
	PatientHN    string    `gorm:"not null;uniqueIndex:idx_hosp_hn" json:"patient_hn"`
	
	// use *string for optional fields to allow NULL in DB
	NationalID   *string   `json:"national_id"`
	PassportID   *string   `json:"passport_id"`
	
	FirstNameTH  string    `gorm:"not null" json:"first_name_th"`
	MiddleNameTH *string   `json:"middle_name_th"`
	LastNameTH   string    `gorm:"not null" json:"last_name_th"`
	
	FirstNameEN  *string   `json:"first_name_en"`
	MiddleNameEN *string   `json:"middle_name_en"`
	LastNameEN   *string   `json:"last_name_en"`
	
	DateOfBirth  time.Time `gorm:"type:date;not null" json:"date_of_birth"`
	
	PhoneNumber  *string   `json:"phone_number"`
	Email        *string   `json:"email"`
	
	Gender       string    `gorm:"not null;check:gender IN ('M','F')" json:"gender"` // M, F (mandatory)
	
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (p *Patient) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}