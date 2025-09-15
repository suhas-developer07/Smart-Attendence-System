package domain

type SubjectPayload struct{
	Code       string `json:"subject_code"`
	Name       string `json:"subject_name"`
	Faculty    string `json:"faculty"`
	Department string `json:"department"`
	Sem        int    `json:"sem"`
}

type Subject struct {
	Id         string `json:"subject_id"`
	Code       string `json:"subject_code"`
	Name       string `json:"subject_name"`
	Faculty    string `json:"faculty"`
	Department string `json:"department"`
	Sem        int    `json:"sem"`
}

type StudentWithSubjects struct {
	ID       int64            `json:"id"`
	Username string           `json:"username"`
	USN      string           `json:"usn"`
	Subjects []Subject `json:"subjects"`
}



type SubjectRepo interface {
	AddSubject(subject SubjectPayload) (int64, error)
	ListSubjects(branch string, sem int) ([]Subject, error)
	//GetStudentsWithSubjects(req GetStudentsWithSubjectsPayload) ([]StudentWithSubjects, error)
}
