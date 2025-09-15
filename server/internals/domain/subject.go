package domain

type SubjectPayload struct {
	Code       string `json:"subject_code"`
	Name       string `json:"subject_name"`
	FacultyID  int64  `json:"faculty_id"`
	Department string `json:"department"`
	Sem        int    `json:"sem"`
}

type Subject struct {
	ID         int64  `json:"subject_id"`
	Code       string `json:"subject_code"`
	Name       string `json:"subject_name"`
	Faculty    string `json:"faculty"` // This returns faculty name (JOINed in repo)
	Department string `json:"department"`
	Sem        int    `json:"sem"`
}

// type StudentWithSubjects struct {
// 	ID       int64            `json:"id"`
// 	Username string           `json:"username"`
// 	USN      string           `json:"usn"`
// 	Subjects []Subject `json:"subjects"`
// }


type SubjectRepo interface {
	AddSubject(subject SubjectPayload) (int64, error)
	GetSubjectsByDeptAndSem(department string, sem int) ([]Subject, error)
	GetSubjectsByFacultyID(facultyID int) ([]Subject, error)
}
