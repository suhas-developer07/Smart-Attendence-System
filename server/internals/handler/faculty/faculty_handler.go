package faculty

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/suhas-developer07/Smart-Attendence-System/server/internals/domain"
	faculty_service "github.com/suhas-developer07/Smart-Attendence-System/server/internals/service/faculty"
)

type FacultyHandler struct {
	FacultyService *faculty_service.FacultyService
}

func NewFacultyHandler(fr *faculty_service.FacultyService) *FacultyHandler {
	return &FacultyHandler{
		FacultyService: fr,
	}
}

func (h *FacultyHandler) RegisterFacultyHandler(c echo.Context) error {
	var req domain.FacultyRegisterPayload
	
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "invalid request payload" + err.Error(),
		})
	}

	id ,err := h.FacultyService.RegisterFaculty(req);

	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Status: "error",
			Error:  "Failed to register faculty" + err.Error(),
		})
	}
	
	return c.JSON(http.StatusOK, domain.SuccessResponse{
		Status:  "success",
		Message: "Faculty registered successfully",
		Data:   map[string]int64{"faculty_id": id},
	})
}

func (h *FacultyHandler) GetFacultyByIDHandler(c echo.Context) error {
	facultyIDParam := c.Param("id")

	facultyID, err := strconv.Atoi(facultyIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "Invalid faculty ID: " + err.Error(),
		})
	}

	faculty, err := h.FacultyService.GetFacultyByID(facultyID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Status: "error",
			Error:  "Failed to get faculty: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.SuccessResponse{
		Status:  "success",
		Message: "Faculty retrieved successfully",
		Data:    faculty,
	})
}