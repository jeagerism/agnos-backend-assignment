package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jeagerism/agnos-backend-assignment/internal/model"
	"github.com/jeagerism/agnos-backend-assignment/internal/repository"
	"github.com/jeagerism/agnos-backend-assignment/internal/request"
	"golang.org/x/crypto/bcrypt"
)

// sentinel errors สำหรับแยกประเภทความผิดพลาดใน service layer
var (
	ErrUserNotFound  = errors.New("user not found")
	ErrWrongPassword = errors.New("wrong password")
)

type StaffService interface {
	Create(ctx context.Context, req request.CreateStaffRequest) error
	Login(ctx context.Context, req request.LoginStaffRequest) (string, error)
}

type staffService struct {
	repo      repository.StaffRepository
	jwtSecret string
}

func NewStaffService(repo repository.StaffRepository, jwtSecret string) StaffService {
	return &staffService{repo: repo, jwtSecret: jwtSecret}
}

func (s *staffService) Create(ctx context.Context, req request.CreateStaffRequest) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.Create(ctx, &model.Staff{
		Username: req.Username,
		Password: string(hashed),
		Hospital: req.Hospital,
	})
}

func (s *staffService) Login(ctx context.Context, req request.LoginStaffRequest) (string, error) {
	staff, err := s.repo.FindByUsername(ctx, req.Username)
	if err != nil {
		return "", ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(staff.Password), []byte(req.Password)); err != nil {
		return "", ErrWrongPassword
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"staff_id": staff.ID.String(),
		"hospital": staff.Hospital,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString([]byte(s.jwtSecret))
}
