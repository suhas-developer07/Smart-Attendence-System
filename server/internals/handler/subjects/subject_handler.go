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
	facultyIDFromJWT,ok := c.Get("faculty_id").(int64)

	if !ok {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "facultyID is not getting from jwt",
		})
	}


	subjects, err := h.SubjectService.GetSubjectsByFacultyID(facultyIDFromJWT)
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

func (h *SubjectHandler) GetSubjectsByStudentIDHandler(c echo.Context) error {
	// studentIDParam := c.Param("student_id")
	
	studentIDFromJWT , ok := c.Get("student_id").(int64)

	if !ok {
		return c.JSON(http.StatusBadRequest,domain.ErrorResponse{
			Status: "error",
			Error: "student id is not geting",
		})
	}
	// if studentIDParam == "" {
	// 	return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
	// 		Status: "error",
	// 		Error:  "student_id query parameter is required",
	// 	})
	// }

	// studentID, err := strconv.ParseInt(studentIDFromJWT, 10, 64)
	// if err != nil {
	// 	return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
	// 		Status: "error",
	// 		Error:  "invalid student_id parameter",
	// 	})
	// }

	subjects, err := h.SubjectService.GetSubjectsByStudentID(studentIDFromJWT)
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
