package repository

import (
	"database/sql"
	"fmt"

	"github.com/suhas-developer07/Smart-Attendence-System/server/internals/domain"
	"golang.org/x/crypto/bcrypt"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// InitTables creates tables in an order that satisfies FK constraints.
func (p *PostgresRepo) InitTables() error {
	queries := []string{
		// 1. faculty first (subjects references faculty)
		`CREATE TABLE IF NOT EXISTS faculty (
		    faculty_id SERIAL PRIMARY KEY,
		    faculty_name VARCHAR(100) NOT NULL,
		    email VARCHAR(100) UNIQUE NOT NULL,
		    password_hash VARCHAR(256) NOT NULL,
		    department VARCHAR(50) NOT NULL,
		    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
		);`,

		// 2. students (add optional face encodings and nfc uid)
		`CREATE TABLE IF NOT EXISTS students (
		    student_id SERIAL PRIMARY KEY,
		    usn VARCHAR(20) UNIQUE NOT NULL,
		    username VARCHAR(100) NOT NULL,
		    department VARCHAR(50) NOT NULL,
		    sem INT NOT NULL,
		    face_encoding BYTEA NULL,
		    nfc_uid VARCHAR(100) UNIQUE NULL
		);`,

		// 3. subjects (references faculty)
		`CREATE TABLE IF NOT EXISTS subjects (
		    subject_id SERIAL PRIMARY KEY,
		    subject_code VARCHAR(20) UNIQUE NOT NULL,
		    subject_name VARCHAR(100) NOT NULL,
		    department VARCHAR(50) NOT NULL,
		    sem INT NOT NULL,
		    faculty_id INT NOT NULL,
		    CONSTRAINT fk_faculty FOREIGN KEY (faculty_id) REFERENCES faculty(faculty_id) ON DELETE RESTRICT
		);`,

		// 4. mapping table
		`CREATE TABLE IF NOT EXISTS student_subjects (
		    student_id INT NOT NULL,
		    subject_id INT NOT NULL,
		    PRIMARY KEY (student_id, subject_id),
		    CONSTRAINT fk_student FOREIGN KEY (student_id) REFERENCES students(student_id) ON DELETE CASCADE,
		    CONSTRAINT fk_subject FOREIGN KEY (subject_id) REFERENCES subjects(subject_id) ON DELETE CASCADE
		);`,

		// 5. attendance table
		`CREATE TABLE IF NOT EXISTS attendance (
		    attendance_id SERIAL PRIMARY KEY,
		    student_id INT NOT NULL,
		    subject_id INT NOT NULL,
		    date DATE NOT NULL,
		    status VARCHAR(10) NOT NULL CHECK (status IN ('Present', 'Absent')),
		    CONSTRAINT fk_student_att FOREIGN KEY (student_id) REFERENCES students(student_id),
		    CONSTRAINT fk_subject_att FOREIGN KEY (subject_id) REFERENCES subjects(subject_id),
		    UNIQUE(student_id, subject_id, date)
		);`,
	}

	for _, query := range queries {
		if _, err := p.db.Exec(query); err != nil {
			return fmt.Errorf("failed to exec init query: %w", err)
		}
	}
	return nil
}

// StudentRegister inserts a new student and auto-assigns existing subjects for that dept+sem.
func (p *PostgresRepo) StudentRegister(student domain.StudentRegisterPayload) (int64, error) {
	// Wrap in a transaction to ensure consistency
	tx, err := p.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer func() {
		// if tx not committed, rollback silently
		_ = tx.Rollback()
	}()

	var id int64
	query := `INSERT INTO students (usn, username, department, sem, face_encoding, nfc_uid) VALUES($1, $2, $3, $4, $5, $6) RETURNING student_id;`

	err = tx.QueryRow(query, student.USN, student.Username, student.Department, student.Sem).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert student: %w", err)
	}

	assignQuery := `
	INSERT INTO student_subjects (student_id, subject_id)
	SELECT $1, sub.subject_id
	FROM subjects sub
	WHERE sub.department = $2 AND sub.sem = $3
	ON CONFLICT DO NOTHING;`

	if _, err := tx.Exec(assignQuery, id, student.Department, student.Sem); err != nil {
		return id, fmt.Errorf("failed to auto-assign subjects: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return id, fmt.Errorf("failed to commit tx: %w", err)
	}

	return id, nil
}

// AddSubject inserts a subject and auto-assigns the subject to existing students of same dept+sem
func (p *PostgresRepo) AddSubject(subject domain.SubjectPayload) (int64, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	var id int64
	query := `INSERT INTO subjects (subject_code, subject_name, faculty_id, department, sem)
	          VALUES($1, $2, $3, $4, $5) RETURNING subject_id;`

	err = tx.QueryRow(query,
		subject.Code, subject.Name, subject.FacultyID, subject.Department, subject.Sem).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert subject: %w", err)
	}

	assignQuery := `
	INSERT INTO student_subjects (student_id, subject_id)
	SELECT s.student_id, $1
	FROM students s
	WHERE s.department = $2 AND s.sem = $3
	ON CONFLICT DO NOTHING;`

	if _, err := tx.Exec(assignQuery, id, subject.Department, subject.Sem); err != nil {
		return id, fmt.Errorf("failed to assign subject to students: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return id, fmt.Errorf("failed to commit tx: %w", err)
	}

	return id, nil
}

// GetSubjectsByDeptAndSem returns subjects for a department+sem and includes faculty name via join
func (p *PostgresRepo) GetSubjectsByDeptAndSem(department string, sem int) ([]domain.Subject, error) {
	query := `SELECT s.subject_id, s.subject_code, s.subject_name, s.department, s.sem, f.faculty_name
	          FROM subjects s
	          JOIN faculty f ON s.faculty_id = f.faculty_id
	          WHERE s.department = $1 AND s.sem = $2;`

	rows, err := p.db.Query(query, department, sem)
	if err != nil {
		return nil, fmt.Errorf("failed to query subjects: %w", err)
	}
	defer rows.Close()

	var subjects []domain.Subject
	for rows.Next() {
		var sub domain.Subject
		if err := rows.Scan(&sub.ID, &sub.Code, &sub.Name, &sub.Department, &sub.Sem, &sub.Faculty); err != nil {
			return nil, fmt.Errorf("failed to scan subject: %w", err)
		}
		subjects = append(subjects, sub)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return subjects, nil
}

func (p *PostgresRepo) UpdateStudentInfo(studentID int, payload domain.StudentUpdatePayload) error {
	query := `
	UPDATE students
	SET username = $2, department = $3, sem = $4
	WHERE student_id = $1;`

	_, err := p.db.Exec(query, studentID, payload.Username, payload.Department, payload.Sem)
	if err != nil {
		return fmt.Errorf("failed to update student info: %w", err)
	}
	return nil
}

func (p *PostgresRepo) GetStudentsByDeptAndSem(department string, sem int) ([]domain.Student, error) {
	query := `SELECT student_id, usn, username, department, sem
	                          FROM students
	                          WHERE department = $1 AND sem = $2;`
	rows, err := p.db.Query(query, department, sem)
	if err != nil {
		return nil, fmt.Errorf("failed to query students: %w", err)
	}
	defer rows.Close()

	var students []domain.Student
	for rows.Next() {
		var stu domain.Student
		if err := rows.Scan(&stu.ID, &stu.USN, &stu.Username, &stu.Department, &stu.Sem); err != nil {
			return nil, fmt.Errorf("failed to scan student: %w", err)
		}
		students = append(students, stu)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return students, nil
}

// GetSubjectsByStudentID returns subject payloads assigned to a student. Returns faculty_id (int) for each subject.
func (p *PostgresRepo) GetSubjectsByStudentID(studentID int64) ([]domain.SubjectPayload, error) {
	var subjects []domain.SubjectPayload
	query := `
	SELECT sub.subject_code, sub.subject_name, sub.faculty_id, sub.department, sub.sem
	FROM subjects sub
	JOIN student_subjects ss ON sub.subject_id = ss.subject_id
	WHERE ss.student_id = $1;`

	rows, err := p.db.Query(query, studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subjects for student: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var subject domain.SubjectPayload
		if err := rows.Scan(&subject.Code, &subject.Name, &subject.FacultyID, &subject.Department, &subject.Sem); err != nil {
			return nil, fmt.Errorf("failed to scan subject: %w", err)
		}
		subjects = append(subjects, subject)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return subjects, nil
}

func (p *PostgresRepo) GetSubjectsByFacultyID(facultyID int) ([]domain.Subject, error) {
	var subjects []domain.Subject
	query := `
    SELECT
        s.subject_id, s.subject_code, s.subject_name, s.department, s.sem, f.faculty_name
    FROM subjects s
    JOIN faculty f ON s.faculty_id = f.faculty_id
    WHERE s.faculty_id = $1;`

	rows, err := p.db.Query(query, facultyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subjects for faculty ID %d: %w", facultyID, err)
	}
	defer rows.Close()

	for rows.Next() {
		var subject domain.Subject
		if err := rows.Scan(
			&subject.ID,
			&subject.Code,
			&subject.Name,
			&subject.Department,
			&subject.Sem,
			&subject.Faculty,
		); err != nil {
			return nil, fmt.Errorf("failed to scan subject row: %w", err)
		}
		subjects = append(subjects, subject)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return subjects, nil
}

// --- Faculty registration & auth helpers ---

func hashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func comparePassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// RegisterFaculty saves a new faculty user (hashes password before storing)
func (p *PostgresRepo) RegisterFaculty(name, email, password, department string) (int64, error) {
	hash, err := hashPassword(password)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	var id int64
	query := `INSERT INTO faculty (faculty_name, email, password_hash, department) VALUES($1, $2, $3, $4) RETURNING faculty_id;`
	if err := p.db.QueryRow(query, name, email, hash, department).Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to insert faculty: %w", err)
	}
	return id, nil
}

// AuthenticateFaculty checks email & password, returns faculty_id on success.
func (p *PostgresRepo) AuthenticateFaculty(email, password string) (int64, error) {
	var id int64
	var pwHash string
	query := `SELECT faculty_id, password_hash FROM faculty WHERE email = $1;`
	if err := p.db.QueryRow(query, email).Scan(&id, &pwHash); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("faculty not found")
		}
		return 0, fmt.Errorf("failed to query faculty: %w", err)
	}
	if err := comparePassword(pwHash, password); err != nil {
		return 0, fmt.Errorf("invalid credentials")
	}
	return id, nil
}

// UpdateStudentFaceEncoding stores a binary face encoding for facial recognition
func (p *PostgresRepo) UpdateStudentFaceEncoding(studentID int64, encoding []byte) error {
	query := `UPDATE students SET face_encoding = $2 WHERE student_id = $1;`
	if _, err := p.db.Exec(query, studentID, encoding); err != nil {
		return fmt.Errorf("failed to update face encoding: %w", err)
	}
	return nil
}

// UpdateStudentNFC stores/updates a student's NFC UID
func (p *PostgresRepo) UpdateStudentNFC(studentID int64, uid string) error {
	query := `UPDATE students SET nfc_uid = $2 WHERE student_id = $1;`
	if _, err := p.db.Exec(query, studentID, uid); err != nil {
		return fmt.Errorf("failed to update nfc uid: %w", err)
	}
	return nil
}

// Note: Attendance marking / reports (core attendance logic) intentionally omitted per user's request.
// We'll implement those in the next step.

// Small utility: GetFacultyByID (useful for admin/faculty screens)
func (p *PostgresRepo) GetFacultyByID(facultyID int) (domain.Faculty, error) {
	var f domain.Faculty
	query := `SELECT faculty_id, faculty_name, email, department, created_at FROM faculty WHERE faculty_id = $1;`
	if err := p.db.QueryRow(query, facultyID).Scan(&f.ID, &f.Name, &f.Email, &f.Department, &f.CreatedAt); err != nil {
		return f, fmt.Errorf("failed to get faculty: %w", err)
	}
	return f, nil
}
