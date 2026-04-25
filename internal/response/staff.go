package response

type CreateStaffResponse struct {
	Message string `json:"message"`
}

type LoginStaffResponse struct {
	Token string `json:"token"`
}
