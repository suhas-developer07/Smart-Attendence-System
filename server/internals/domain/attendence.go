package domain

import "time"

type Attendance struct {
	ID string `json:"attendance_id"`
	StudentID   string `json:"student_id"`
	SubjectID   string `json:"subject_id"`
	Date         string `json:"date"`
	Status       string `json:"status"`
	RecordedAt string `json:"recorded_at"`
	CreatedAt string `json:"created_at"`
}

type AttendancePayload struct {
	StudentID int64  `json:"student_id" validate:"required"`
	SubjectID int64  `json:"subject_id" validate:"required"`
	Status    string `json:"status" validate:"required,oneof=Present Absent"`
	RecordedAt time.Time `json:"recorded_at" validate:"required"`
}

type SubjectSummary struct {
    SubjectID     int64   `json:"subject_id"`
    SubjectName   string  `json:"subject_name"`
    TotalClasses  int     `json:"total_classes"`
    Attended      int     `json:"attended"`
    Percentage    float64 `json:"percentage"`
}

type StudentSummary struct {
    StudentID    int64   `json:"student_id"`
    StudentName  string  `json:"student_name"`
    TotalClasses int     `json:"total_classes"`
    Attended     int     `json:"attended"`
    Percentage   float64 `json:"percentage"`
}

type ClassAttendance struct {
    StudentID   int64     `json:"student_id"`
    StudentName string    `json:"student_name"`
    Date        time.Time `json:"date"`
    Status      string    `json:"status"`
}


// domain/attendance.go

type AttendanceWithNames struct {
    ID          int64     `json:"attendance_id"`
    StudentID   int64     `json:"student_id"`
    StudentName string    `json:"student_name"`
    SubjectID   int64     `json:"subject_id"`
    SubjectName string    `json:"subject_name"`
    Date        time.Time `json:"date"`
    Status      string    `json:"status"`
    RecordedAt  time.Time `json:"recorded_at"`
    CreatedAt   time.Time `json:"created_at"`
}

type StudentHistory struct {
    ID          int64     `json:"attendance_id"`
    Date        time.Time `json:"date"`
    Status      string    `json:"status"`
    SubjectID   int64     `json:"subject_id"`
    SubjectName string    `json:"subject_name"`
    RecordedAt  time.Time `json:"recorded_at"`
}


type AttendanceRepository interface {
	MarkAttendance(attendance *AttendancePayload) (int64, error)
	GetAttendanceByStudentAndSubject(studentID, subjectID int64) ([]AttendanceWithNames, error)
	GetAttendanceBySubject(subjectID int64, fromDate, toDate time.Time) ([]AttendanceWithNames, error)
	AssignSubjectToTimeRange(facultyID int, subjectID int64, start, end time.Time) (int64, int64, error)
	GetAttendanceSummaryBySubject(subjectID int64) ([]StudentSummary, error)
	GetClassAttendance(subjectID int64, date time.Time) ([]ClassAttendance, error)
	GetStudentAttendanceHistory(studentID int64, subjectID int64) ([]StudentHistory, error) 

}