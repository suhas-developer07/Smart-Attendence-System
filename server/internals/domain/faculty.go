package domain

type Faculty struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Department string `json:"department"`
	CreatedAt  string `json:"created_at"`
}

type FacultyRegisterPayload struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Department string `json:"department"`
}

type FacultyLoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type FacultyRepo interface {
	GetFacultyByID(facultyID int) (Faculty, error)
	CreateFaculty(req FacultyRegisterPayload) (int64, error)
	AuthenticateFaculty(req FacultyLoginPayload) (int64, error)
}
