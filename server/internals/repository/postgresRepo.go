package repository

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/suhas-developer07/Smart-Attendence-System/server/internals/domain"
	"github.com/suhas-developer07/Smart-Attendence-System/server/pkg/utils"
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
		// 1. admins
		`CREATE TABLE IF NOT EXISTS admins (
			admin_id SERIAL PRIMARY KEY,
			username VARCHAR(100) NOT NULL UNIQUE,
			email VARCHAR(100) UNIQUE NOT NULL,
			password_hash VARCHAR(256) NOT NULL,
			created_at TIMESTAMPTZ DEFAULT now()
		);`,

		// 2. faculty
		`CREATE TABLE IF NOT EXISTS faculty (
			faculty_id SERIAL PRIMARY KEY,
			faculty_name VARCHAR(100) NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			password_hash VARCHAR(256) NOT NULL,
			department VARCHAR(50) NOT NULL,
			created_at TIMESTAMPTZ DEFAULT now()
		);`,

		// 3. faculty api keys (admin issues keys to faculty for access)
		`CREATE TABLE IF NOT EXISTS faculty_api_keys (
			key_id SERIAL PRIMARY KEY,
			faculty_id INT NOT NULL,
			api_key VARCHAR(256) NOT NULL UNIQUE,
			created_at TIMESTAMPTZ DEFAULT now(),
			expires_at TIMESTAMPTZ NULL,
			CONSTRAINT fk_faculty_key FOREIGN KEY (faculty_id) REFERENCES faculty(faculty_id) ON DELETE CASCADE
		);`,

		// 4. students
		`CREATE TABLE IF NOT EXISTS students (
			student_id SERIAL PRIMARY KEY,
			usn VARCHAR(50) UNIQUE NOT NULL,
			username VARCHAR(100) NOT NULL,
			department VARCHAR(50) NOT NULL,
			sem INT NOT NULL,
			face_encoding BYTEA NULL,
			nfc_uid VARCHAR(100) UNIQUE NULL,
			created_at TIMESTAMPTZ DEFAULT now()
		);`,

		// 5. subjects
		`CREATE TABLE IF NOT EXISTS subjects (
			subject_id SERIAL PRIMARY KEY,
			subject_code VARCHAR(50) UNIQUE NOT NULL,
			subject_name VARCHAR(150) NOT NULL,
			department VARCHAR(50) NOT NULL,
			sem INT NOT NULL,
			faculty_id INT NOT NULL,
			CONSTRAINT fk_faculty_sub FOREIGN KEY (faculty_id) REFERENCES faculty(faculty_id) ON DELETE RESTRICT
		);`,

		// 6. student_subjects mapping
		`CREATE TABLE IF NOT EXISTS student_subjects (
			student_id INT NOT NULL,
			subject_id INT NOT NULL,
			PRIMARY KEY (student_id, subject_id),
			CONSTRAINT fk_student_sub FOREIGN KEY (student_id) REFERENCES students(student_id) ON DELETE CASCADE,
			CONSTRAINT fk_subject_sub FOREIGN KEY (subject_id) REFERENCES subjects(subject_id) ON DELETE CASCADE
		);`,

		// 7. attendance
     	`CREATE TABLE IF NOT EXISTS attendance (
			attendance_id SERIAL PRIMARY KEY,
			student_id INT NOT NULL,
			subject_id INT  NULL,
			date DATE NOT NULL,
			status VARCHAR(20) NOT NULL CHECK (status IN ('Present','Absent','Late')),
			recorded_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			created_at TIMESTAMPTZ DEFAULT now(),
			CONSTRAINT fk_attendance_student FOREIGN KEY (student_id) REFERENCES students(student_id) ON DELETE CASCADE,
			CONSTRAINT fk_attendance_subject FOREIGN KEY (subject_id) REFERENCES subjects(subject_id) ON DELETE CASCADE,
			UNIQUE(student_id, subject_id, date)
);`,

	}

	for _, q := range queries {
		if _, err := p.db.Exec(q); err != nil {
			return fmt.Errorf("failed to exec init query: %w", err)
		}
	}
	return nil
}

// ------------------------ Student & Subject Operations ------------------------

func (p *PostgresRepo) StudentRegister(student domain.StudentRegisterPayload) (int64, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	fmt.Println("Department:",student.Department)

	var id int64
	query := `INSERT INTO students (usn, username, department, sem)
	          VALUES ($1, $2, $3, $4) RETURNING student_id;`
	err = tx.QueryRow(query, student.USN, student.Username, student.Department, student.Sem).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert student: %w", err)
	}

	// Auto assign existing subjects for dept+sem
	assignQuery := `
	INSERT INTO student_subjects (student_id, subject_id)
	SELECT $1, s.subject_id FROM subjects s WHERE s.department = $2 AND s.sem = $3
	ON CONFLICT DO NOTHING;`
	if _, err := tx.Exec(assignQuery, id, student.Department, student.Sem); err != nil {
		return id, fmt.Errorf("auto-assign subjects: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return id, fmt.Errorf("commit tx: %w", err)
	}
	return id, nil
}

func (p *PostgresRepo) UpdateStudentInfo(studentID int, payload domain.StudentUpdatePayload) error {
	query := `UPDATE students SET username = $2, department = $3, sem = $4 WHERE student_id = $1;`
	if _, err := p.db.Exec(query, studentID, payload.Username, payload.Department, payload.Sem); err != nil {
		return fmt.Errorf("update student: %w", err)
	}
	return nil
}

func (p *PostgresRepo) UpdateStudentFaceEncoding(studentID int64, encoding []byte) error {
	query := `UPDATE students SET face_encoding = $2 WHERE student_id = $1;`
	if _, err := p.db.Exec(query, studentID, encoding); err != nil {
		return fmt.Errorf("update face encoding: %w", err)
	}
	return nil
}

func (p *PostgresRepo) UpdateStudentNFC(studentID int64, uid string) error {
	query := `UPDATE students SET nfc_uid = $2 WHERE student_id = $1;`
	if _, err := p.db.Exec(query, studentID, uid); err != nil {
		return fmt.Errorf("update nfc uid: %w", err)
	}
	return nil
}

func (p *PostgresRepo) AddSubject(subject domain.SubjectPayload) (int64, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	var id int64
	query := `INSERT INTO subjects (subject_code, subject_name, faculty_id, department, sem)
	          VALUES ($1, $2, $3, $4, $5) RETURNING subject_id;`
	if err := tx.QueryRow(query, subject.Code, subject.Name, subject.FacultyID, subject.Department, subject.Sem).Scan(&id); err != nil {
		return 0, fmt.Errorf("insert subject: %w", err)
	}

	assignQuery := `
	INSERT INTO student_subjects (student_id, subject_id)
	SELECT s.student_id, $1 FROM students s WHERE s.department = $2 AND s.sem = $3
	ON CONFLICT DO NOTHING;`
	if _, err := tx.Exec(assignQuery, id, subject.Department, subject.Sem); err != nil {
		return id, fmt.Errorf("assign subject to students: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return id, fmt.Errorf("commit tx: %w", err)
	}
	return id, nil
}

// ------------------------ Subjects queries ------------------------
//subjects of a particular department and sem
func (p *PostgresRepo) GetSubjectsByDeptAndSem(department string, sem int) ([]domain.Subject, error) {
	q := `SELECT s.subject_id, s.subject_code, s.subject_name, s.department, s.sem, f.faculty_name
	      FROM subjects s JOIN faculty f ON s.faculty_id = f.faculty_id
	      WHERE s.department = $1 AND s.sem = $2;`
	rows, err := p.db.Query(q, department, sem)
	if err != nil {
		return nil, fmt.Errorf("query subjects: %w", err)
	}
	defer rows.Close()

	var list []domain.Subject
	for rows.Next() {
		var s domain.Subject
		if err := rows.Scan(&s.ID, &s.Code, &s.Name, &s.Department, &s.Sem, &s.Faculty); err != nil {
			return nil, fmt.Errorf("scan subject: %w", err)
		}
		list = append(list, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}
	return list, nil
}

func (p *PostgresRepo) GetSubjectsByStudentID(studentID int64) ([]domain.SubjectPayload, error) {
	q := `SELECT sub.subject_code, sub.subject_name, sub.faculty_id, sub.department, sub.sem
	      FROM subjects sub JOIN student_subjects ss ON sub.subject_id = ss.subject_id
	      WHERE ss.student_id = $1;`
	rows, err := p.db.Query(q, studentID)
	if err != nil {
		return nil, fmt.Errorf("get subjects for student: %w", err)
	}
	defer rows.Close()

	var out []domain.SubjectPayload
	for rows.Next() {
		var sp domain.SubjectPayload
		if err := rows.Scan(&sp.Code, &sp.Name, &sp.FacultyID, &sp.Department, &sp.Sem); err != nil {
			return nil, fmt.Errorf("scan subj payload: %w", err)
		}
		out = append(out, sp)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}
	return out, nil
}
//subjects handled by faculty
func (p *PostgresRepo) GetSubjectsByFacultyID(facultyID int) ([]domain.Subject, error) {
	q := `SELECT s.subject_id, s.subject_code, s.subject_name, s.department, s.sem, f.faculty_name
	      FROM subjects s JOIN faculty f ON s.faculty_id = f.faculty_id
	      WHERE s.faculty_id = $1;`
	rows, err := p.db.Query(q, facultyID)
	if err != nil {
		return nil, fmt.Errorf("query subjects faculty: %w", err)
	}
	defer rows.Close()

	var list []domain.Subject
	for rows.Next() {
		var s domain.Subject
		if err := rows.Scan(&s.ID, &s.Code, &s.Name, &s.Department, &s.Sem, &s.Faculty); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		list = append(list, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}
	return list, nil
}

// ------------------------ Faculty & Admin ------------------------

func (p *PostgresRepo) CreateFaculty(req domain.FacultyRegisterPayload) (int64, error) {
	pwHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return 0, fmt.Errorf("hash password: %w", err)
	}
	var id int64
	q := `INSERT INTO faculty (faculty_name, email, password_hash, department) VALUES ($1, $2, $3, $4) RETURNING faculty_id;`
	if err := p.db.QueryRow(q, req.Name, req.Email, pwHash, req.Department).Scan(&id); err != nil {
		return 0, fmt.Errorf("insert faculty: %w", err)
	}
	return id, nil
}

func (p *PostgresRepo) AuthenticateFaculty(req domain.FacultyLoginPayload) (int64, error) {
	var id int64
	var pwHash string
	q := `SELECT faculty_id, password_hash FROM faculty WHERE email = $1;`
	if err := p.db.QueryRow(q, req.Email).Scan(&id, &pwHash); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("faculty not found")
		}
		return 0, fmt.Errorf("query faculty: %w", err)
	}
	if err := utils.ComparePassword(pwHash, req.Password); err != nil {
		return 0, fmt.Errorf("invalid credentials")
	}
	return id, nil
}

// Get all faculty
func (p *PostgresRepo) GetAllFaculty() ([]domain.Faculty, error) {
    query := `SELECT faculty_id, faculty_name, email, department FROM faculty`

    rows, err := p.db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("get faculty: %w", err)
    }
    defer rows.Close()

    var facultyList []domain.Faculty
    for rows.Next() {
        var f domain.Faculty
        if err := rows.Scan(&f.ID, &f.Name, &f.Email, &f.Department); err != nil {
            return nil, err
        }
        facultyList = append(facultyList, f)
    }
    return facultyList, nil
}

func (p *PostgresRepo) GetFacultyByDepartment(department string) ([]domain.Faculty, error) {
    query := `SELECT faculty_id, faculty_name, email, department FROM faculty WHERE department = $1`

    rows, err := p.db.Query(query, department)
    if err != nil {
        return nil, fmt.Errorf("get faculty by department: %w", err)
    }
    defer rows.Close()

    var facultyList []domain.Faculty
    for rows.Next() {
        var f domain.Faculty
        if err := rows.Scan(&f.ID, &f.Name, &f.Email, &f.Department); err != nil {
            return nil, err
        }
        facultyList = append(facultyList, f)
    }
    return facultyList, nil
}

func (p *PostgresRepo) GetFacultyByID(facultyID int) (domain.Faculty, error) {
	var f domain.Faculty
	q := `SELECT faculty_id, faculty_name, email, department, created_at FROM faculty WHERE faculty_id = $1;`
	if err := p.db.QueryRow(q, facultyID).Scan(&f.ID, &f.Name, &f.Email, &f.Department, &f.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return domain.Faculty{}, fmt.Errorf("faculty not found")
		}
		return domain.Faculty{}, fmt.Errorf("query faculty: %w", err)
	}
	return f, nil
}


func (p *PostgresRepo) CreateAdmin(username, email, password string) (int64, error) {
	pwHash, err := utils.HashPassword(password)
	if err != nil {
		return 0, fmt.Errorf("hash password: %w", err)
	}
	var id int64
	q := `INSERT INTO admins (username, email, password_hash) VALUES ($1, $2, $3) RETURNING admin_id;`
	if err := p.db.QueryRow(q, username, email, pwHash).Scan(&id); err != nil {
		return 0, fmt.Errorf("insert admin: %w", err)
	}
	return id, nil
}

// CreateFacultyAPIKey creates an API key (random string) for a faculty (admin action).
// utils.GenerateRandomKey() is assumed to exist — replace with your own generator.
func (p *PostgresRepo) CreateFacultyAPIKey(facultyID int, expiresAt *time.Time) (string, error) {
	apiKey, err := utils.GenerateRandomKey(48)
	if err != nil {
		return "", fmt.Errorf("generate key: %w", err)
	}
	q := `INSERT INTO faculty_api_keys (faculty_id, api_key, expires_at) VALUES ($1, $2, $3);`
	if _, err := p.db.Exec(q, facultyID, apiKey, expiresAt); err != nil {
		return "", fmt.Errorf("insert api key: %w", err)
	}
	return apiKey, nil
}

func (p *PostgresRepo) ValidateFacultyAPIKey(apiKey string) (int, error) {
	var facultyID int
	var expires sql.NullTime
	q := `SELECT faculty_id, expires_at FROM faculty_api_keys WHERE api_key = $1;`
	if err := p.db.QueryRow(q, apiKey).Scan(&facultyID, &expires); err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("invalid api key")
		}
		return 0, fmt.Errorf("query api key: %w", err)
	}
	if expires.Valid && time.Now().After(expires.Time) {
		return 0, errors.New("api key expired")
	}
	return facultyID, nil
}

// ------------------------ Attendance ------------------------

func (p *PostgresRepo) MarkAttendance(req *domain.AttendancePayload) (int64, error) {
	var id int64
	q := `INSERT INTO attendance (student_id, subject_id, date, status, recorded_at)
	      VALUES ($1, $2, $3, $4, $5)
	      ON CONFLICT (student_id, subject_id, date)
	      DO UPDATE SET status = EXCLUDED.status, recorded_at = EXCLUDED.recorded_at
	      RETURNING attendance_id;`

	err := p.db.QueryRow(q, req.StudentID, req.SubjectID, req.RecordedAt.Format("2006-01-02"), req.Status, req.RecordedAt).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("mark attendance: %w", err)
	}
	return id, nil
}

func (p *PostgresRepo) GetAttendanceByStudentAndSubject(studentID, subjectID int64) ([]domain.Attendance, error) {
	q := `SELECT attendance_id, student_id, subject_id, date, status, recorded_at, created_at
	      FROM attendance
	      WHERE student_id = $1 AND subject_id = $2
	      ORDER BY date ASC;`

	rows, err := p.db.Query(q, studentID, subjectID)
	if err != nil {
		return nil, fmt.Errorf("query attendance: %w", err)
	}
	defer rows.Close()

	var list []domain.Attendance
	for rows.Next() {
		var a domain.Attendance
		if err := rows.Scan(&a.ID, &a.StudentID, &a.SubjectID, &a.Date, &a.Status, &a.RecordedAt, &a.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan attendance: %w", err)
		}
		list = append(list, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}
	return list, nil
}

func (p *PostgresRepo) GetAttendanceBySubject(subjectID int64, fromDate, toDate time.Time) ([]domain.Attendance, error) {
	q := `SELECT attendance_id, student_id, subject_id, date, status, recorded_at, created_at
	      FROM attendance
	      WHERE subject_id = $1 AND date BETWEEN $2 AND $3
	      ORDER BY date ASC;`

	rows, err := p.db.Query(q, subjectID, fromDate, toDate)
	if err != nil {
		return nil, fmt.Errorf("query attendance by subject: %w", err)
	}
	defer rows.Close()

	var list []domain.Attendance
	for rows.Next() {
		var a domain.Attendance
		if err := rows.Scan(&a.ID, &a.StudentID, &a.SubjectID, &a.Date, &a.Status, &a.RecordedAt, &a.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan attendance: %w", err)
		}
		list = append(list, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}
	return list, nil
}


// AssignSubjectToTimeRange assigns subjectID to all attendance rows in the given time range
// that currently have subject_id IS NULL. The caller must be the faculty that owns the subject.
// Returns (numberUpdated, numberSkippedDueToExistingRecord, error).
func (p *PostgresRepo) AssignSubjectToTimeRange(facultyID int, subjectID int64, start, end time.Time) (int64, int64, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return 0, 0, fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	// 1) Verify the subject belongs to the faculty (authorization)
	var ownerID int
	if err := tx.QueryRow(`SELECT faculty_id FROM subjects WHERE subject_id = $1 FOR UPDATE;`, subjectID).Scan(&ownerID); err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, errors.New("subject not found")
		}
		return 0, 0, fmt.Errorf("query subject owner: %w", err)
	}
	if ownerID != facultyID {
		return 0, 0, errors.New("not authorized to assign this subject")
	}

	// 2) Count candidate NULL-subject rows in the time range
	var totalCandidates int64
	if err := tx.QueryRow(
		`SELECT COUNT(*) FROM attendance WHERE subject_id IS NULL AND recorded_at >= $1 AND recorded_at <= $2;`,
		start, end,
	).Scan(&totalCandidates); err != nil {
		return 0, 0, fmt.Errorf("count candidates: %w", err)
	}

	if totalCandidates == 0 {
		// nothing to do
		if err := tx.Commit(); err != nil {
			return 0, 0, fmt.Errorf("commit tx: %w", err)
		}
		return 0, 0, nil
	}

	// 3) Update only those NULL rows for which there is NO existing attendance
	//    row for (same student, same date, this subject). This prevents UNIQUE() conflicts.
	//
	// The CTE selects candidate attendance rows (NULL subject, in time range),
	// then filters out candidates where a conflicting attendance already exists.
	updateSQL := `
	WITH candidates AS (
	  SELECT a.attendance_id, a.student_id, a.date
	  FROM attendance a
	  WHERE a.subject_id IS NULL AND a.recorded_at >= $1 AND a.recorded_at <= $2
	),
	to_update AS (
	  SELECT c.attendance_id
	  FROM candidates c
	  LEFT JOIN attendance existing
	    ON existing.student_id = c.student_id
	   AND existing.subject_id = $3
	   AND existing.date = c.date
	  WHERE existing.attendance_id IS NULL
	)
	UPDATE attendance
	SET subject_id = $3
	WHERE attendance_id IN (SELECT attendance_id FROM to_update)
	RETURNING attendance_id;
	`

	rows, err := tx.Query(updateSQL, start, end, subjectID)
	if err != nil {
		return 0, 0, fmt.Errorf("update attendance: %w", err)
	}
	var updatedCount int64
	for rows.Next() {
		updatedCount++
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return 0, 0, fmt.Errorf("rows err: %w", err)
	}

	// 4) skipped = totalCandidates - updatedCount (these were left because of conflicts)
	skipped := totalCandidates - updatedCount
	if skipped < 0 {
		skipped = 0
	}

	if err := tx.Commit(); err != nil {
		return 0, 0, fmt.Errorf("commit tx: %w", err)
	}
	return updatedCount, skipped, nil
}

// Get summary for one student across all subjects
func (p *PostgresRepo) GetAttendanceSummaryByStudent(studentID int64) ([]domain.SubjectSummary, error) {
    q := `
    SELECT s.subject_id, subj.subject_name,
           COUNT(*) AS total_classes,
           SUM(CASE WHEN a.status = 'Present' THEN 1 ELSE 0 END) AS attended,
           ROUND(100.0 * SUM(CASE WHEN a.status = 'Present' THEN 1 ELSE 0 END) / COUNT(*), 2) AS percentage
    FROM attendance a
    JOIN subjects subj ON a.subject_id = subj.subject_id
    JOIN student_subjects s ON s.subject_id = subj.subject_id AND s.student_id = a.student_id
    WHERE a.student_id = $1 AND a.subject_id IS NOT NULL
    GROUP BY s.subject_id, subj.subject_name;
    `

    rows, err := p.db.Query(q, studentID)
    if err != nil {
        return nil, fmt.Errorf("get student summary: %w", err)
    }
    defer rows.Close()

    var list []domain.SubjectSummary
    for rows.Next() {
        var s domain.SubjectSummary
        if err := rows.Scan(&s.SubjectID, &s.SubjectName, &s.TotalClasses, &s.Attended, &s.Percentage); err != nil {
            return nil, err
        }
        list = append(list, s)
    }
    return list, nil
}

// Get summary for one subject across all students
func (p *PostgresRepo) GetAttendanceSummaryBySubject(subjectID int64) ([]domain.StudentSummary, error) {
    q := `
    SELECT st.student_id, st.student_name,
           COUNT(*) AS total_classes,
           SUM(CASE WHEN a.status = 'Present' THEN 1 ELSE 0 END) AS attended,
           ROUND(100.0 * SUM(CASE WHEN a.status = 'Present' THEN 1 ELSE 0 END) / COUNT(*), 2) AS percentage
    FROM attendance a
    JOIN students st ON a.student_id = st.student_id
    WHERE a.subject_id = $1
    GROUP BY st.student_id, st.student_name;
    `

    rows, err := p.db.Query(q, subjectID)
    if err != nil {
        return nil, fmt.Errorf("get subject summary: %w", err)
    }
    defer rows.Close()

    var list []domain.StudentSummary
    for rows.Next() {
        var s domain.StudentSummary
        if err := rows.Scan(&s.StudentID, &s.StudentName, &s.TotalClasses, &s.Attended, &s.Percentage); err != nil {
            return nil, err
        }
        list = append(list, s)
    }
    return list, nil
}


func (p *PostgresRepo) GetClassAttendance(subjectID int64, date time.Time) ([]domain.ClassAttendance, error) {
    q := `
    SELECT st.student_id, st.student_name, a.date, a.status
    FROM attendance a
    JOIN students st ON a.student_id = st.student_id
    WHERE a.subject_id = $1 AND a.date = $2
    ORDER BY st.student_name;
    `

    rows, err := p.db.Query(q, subjectID, date)
    if err != nil {
        return nil, fmt.Errorf("get class attendance: %w", err)
    }
    defer rows.Close()

    var list []domain.ClassAttendance
    for rows.Next() {
        var ca domain.ClassAttendance
        if err := rows.Scan(&ca.StudentID, &ca.StudentName, &ca.Date, &ca.Status); err != nil {
            return nil, err
        }
        list = append(list, ca)
    }
    return list, nil
}

func (p *PostgresRepo) GetStudentAttendanceHistory(studentID int64, subjectID int64) ([]domain.Attendance, error) {
    q := `SELECT attendance_id, student_id, subject_id, date, status, recorded_at, created_at
          FROM attendance
          WHERE student_id = $1 AND subject_id = $2
          ORDER BY date ASC;`

    rows, err := p.db.Query(q, studentID, subjectID)
    if err != nil {
        return nil, fmt.Errorf("get student history: %w", err)
    }
    defer rows.Close()

    var list []domain.Attendance
    for rows.Next() {
        var a domain.Attendance
        if err := rows.Scan(&a.ID, &a.StudentID, &a.SubjectID, &a.Date, &a.Status, &a.RecordedAt, &a.CreatedAt); err != nil {
            return nil, err
        }
        list = append(list, a)
    }
    return list, nil
}

func (p *PostgresRepo) ExportSubjectAttendanceCSV(subjectID int64, fromDate, toDate time.Time, filePath string) error {
    q := `
    SELECT st.student_id, st.student_name, a.date, a.status
    FROM attendance a
    JOIN students st ON a.student_id = st.student_id
    WHERE a.subject_id = $1 AND a.date BETWEEN $2 AND $3
    ORDER BY a.date, st.student_name;
    `

    rows, err := p.db.Query(q, subjectID, fromDate, toDate)
    if err != nil {
        return fmt.Errorf("export query: %w", err)
    }
    defer rows.Close()

    f, err := os.Create(filePath)
    if err != nil {
        return fmt.Errorf("create file: %w", err)
    }
    defer f.Close()

    writer := csv.NewWriter(f)
    defer writer.Flush()

    // header
    writer.Write([]string{"Student ID", "Student Name", "Date", "Status"})

    // rows
    for rows.Next() {
        var studentID int64
        var studentName, status string
        var date time.Time
        if err := rows.Scan(&studentID, &studentName, &date, &status); err != nil {
            return err
        }
        record := []string{
            fmt.Sprintf("%d", studentID),
            studentName,
            date.Format("2006-01-02"),
            status,
        }
        writer.Write(record)
    }

    return writer.Error()
}


