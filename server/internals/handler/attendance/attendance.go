package attendance_handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/suhas-developer07/Smart-Attendence-System/server/internals/domain"
	attendence_service "github.com/suhas-developer07/Smart-Attendence-System/server/internals/service/attendence"
)

type AttendanceHandler struct {
	AttendanceService *attendence_service.AttendanceService
}

func NewAttendanceHandler(ar *attendence_service.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{
		AttendanceService: ar,
	}
}

func (h *AttendanceHandler) MarkAttendanceHandler(c echo.Context) error {
	var req domain.AttendancePayload

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "invalid request payload" + err.Error(),
		})
	}
	id, err := h.AttendanceService.MarkAttendance(&req)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Status: "error",
			Error:  "Failed to mark attendance" + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.SuccessResponse{
		Status:  "success",
		Message: "Attendance marked successfully",
		Data:    map[string]int64{"attendance_id": id},
	})
}

func (h *AttendanceHandler) GetAttendanceByStudentAndSubjectHandler(c echo.Context) error {
	var usn = c.QueryParam("usn")
	var subjectID = c.QueryParam("subject_id")

	if usn == "" || subjectID == "" {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "usn and subject_id are required",
		})
	}

	subjectIDInt, err := strconv.ParseInt(subjectID, 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "Invalid subjectID" + err.Error(),
		})
	}

	attendances, err := h.AttendanceService.GetAttendanceByStudentAndSubject(usn, subjectIDInt)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Status: "error",
			Error:  "Failed to fetch attendance" + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.SuccessResponse{
		Status:  "success",
		Message: "Attendance fetched successfully",
		Data:    attendances,
	})
}

//this function returns attendance of a subject on a particular date with all fields in the AttendanceWithNames struct struct 
func (h *AttendanceHandler) GetAttendanceBySubjectAndDateHandler(c echo.Context) error {
	subjectID := c.QueryParam("subject_id")
	dateStr := c.QueryParam("date") // e.g. "2025-06-12"

	if subjectID == "" || dateStr == "" {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "subject_id and date are required",
		})
	}

	subjectIDInt, err := strconv.ParseInt(subjectID, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "invalid subject_id: " + err.Error(),
		})
	}

	// Parse only the date part (no time, assume YYYY-MM-DD)
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "invalid date: " + err.Error(),
		})
	}

	attendances, err := h.AttendanceService.GetAttendanceBySubjectAndDate(subjectIDInt, date)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Status: "error",
			Error:  "failed to fetch attendance: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.SuccessResponse{
		Status:  "success",
		Message: "Attendance fetched successfully",
		Data:    attendances,
	})
}
func (h *AttendanceHandler) AssignSubjectToTimeRangeHandler(c echo.Context) error {
	var req struct {
		FacultyID int    `json:"faculty_id"`
		SubjectID int64  `json:"subject_id"`
		ClassDate string `json:"class_date"` // YYYY-MM-DD
		Start     string `json:"start"`      // HH:MM
		End       string `json:"end"`        // HH:MM
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error", Error: "invalid request: " + err.Error(),
		})
	}

	loc, _ := time.LoadLocation("Asia/Kolkata")
	classDate, _ := time.ParseInLocation("2006-01-02", req.ClassDate, loc)
	startTime, _ := time.Parse("15:04", req.Start)
	endTime, _ := time.Parse("15:04", req.End)

	updated, skipped, err := h.AttendanceService.AssignSubjectToTimeRange(
		req.FacultyID, req.SubjectID, classDate, startTime, endTime,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Status: "error", Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.SuccessResponse{
		Status: "success",
		Message: "Subject assigned to attendance within time range",
		Data: map[string]int64{"updatedCount": updated, "skipped": skipped},
	})
}

//This function returns attendance summary of a subject means whole attendence of all students in that subject 
func (h *AttendanceHandler) GetAttendanceSummaryBySubjectHandler(c echo.Context) error {
	var subjectID = c.Param("subject_id")

	if subjectID == "" {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "subject_id is required",
		})
	}
	subjectIDInt, err := strconv.ParseInt(subjectID, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "invalid subject_id" + err.Error(),
		})
	}

	summaries, err := h.AttendanceService.GetAttendanceSummaryBySubject(subjectIDInt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Status: "error",
			Error:  "Failed to fetch attendance summary" + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.SuccessResponse{
		Status:  "success",
		Message: "Attendance summary fetched successfully",
		Data:    summaries,
	})
}

// this returns attendance of all students in a class on a particular date just returns usn,studentname,date,status
func (h *AttendanceHandler) GetClassAttendanceHandler(c echo.Context) error {
	subjectID := c.QueryParam("subject_id")
	date := c.QueryParam("date")

	if subjectID == "" || date == "" {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "subject_id and date are required",
		})
	}

	subjectIDInt, err := strconv.ParseInt(subjectID, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "invalid subject_id: " + err.Error(),
		})
	}

	dateTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "invalid date: " + err.Error(),
		})
	}

	classAttendance, err := h.AttendanceService.GetClassAttendance(subjectIDInt, dateTime)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Status: "error",
			Error:  "failed to get class attendance: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.SuccessResponse{
		Status:  "success",
		Message: "Attendance fetched successfully",
		Data:    classAttendance,
	})
}

//this function returns attendance history of a student in a particular subject. with all dates and status
func (h *AttendanceHandler) GetStudentAttendanceHistoryHandler(c echo.Context) error {
	var usn = c.QueryParam("usn")
	var subjectID = c.QueryParam("subject_id")

	if subjectID == "" || usn == "" {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "usn and subject_id are required",
		})
	}

	subjectIDInt, err := strconv.ParseInt(subjectID, 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "Invalid subjectID" + err.Error(),
		})
	}

	attendanceHistory, err := h.AttendanceService.GetStudentAttendanceHistory(usn, subjectIDInt)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Status: "error",
			Error:  "Failed to fetch attendance history" + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.SuccessResponse{
		Status:  "success",
		Message: "Attendance history fetched successfully",
		Data:    attendanceHistory,
	})
}
