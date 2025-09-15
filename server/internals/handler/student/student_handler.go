package student_handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/suhas-developer07/Smart-Attendence-System/server/internals/domain"
	student_service "github.com/suhas-developer07/Smart-Attendence-System/server/internals/service/student"
)

type StudentHandler struct {
	studentRepo *student_service.StudentService
}

func NewStudentHandler(r *student_service.StudentService) *StudentHandler {
	return &StudentHandler{studentRepo: r}
}
func (h *StudentHandler) StudentRegisterHandler(c echo.Context) error {

	semStr := c.FormValue("sem")
	sem, err := strconv.Atoi(semStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "Invalid sem value: " + err.Error(),
		})
	}
	// Parse form fields
	req := domain.StudentRegisterPayload{
		USN:      c.FormValue("usn"),
		Username: c.FormValue("username"),
		Department:   c.FormValue("department"),
		Sem:      sem,
	}

	// Create directory for this user
	userDir := filepath.Join("images", req.USN)
	if err := os.MkdirAll(userDir, os.ModePerm); err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Status: "error",
			Error:  "Failed to create Students directory: " + err.Error(),
		})
	}

	// Get uploaded images
	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "Invalid form data: " + err.Error(),
		})
	}

	files := form.File["images"]
	// if len(files) != 5 {
	// 	return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
	// 		Status: "error",
	// 		Error:  "You must upload exactly 5 images",
	// 	})
	// }

	// Save each image to userDir
	for i, fileHeader := range files {
		src, err := fileHeader.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
				Status: "error",
				Error:  "Failed to open uploaded file: " + err.Error(),
			})
		}
		defer src.Close()

		// Save with sequential naming
		dstPath := filepath.Join(userDir, fmt.Sprintf("image_%d.jpg", i+1))
		dst, err := os.Create(dstPath)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
				Status: "error",
				Error:  "Failed to save file: " + err.Error(),
			})
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
				Status: "error",
				Error:  "Failed to write file: " + err.Error(),
			})
		}
	}

	id, err := h.studentRepo.RegisterStudent(req)
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
