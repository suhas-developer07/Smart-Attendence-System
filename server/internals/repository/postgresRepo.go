package repository

import (
	"database/sql"
	"fmt"

	"github.com/suhas-developer07/Smart-Attendence-System/server/internals/domain"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (p *PostgresRepo) InitTables() error {
	queries := []string{
		`
	CREATE TABLE IF NOT EXISTS students (
	    student_id SERIAL PRIMARY KEY,
	    usn VARCHAR(20) UNIQUE NOT NULL,
	    username VARCHAR(100) NOT NULL,
	    department VARCHAR(50) NOT NULL,
	    sem INT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS subjects (
	    subject_id SERIAL PRIMARY KEY,
	    subject_code VARCHAR(20) UNIQUE NOT NULL,
	    subject_name VARCHAR(100) NOT NULL,
	    faculty VARCHAR(100),
	    department VARCHAR(50) NOT NULL,
	    sem INT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS student_subjects (
	    student_id INT NOT NULL,
	    subject_id INT NOT NULL,
	    PRIMARY KEY (student_id, subject_id),
	    CONSTRAINT fk_student FOREIGN KEY (student_id) REFERENCES students(student_id) ON DELETE CASCADE,
	    CONSTRAINT fk_subject FOREIGN KEY (subject_id) REFERENCES subjects(subject_id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS attendance (
	    attendance_id SERIAL PRIMARY KEY,
	    student_id INT NOT NULL,
	    subject_id INT NOT NULL,
	    date DATE NOT NULL,
	    status VARCHAR(10) NOT NULL CHECK (status IN ('Present', 'Absent')),
	    CONSTRAINT fk_student FOREIGN KEY (student_id) REFERENCES students(student_id),
	    CONSTRAINT fk_subject FOREIGN KEY (subject_id) REFERENCES subjects(subject_id),
	    UNIQUE(student_id, subject_id, date)
	);
	`,
	}

	for _, query := range queries {
		if _, err := p.db.Exec(query); err != nil {
			return err
		}
	}
	return nil
}

func (p *PostgresRepo) StudentRegister(student domain.StudentRegisterPayload) (int64, error) {
	var id int64
	query := `INSERT INTO students (usn, username, department, sem) VALUES($1, $2, $3, $4) RETURNING student_id;`

	err := p.db.QueryRow(query, student.USN, student.Username, student.Branch, student.Sem).Scan(&id)
	if err != nil {
		return 0, err
	}

	// Auto-assign already existing subjects of that dept+sem to this new student
	assignQuery := `
	INSERT INTO student_subjects (student_id, subject_id)
	SELECT $1, sub.subject_id
	FROM subjects sub
	WHERE sub.department = $2 AND sub.sem = $3
	ON CONFLICT DO NOTHING;`

	_, err = p.db.Exec(assignQuery, id, student.Branch, student.Sem)
	if err != nil {
		return id, fmt.Errorf("failed to auto-assign subjects: %w", err)
	}

	return id, nil
}

func (p *PostgresRepo) AddSubject(subject domain.SubjectPayload) (int64, error) {
	var id int64
	query := `INSERT INTO subjects (subject_code, subject_name, faculty, department, sem)
	          VALUES($1, $2, $3, $4, $5) RETURNING subject_id;`

	err := p.db.QueryRow(query,
		subject.Code, subject.Name, subject.Faculty, subject.Department, subject.Sem).Scan(&id)
	if err != nil {
		return 0, err	
	}	

	// Auto-assign this subject to all existing students of that dept+sem
	assignQuery := `
	INSERT INTO student_subjects (student_id, subject_id)
	SELECT s.student_id, $1
	FROM students s
	WHERE s.department = $2 AND s.sem = $3
	ON CONFLICT DO NOTHING;`

	_, err = p.db.Exec(assignQuery, id, subject.Department, subject.Sem)
	if err != nil {
		return id, fmt.Errorf("failed to assign subject to students: %w", err)
	}

	return id, nil
}

func (p *PostgresRepo) GetStudentsWithSubjects(department string, sem int) ([]domain.StudentWithSubjects, error) {
	rows, err := p.db.Query(`
		SELECT s.student_id, s.username, s.usn,
		       sub.subject_id, sub.subject_code, sub.subject_name, sub.faculty
		FROM students s
		JOIN student_subjects ss ON s.student_id = ss.student_id
		JOIN subjects sub ON ss.subject_id = sub.subject_id
		WHERE s.department = $1 AND s.sem = $2
		ORDER BY s.student_id;`, department, sem)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Map student -> subjects
	studentMap := make(map[int64]*domain.StudentWithSubjects)

	for rows.Next() {
		var sid int64
		var username, usn string
		var subject domain.Subject

		err := rows.Scan(&sid, &username, &usn,
			&subject.Id, &subject.Code, &subject.Name, &subject.Faculty)
		if err != nil {
			return nil, err
		}

		if _, exists := studentMap[sid]; !exists {
			studentMap[sid] = &domain.StudentWithSubjects{
				ID:       sid,
				Username: username,
				USN:      usn,
				Subjects: []domain.Subject{},
			}
		}
		studentMap[sid].Subjects = append(studentMap[sid].Subjects, subject)
	}

	// Convert map -> slice
	students := make([]domain.StudentWithSubjects, 0, len(studentMap))
	for _, s := range studentMap {
		students = append(students, *s)
	}
	return students, nil
}
