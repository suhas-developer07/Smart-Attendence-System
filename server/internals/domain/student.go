package domain
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
	Department string `json:"department"`
	Sem        int    `json:"sem"`
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

type StudentRepo interface {
	StudentRegister(student StudentRegisterPayload) (int64, error)
	UpdateStudentInfo(studentID int, payload StudentUpdatePayload) error
	GetStudentsByDeptAndSem(department string, sem int) ([]Student, error)
	GetSubjectsByStudentID(studentID int64) ([]SubjectPayload, error)
}

type Repository interface {
    StudentRepo
    SubjectRepo
    FacultyRepo
}

