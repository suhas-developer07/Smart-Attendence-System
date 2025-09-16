package domain

type Attendance struct {
	ID string `json:"attendance_id"`
	StudentID   string `json:"student_id"`
	SubjectID   string `json:"subject_id"`
	Date         string `json:"date"`
	Status       string `json:"status"`
	RecordedAt string `json:"recorded_at"`
	CreatedAt string `json:"created_at"`
}