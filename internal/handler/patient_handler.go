package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jeagerism/agnos-backend-assignment/internal/request"
	"github.com/jeagerism/agnos-backend-assignment/internal/response"
	"github.com/jeagerism/agnos-backend-assignment/internal/service"
)

type PatientHandler struct {
	service service.PatientService
}

const (
	defaultPage  = 1
	defaultLimit = 10
)

func NewPatientHandler(svc service.PatientService) *PatientHandler {
	return &PatientHandler{service: svc}
}


func (h *PatientHandler) Search(c *gin.Context) {
	var req request.SearchPatientRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Hospital = c.GetString("hospital")
	req.Page, req.Limit = normalizePageLimit(req.Page, req.Limit)
	log.Printf("\n\n handler layer : \n\n hospital: %s\n\n", req.Hospital)

	patients, total, err := h.service.Search(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.SearchPatientResponse{
		Message: "patients searched successfully",
		Data:    patients,
		Paginate: buildPaginate(req.Page, req.Limit, total),
	})
}


func (h *PatientHandler) GetPatientByID(c *gin.Context) {
	id := c.Param("id")
	patient, err := h.service.GetPatientByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.GetPatientByIDResponse{
		Message: "patient retrieved successfully",
		Data:    patient,
	})
}
func buildPaginate(page, limit int, total int64) response.PaginateResponse {
	return response.PaginateResponse{
		Page:    page,
		PerPage: limit,
		Total:   total,
	}
}

func normalizePageLimit(page, limit int) (int, int) {
	if page <= 0 {
		page = defaultPage
	}
	if limit <= 0 {
		limit = defaultLimit
	}
	return page, limit
}