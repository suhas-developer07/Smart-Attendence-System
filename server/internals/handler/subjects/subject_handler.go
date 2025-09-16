package subject_handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/suhas-developer07/Smart-Attendence-System/server/internals/domain"
	"github.com/suhas-developer07/Smart-Attendence-System/server/internals/service/subject"
)


type SubjectHandler struct {
	SubjectService *subject_service.SubjectService
}

func NewSubjectHandler(fr *subject_service.SubjectService) *SubjectHandler {
	return &SubjectHandler{
		SubjectService: fr,
	}
}


func (h *SubjectHandler) AddSubjectHandler(c echo.Context) error {

	var req domain.SubjectPayload

     if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
            Status: "error",
            Error:  "invalid request payload"+err.Error(),
        })
	}
	id, err := h.SubjectService.AddSubject(req)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Status: "error",
			Error:  "Failed to add the subjects" + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.SuccessResponse{
		Status:  "success",
		Message: "Subjects added successfully",
		Data:    map[string]int64{"subject_id": id},
	})
}

func (h *SubjectHandler) GetSubjectsByDeptAndSemHandler(c echo.Context) error {
	department := c.QueryParam("department")
	semParam := c.QueryParam("sem")

	if department == "" || semParam == "" {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "department and sem query parameters are required",
		})
	}

	sem, err := strconv.Atoi(semParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "invalid sem parameter",
		})
	}

	subjects, err := h.SubjectService.GetSubjectsByDeptAndSem(department, sem)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Status: "error",
			Error:  "Failed to fetch subjects" + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.SuccessResponse{
		Status:  "success",
		Message: "Fetched subjects successfully",
		Data:    subjects,
	})
}

func (h *SubjectHandler) GetSubjectsByFacultyIDHandler(c echo.Context) error {
	facultyIDParam := c.QueryParam("faculty_id")

	if facultyIDParam == "" {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "faculty_id query parameter is required",
		})
	}

	facultyID, err := strconv.Atoi(facultyIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "invalid faculty_id parameter",
		})
	}

	subjects, err := h.SubjectService.GetSubjectsByFacultyID(facultyID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Status: "error",
			Error:  "Failed to fetch subjects" + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.SuccessResponse{
		Status:  "success",
		Message: "Fetched subjects successfully",
		Data:    subjects,
	})
}

