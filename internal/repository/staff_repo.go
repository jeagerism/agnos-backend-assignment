package repository

import (
	"context"

	"github.com/jeagerism/agnos-backend-assignment/internal/model"
	"gorm.io/gorm"
)

type StaffRepository interface {
	Create(ctx context.Context, staff *model.Staff) error
	FindByUsername(ctx context.Context, username string) (*model.Staff, error)
}

type staffRepository struct {
	db *gorm.DB
}

func NewStaffRepository(db *gorm.DB) StaffRepository {
	return &staffRepository{db: db}
}

func (r *staffRepository) Create(ctx context.Context, staff *model.Staff) error {
	return r.db.WithContext(ctx).Create(staff).Error
}

func (r *staffRepository) FindByUsername(ctx context.Context, username string) (*model.Staff, error) {
	var s model.Staff
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}
