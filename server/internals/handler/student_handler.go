package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/suhas-developer07/Smart-Attendence-System/server/internals/domain"
	"github.com/suhas-developer07/Smart-Attendence-System/server/internals/service"
)

type StudentHandler struct {
	repo *service.StudentService
}

func NewStudentHandler(r *service.StudentService) *StudentHandler {
	return &StudentHandler{repo: r}
}

func (h *StudentHandler) StudentRegisterHandler(c echo.Context) error{
	var req domain.StudentRegisterPayload

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "Invalid body format",
		})
		
	}
	id, err := h.repo.RegisterStudentService(req)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Status: "error",
			Error:  "Failed to register student: " + err.Error(),
		})
		
	}
	return c.JSON(http.StatusOK, domain.SuccessResponse{
		Status:  "success",
		Message: "Student registered successfully",
		Data:    map[string]int64{"student_id": id},
	})

}
