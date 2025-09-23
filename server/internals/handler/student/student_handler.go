package student_handler

import (

	"net/http"
	
	"github.com/labstack/echo/v4"
	"github.com/suhas-developer07/Smart-Attendence-System/server/internals/domain"
	student_service "github.com/suhas-developer07/Smart-Attendence-System/server/internals/service/student"
)

type StudentHandler struct {
	StudentService *student_service.StudentService
}

func NewStudentHandler(r *student_service.StudentService) *StudentHandler {
	return &StudentHandler{StudentService: r}
}


func (h *StudentHandler)StudentRegisterHandler(c echo.Context)error{

	var req domain.StudentRegisterPayload

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "Invalid request payload: " + err.Error(),
		})
	}

	id, err := h.StudentService.RegisterStudent(req)
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


func (h *StudentHandler) LoginStudentHandler(c echo.Context) error {
	var req domain.StudentLoginPayload
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "Invalid request payload: " + err.Error(),
		})
	}

	// Authenticate student
	token, err := h.StudentService.LoginStudent(req.USN, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, domain.ErrorResponse{
			Status: "error",
			Error:  "Authentication failed: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.SuccessResponse{
		Status:  "success",
		Message: "Login successful",
		Data:    map[string]string{"token": token},
	})
}

// func (h *StudentHandler) GetStudentsByDeptAndSemHandler(c echo.Context) error {
// 	department := c.QueryParam("department")
// 	semStr := c.QueryParam("sem")
// 	sem, err := strconv.Atoi(semStr)
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
// 			Status: "error",
// 			Error:  "Invalid sem value: " + err.Error(),
// 		})
// 	}

// 	students, err := h.StudentService.GetStudentsByDeptAndSem(department, sem)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
// 			Status: "error",
// 			Error:  "Failed to get students: " + err.Error(),
// 		})
// 	}

// 	return c.JSON(http.StatusOK, domain.SuccessResponse{
// 		Status:  "success",
// 		Message: "Students retrieved successfully",
// 		Data:    students,
// 	})
// }

// func (h *StudentHandler) GetSubjectsByStudentIDHandler(c echo.Context) error {
// 	studentIDStr := c.Param("student_id")
// 	studentID, err := strconv.ParseInt(studentIDStr, 10, 64)
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
// 			Status: "error",
// 			Error:  "Invalid student ID: " + err.Error(),
// 		})
// 	}

// 	subjects, err := h.StudentService.GetSubjectsByStudentID(studentID)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
// 			Status: "error",
// 			Error:  "Failed to get subjects: " + err.Error(),
// 		})
// 	}

// 	return c.JSON(http.StatusOK, domain.SuccessResponse{
// 		Status:  "success",
// 		Message: "Subjects retrieved successfully",
// 		Data:    subjects,
// 	})
// }

