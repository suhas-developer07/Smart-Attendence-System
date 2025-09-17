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

func NewAttendanceHandler(ar *attendence_service.AttendanceService) *AttendanceHandler{
	return &AttendanceHandler{
		AttendanceService: ar,
	}
}

func (h *AttendanceHandler) MarkAttendanceHandler(c echo.Context)error {
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
		Data:   map[string]int64{"attendance_id": id},
	})
}

func (h *AttendanceHandler) GetAttendanceByStudentAndSubjectHandler(c echo.Context)error {
	var studentID = c.QueryParam("student_id")
	var subjectID = c.QueryParam("subject_id")

	if studentID == "" || subjectID == "" {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "student_id and subject_id are required",
		})
	}
	studentIDInt, err := strconv.ParseInt(studentID, 10, 64)
	
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "Invalid studentID" + err.Error(),
		})
	}
	subjectIDInt, err := strconv.ParseInt(subjectID, 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "Invalid subjectID" + err.Error(),
		})
	}

	attendances, err := h.AttendanceService.GetAttendanceByStudentAndSubject(studentIDInt, subjectIDInt)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Status: "error",
			Error:  "Failed to fetch attendance" + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.SuccessResponse{
		Status:  "success",
		Message: "Attendance fetched successfully",
		Data:   attendances,
	})
}

func (h *AttendanceHandler) GetAttendanceBySubjectHandler(c echo.Context)error {
	var subjectID = c.QueryParam("subject_id")
	var fromDate = c.QueryParam("from_date")
	var toDate = c.QueryParam("to_date")

	if subjectID == "" || fromDate == "" || toDate == "" {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "subject_id, from_date and to_date are required",
		})
	}

	// Parse subjectID to int64
	subjectIDInt, err := strconv.ParseInt(subjectID, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "invalid subject_id" + err.Error(),
		})
	}
	
	// Parse fromDate and toDate to time.Time
	fromDateTime, err := time.Parse("2006-01-02", fromDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "invalid from_date" + err.Error(),
		})
	}
	toDateTime, err := time.Parse("2006-01-02", toDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "invalid to_date" + err.Error(),
		})
	}

	attendances, err := h.AttendanceService.GetAttendanceBySubject(subjectIDInt, fromDateTime, toDateTime)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Status: "error",
			Error:  "Failed to fetch attendance" + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.SuccessResponse{
		Status:  "success",
		Message: "Attendance fetched successfully",
		Data:   attendances,
	})
}

func (h *AttendanceHandler) AssignSubjectToTimeRangeHandler(c echo.Context)error {
	var req struct {
		FacultyID int    `json:"faculty_id"`
		SubjectID int64  `json:"subject_id"`
		Start     string `json:"start"`
		End       string `json:"end"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "invalid request payload" + err.Error(),
		})
	}

	if req.FacultyID == 0 || req.SubjectID == 0 || req.Start == "" || req.End == "" {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "faculty_id, subject_id, start and end are required",
		})
	}

	startTime, err := time.Parse("2006-01-02", req.Start)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "invalid start date" + err.Error(),
		})
	}
	endTime, err := time.Parse("2006-01-02", req.End)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "invalid end date" + err.Error(),
		})
	}

	updatedCount, skipped, err := h.AttendanceService.AssignSubjectToTimeRange(req.FacultyID, req.SubjectID, startTime, endTime)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Status: "error",
			Error:  "Failed to assign subject to time range" + err.Error(),
		})
	}
	
	return c.JSON(http.StatusOK, domain.SuccessResponse{
		Status:  "success",
		Message: "Subject assigned to time range",
		Data: map[string]int64{
			"updatedCount": updatedCount,
			"skipped":      skipped,
		},
	},
)
}

func (h *AttendanceHandler) GetAttendanceSummaryBySubjectHandler(c echo.Context) error {
	var subjectID = c.QueryParam("subject_id")

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
		Data:   summaries,
	})
}

func (h *AttendanceHandler) GetClassAttendanceHandler(c echo.Context)error {
	var subjectID = c.QueryParam("subject_id")
	var date = c.QueryParam("date")

	if subjectID == "" || date == ""{
		return c.JSON(http.StatusBadRequest,domain.ErrorResponse{
			Status : "error",
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

	dateTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "invalid date" + err.Error(),
		})
	}

	ClassAttendance, err := h.AttendanceService.GetClassAttendance(subjectIDInt,dateTime)

	if err != nil {
		return c.JSON(http.StatusBadRequest,domain.ErrorResponse{
			Status: "error",
		    Error: "failed to get Class Attendance" +err.Error(),
		})
	}

	return  c.JSON(http.StatusOK, domain.SuccessResponse{
		Status: "success",
		Message: "Attendance fetched successfull",
		Data: ClassAttendance,
	})
}

func (h *AttendanceHandler) GetStudentAttendanceHistoryHandler(c echo.Context) error {
	var studentID = c.QueryParam("student_id")
	var subjectID = c.QueryParam("subject_id")

	if subjectID == "" || studentID == ""{
		return c.JSON(http.StatusBadRequest,domain.ErrorResponse{
			Status: "error",
			Error: "student_id and subject_id are required",
		})
	}

	studentIDInt, err := strconv.ParseInt(studentID, 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "Invalid studentID" + err.Error(),
		})
	}
	subjectIDInt, err := strconv.ParseInt(subjectID, 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status: "error",
			Error:  "Invalid subjectID" + err.Error(),
		})
	}

	attendanceHistory ,err  := h.AttendanceService.GetStudentAttendanceHistory(studentIDInt,subjectIDInt)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Status: "error",
			Error:  "Failed to fetch attendance history" + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.SuccessResponse{
		Status:  "success",
		Message: "Attendance history fetched successfully",
		Data:   attendanceHistory,
	})
}
