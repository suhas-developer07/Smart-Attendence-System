package repository

import (
	"database/sql"

	"github.com/elastic/go-elasticsearch/v8/typedapi/sql/query"
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

func (p *PostgresRepo) InsertStudents(domain.StudentRegisterPayload)error{
  query := `INSERT INTO students (usn,username,department,sem) VALUES($1,$2,$3,$4)`

  _,err := 
}
