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
	USN      string    `json:"usn" validate:"required"`
	Status   string    `json:"status" validate:"required,oneof=Present Absent"`
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
    USN         string   `json:"usn"`
    StudentName  string  `json:"student_name"`
    TotalClasses int     `json:"total_classes"`
    Attended     int     `json:"attended"`
    Percentage   float64 `json:"percentage"`
}

type ClassAttendance struct {
    USN       string    `json:"usn"`
    StudentName string    `json:"student_name"`
    Date        time.Time `json:"date"`
    Status      string    `json:"status"`
}

type AttendanceWithNames struct {
    ID          int64     `json:"attendance_id"`
    USN        string    `json:"usn"`
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
	//MarkAttendance(attendance *AttendancePayload) (int64, error)
    BulkMarkAttendance(attendances []AttendancePayload) (int, error)
	GetAttendanceByStudentAndSubject(usn string, subjectCode string) ([]AttendanceWithNames, error)
	GetAttendanceBySubjectAndDate(subjectCode string, date time.Time) ([]AttendanceWithNames, error)
	AssignSubjectToTimeRange(facultyID int, subjectCode string, classDate time.Time, start time.Time, end time.Time) (int64, int64, error)
	GetAttendanceSummaryBySubject(subjectCode string) ([]StudentSummary, error)
	GetClassAttendance(subjectCode string, date time.Time) ([]ClassAttendance, error)
	GetStudentAttendanceHistory(usn string, subjectCode string) ([]StudentHistory, error)
    GetAttendanceSummaryByStudent(usn string) ([]SubjectSummary, error)
}