package repository

import (
	"database/sql"

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
	    department VARCHAR(50),
	    sem INT
	);

	CREATE TABLE IF NOT EXISTS subjects (
	    subject_id SERIAL PRIMARY KEY,
	    subject_code VARCHAR(20) UNIQUE NOT NULL,
	    subject_name VARCHAR(100) NOT NULL,
	    faculty VARCHAR(100)
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
	return id, nil
}
