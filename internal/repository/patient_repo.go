package repository

import (
	"context"
	"errors"

	"github.com/jeagerism/agnos-backend-assignment/internal/model"
	"gorm.io/gorm"
)

const (
	defaultSearchLimit = 50
	maxSearchLimit     = 200
)

var ErrHospitalRequired = errors.New("hospital is required")

type PatientSearchFilter struct {
	Hospital string

	PassportID *string
	NationalID *string

	FirstName   string
	MiddleName  *string
	LastName    *string
	DateOfBirth *string

	PhoneNumber *string
	Email       *string
	Page        int
	Limit       int
}

type PatientRepository interface {
	Search(ctx context.Context, filter PatientSearchFilter) ([]model.Patient, int64, error)
	GetPatientByID(ctx context.Context, id string) (model.Patient, error)
}

type patientRepository struct {
	db *gorm.DB
}

func NewPatientRepository(db *gorm.DB) PatientRepository {
	return &patientRepository{db: db}
}

func (r *patientRepository) GetPatientByID(ctx context.Context, id string) (response model.Patient,err error) {
	if err := r.db.WithContext(ctx).Where("national_id = ? OR passport_id = ?", id, id).First(&response).Error; err != nil {
		return response, err
	}
	return response, nil
}

func (r *patientRepository) Search(ctx context.Context, filter PatientSearchFilter) ([]model.Patient, int64, error) {
	if filter.Hospital == "" {
		return nil, 0, ErrHospitalRequired
	}

	var patients []model.Patient
	query := r.db.WithContext(ctx).Model(&model.Patient{}).Where("hospital = ?", filter.Hospital)

	limit, offset := normalizePagination(filter.Page, filter.Limit)

	// search by national ID or passport ID (exact match)
	query = addExactStringFilter(query, "national_id", filter.NationalID)
	query = addExactStringFilter(query, "passport_id", filter.PassportID)

	// first name supports single input across TH/EN
	query = addNameLikeFilter(query, "first_name_th", "first_name_en", filter.FirstName)
	query = addOptionalNameLikeFilter(query, "middle_name_th", "middle_name_en", filter.MiddleName)
	query = addOptionalNameLikeFilter(query, "last_name_th", "last_name_en", filter.LastName)

	// date of birth / contact
	query = addExactStringFilter(query, "date_of_birth", filter.DateOfBirth)
	query = addStringFilter(query, "phone_number", filter.PhoneNumber)
	query = addStringFilter(query, "email", filter.Email)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&patients).Error; err != nil {
		return nil, 0, err
	}
	return patients, total, nil
}

func addExactStringFilter(query *gorm.DB, column string, value *string) *gorm.DB {
	if value == nil {
		return query
	}
	return query.Where(column+" = ?", *value)
}

func addStringFilter(query *gorm.DB, column string, value *string) *gorm.DB {
	return addExactStringFilter(query, column, value)
}

func addNameLikeFilter(query *gorm.DB, colTH, colEN, value string) *gorm.DB {
	if value == "" {
		return query
	}
	keyword := "%" + value + "%"
	return query.Where("("+colTH+" ILIKE ? OR "+colEN+" ILIKE ?)", keyword, keyword)
}

func addOptionalNameLikeFilter(query *gorm.DB, colTH, colEN string, value *string) *gorm.DB {
	if value == nil {
		return query
	}
	return addNameLikeFilter(query, colTH, colEN, *value)
}

func normalizePagination(page, limit int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = defaultSearchLimit
	}
	if limit > maxSearchLimit {
		limit = maxSearchLimit
	}
	offset := (page - 1) * limit
	return limit, offset
}
