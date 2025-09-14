package domain

type Student struct {
	Student_id string `json:"student_id"`
	USN        string `json:"usn"`
	Username   string `json:"username"`
	Branch     string `json:"branch"`
	Sem        string `json:"sem"`
}

type StudentRegisterPayload struct {
	USN      string `json:"usn"`
	Username string `json:"username"`
	Branch   string `json:"branch"`
	Sem      string `json:"sem"`
}

type StudentRepo interface {
	StudentRegister(student StudentRegisterPayload) (int64, error)
}
