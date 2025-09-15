package subject_handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/suhas-developer07/Smart-Attendence-System/server/internals/domain"
	"github.com/suhas-developer07/Smart-Attendence-System/server/internals/service/subject"
)


type SubjectHandler struct {
	SubjectRepo *subject_service.SubjectService
}

func NewSubjectHandler(fr *subject_service.SubjectService) *SubjectHandler {
	return &SubjectHandler{
		SubjectRepo: fr,
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
	id, err := h.SubjectRepo.AddSubject(req)

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

func (h *SubjectHandler) ListSubjectsHandler(c echo.Context) error {
	branch := c.QueryParam("branch")
	semParam := c.QueryParam("sem")

	if branch == "" || semParam == "" {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "branch and sem query parameters are required",
		})
	}

	sem, err := strconv.Atoi(semParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "invalid sem parameter",
		})
	}

	subjects, err := h.SubjectRepo.ListSubjects(branch, sem)
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


// func (h *FacultyHandler) GetStudentsWithSubjectsHandler(c echo.Context) error {
// 	var req domain.GetStudentsWithSubjectsPayload

// 	if err := c.Bind(&req); err != nil {
// 		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
// 			Status: "error",
// 			Error:  "invalid request payload" + err.Error(),
// 		})
// 	}

// 	students, err := h.facultyRepo.GetStudentsWithSubjects(req)

// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
// 			Status: "error",
// 			Error:  "Failed to fetch students with subjects" + err.Error(),
// 		})
// 	}

// 	return c.JSON(http.StatusOK, domain.SuccessResponse{
// 		Status:  "success",
// 		Message: "Fetched students with subjects successfully",
// 		Data:    students,
// 	})
// }

