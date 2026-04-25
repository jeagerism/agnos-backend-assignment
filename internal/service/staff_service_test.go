package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jeagerism/agnos-backend-assignment/internal/model"
	"github.com/jeagerism/agnos-backend-assignment/internal/request"
	"github.com/jeagerism/agnos-backend-assignment/internal/service"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// mockStaffRepo จำลอง repository.StaffRepository โดยไม่ต้องเชื่อมต่อ DB จริง
// แต่ละ field เป็น function ที่เราสามารถกำหนดพฤติกรรมได้ในแต่ละ test case
type mockStaffRepo struct {
	createFn         func(ctx context.Context, s *model.Staff) error
	findByUsernameFn func(ctx context.Context, username string) (*model.Staff, error)
}

func (m *mockStaffRepo) Create(ctx context.Context, s *model.Staff) error {
	return m.createFn(ctx, s)
}

func (m *mockStaffRepo) FindByUsername(ctx context.Context, username string) (*model.Staff, error) {
	return m.findByUsernameFn(ctx, username)
}

const testSecret = "test-secret"

// ─── Create ───────────────────────────────────────────────────────────────────

func TestCreate(t *testing.T) {
	tests := []struct {
		name      string
		setupRepo func() *mockStaffRepo
		req       request.CreateStaffRequest
		wantErr   bool
	}{
		{
			name: "success: password is hashed before saving",
			setupRepo: func() *mockStaffRepo {
				return &mockStaffRepo{
					createFn: func(ctx context.Context, s *model.Staff) error {
						// ตรวจสอบว่า password ถูก hash แล้ว (ไม่ใช่ plain text)
						assert.NotEqual(t, "password123", s.Password)
						// ตรวจสอบว่า hash นั้น verify กับ plain text ได้
						err := bcrypt.CompareHashAndPassword([]byte(s.Password), []byte("password123"))
						assert.NoError(t, err)
						return nil
					},
				}
			},
			req:     request.CreateStaffRequest{Username: "staff_a", Password: "password123", Hospital: "Hospital A"},
			wantErr: false,
		},
		{
			name: "fail: duplicate username returns error",
			setupRepo: func() *mockStaffRepo {
				return &mockStaffRepo{
					createFn: func(ctx context.Context, s *model.Staff) error {
						return errors.New("duplicate key value violates unique constraint")
					},
				}
			},
			req:     request.CreateStaffRequest{Username: "staff_a", Password: "password123", Hospital: "Hospital A"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewStaffService(tt.setupRepo(), testSecret)
			err := svc.Create(context.Background(), tt.req)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// ─── Login ────────────────────────────────────────────────────────────────────

func TestLogin(t *testing.T) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	tests := []struct {
		name       string
		setupRepo  func() *mockStaffRepo
		req        request.LoginStaffRequest
		wantToken  bool
		wantErr    error
	}{
		{
			name: "success: valid credentials return JWT token",
			setupRepo: func() *mockStaffRepo {
				return &mockStaffRepo{
					findByUsernameFn: func(ctx context.Context, username string) (*model.Staff, error) {
						// จำลอง staff ที่มีอยู่ใน DB พร้อม hashed password
						return &model.Staff{Username: username, Password: string(hashed), Hospital: "Hospital A"}, nil
					},
				}
			},
			req:      request.LoginStaffRequest{Username: "staff_a", Password: "password123"},
			wantToken: true,
			wantErr:  nil,
		},
		{
			name: "fail: wrong password returns ErrWrongPassword",
			setupRepo: func() *mockStaffRepo {
				return &mockStaffRepo{
					findByUsernameFn: func(ctx context.Context, username string) (*model.Staff, error) {
						return &model.Staff{Username: username, Password: string(hashed), Hospital: "Hospital A"}, nil
					},
				}
			},
			req:      request.LoginStaffRequest{Username: "staff_a", Password: "wrongpassword"},
			wantToken: false,
			// service คืน ErrWrongPassword → handler map เป็น "invalid credentials" ก่อน response
			wantErr: service.ErrWrongPassword,
		},
		{
			name: "fail: user not found returns ErrUserNotFound",
			setupRepo: func() *mockStaffRepo {
				return &mockStaffRepo{
					findByUsernameFn: func(ctx context.Context, username string) (*model.Staff, error) {
						return nil, errors.New("record not found")
					},
				}
			},
			req:      request.LoginStaffRequest{Username: "unknown_user", Password: "password123"},
			wantToken: false,
			// service คืน ErrUserNotFound → handler map เป็น "invalid credentials" ก่อน response
			// ทำให้ client ไม่รู้ว่า username มีอยู่จริงหรือไม่ (ป้องกัน username enumeration)
			wantErr: service.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewStaffService(tt.setupRepo(), testSecret)
			token, err := svc.Login(context.Background(), tt.req)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}
		})
	}
}
