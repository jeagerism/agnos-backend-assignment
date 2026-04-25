package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jeagerism/agnos-backend-assignment/internal/request"
	"github.com/jeagerism/agnos-backend-assignment/internal/response"
	"github.com/jeagerism/agnos-backend-assignment/internal/service"
)

type StaffHandler struct {
	service service.StaffService
}

func NewStaffHandler(svc service.StaffService) *StaffHandler {
	return &StaffHandler{service: svc}
}

func (h *StaffHandler) Create(c *gin.Context) {
	var req request.CreateStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Create(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response.CreateStaffResponse{
		Message: "staff created successfully",
	})
}

func (h *StaffHandler) Login(c *gin.Context) {
	var req request.LoginStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		// log real error to know where it failed (server side only)	
		log.Printf("[login] failed for user %q: %v", req.Username, err)
		// response out as generic to prevent username enumeration
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, response.LoginStaffResponse{
		Token: token,
	})
}
