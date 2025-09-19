package faculty

import (
	"net/http"
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

func (h *FacultyHandler) AuthenticateFacultyHandler(c echo.Context) error {
	var req domain.FacultyLoginPayload

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "invalid request payload" + err.Error(),
		})
	}

	token, err := h.FacultyService.AuthenticateFaculty(req)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, domain.ErrorResponse{
			Status: "error",
			Error:  "Authentication failed: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.SuccessResponse{
		Status:  "success",
		Message: "Authentication successful",
		Data:    map[string]string{"token": token},
	})
}
func (h *FacultyHandler) GetFacultyByIDHandler(c echo.Context) error {
	facultyID := c.Get("faculty_id").(int64)

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

func (h *FacultyHandler) GetAllFacultyHandler(c echo.Context) error {
	faculties, err := h.FacultyService.GetAllFaculty()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Status: "error",
			Error:  "Failed to get faculties: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.SuccessResponse{
		Status:  "success",
		Message: "Faculties retrieved successfully",
		Data:    faculties,
	})
}

func (h *FacultyHandler) GetFacultyByDepartmentHandler(c echo.Context) error {
	department := c.Param("dept")
	if department == "" {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "department query parameter is required",
		})
	}

	faculties, err := h.FacultyService.GetFacultyByDepartment(department)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Status: "error",
			Error:  "Failed to get faculties: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.SuccessResponse{
		Status:  "success",
		Message: "Faculties retrieved successfully",
		Data:    faculties,
	})
}