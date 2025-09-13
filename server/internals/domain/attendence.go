package domain

type Attendence struct {
	Attendence_id string `json:"attendence_id"`
	Student_id   string `json:"student_id"`
	Subject_id   string `json:"subject_id"`
	Date         string `json:"date"`
	Status       string `json:"status"`
}