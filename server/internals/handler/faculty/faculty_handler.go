package faculty_handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/suhas-developer07/Smart-Attendence-System/server/internals/domain"
	faculty_service "github.com/suhas-developer07/Smart-Attendence-System/server/internals/service/faculty"
)


type FacultyHandler struct{
     facultyRepo *faculty_service.FacultyService

}

func NewFacultyHandler(fr *faculty_service.FacultyService) *FacultyHandler 	{
	return &FacultyHandler{
		facultyRepo: fr,
	}
}


func (h *FacultyHandler) AddSubject(c echo.Context) error {

	var req domain.SubjectPayload

     if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
            Status: "error",
            Error:  "invalid request payload"+err.Error(),
        })
	}
	id, err := h.facultyRepo.AddSubject(req)

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


