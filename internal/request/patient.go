package request

// SearchPatientRequest อิงตาม Task 4
type SearchPatientRequest struct {
	Hospital string `json:"-"`

	PassportID  *string `form:"passport_id"`
	NationalID  *string `form:"national_id"`
	FirstName   string  `form:"first_name"`
	MiddleName  *string `form:"middle_name"`
	LastName    *string `form:"last_name"`
	DateOfBirth *string `form:"date_of_birth"`
	PhoneNumber *string `form:"phone_number"`
	Email       *string `form:"email"`

	Page   int `form:"page" default:"1"`
	Limit  int `form:"limit" default:"10"`
}
