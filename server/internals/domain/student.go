package domain

import "time"

type Student struct {
	ID           int64   `json:"student_id"`
	USN          string  `json:"usn"`
	Username     string  `json:"username"`
	Department   string  `json:"department"`
	Sem          int     `json:"sem"`
	FaceEncoding []byte  `json:"face_encoding,omitempty"`
	NFCUID       *string `json:"nfc_uid,omitempty"`
}

type StudentRegisterPayload struct {
	USN        string `json:"usn"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Department string `json:"department"`
	Sem        int    `json:"sem"`
}

type StudentLoginPayload struct {
	USN string `json:"usn" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type GetStudentsWithSubjectsPayload struct {
	Department string `json:"department"`
	Sem        int    `json:"sem"`
}

type StudentUpdatePayload struct {
	Username   string `json:"username"`
	Department string `json:"department"`
	Sem        int    `json:"sem"`
}

type StudentSummary struct {
    USN         string   `json:"usn"`
    StudentName  string  `json:"student_name"`
    TotalClasses int     `json:"total_classes"`
    Attended     int     `json:"attended"`
    Percentage   float64 `json:"percentage"`
}

type StudentHistory struct {
    ID          int64     `json:"attendance_id"`
    Date        time.Time `json:"date"`
    Status      string    `json:"status"`
    SubjectID   int64     `json:"subject_id"`
    SubjectName string    `json:"subject_name"`
    RecordedAt  time.Time `json:"recorded_at"`
}

type StudentRepo interface {
	StudentRegister(student StudentRegisterPayload) (int64, error)
	UpdateStudentInfo(studentID int, payload StudentUpdatePayload) error
	// GetStudentsByDeptAndSem(department string, sem int) ([]Student, error)
	GetSubjectsByStudentID(studentID int64) ([]SubjectPayload, error)
	LoginStudent(usn, password string) (string, error)
}

