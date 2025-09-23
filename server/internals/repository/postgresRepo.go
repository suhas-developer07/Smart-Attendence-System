package repository

import (
	"database/sql"
	"fmt"
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

		`CREATE TABLE IF NOT EXISTS faculty (
			faculty_id SERIAL PRIMARY KEY,
			faculty_name VARCHAR(100) NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			password_hash VARCHAR(256) NOT NULL,
			department VARCHAR(50) NOT NULL,
			created_at TIMESTAMPTZ DEFAULT now()
		);`,

		`CREATE TABLE IF NOT EXISTS students (
			student_id SERIAL PRIMARY KEY,
			usn VARCHAR(50) UNIQUE NOT NULL,
			username VARCHAR(100) NOT NULL,
			password_hash VARCHAR(256) NULL,
			department VARCHAR(50) NOT NULL,
			sem INT NOT NULL,
			face_encoding BYTEA NULL,
			nfc_uid VARCHAR(100) UNIQUE NULL,
			created_at TIMESTAMPTZ DEFAULT now()
		);`,


		`CREATE TABLE IF NOT EXISTS subjects (
			subject_id SERIAL PRIMARY KEY,
			subject_code VARCHAR(50) UNIQUE NOT NULL,
			subject_name VARCHAR(150) NOT NULL,
			department VARCHAR(50) NOT NULL,
			sem INT NOT NULL,
			faculty_id INT NOT NULL,
			CONSTRAINT fk_faculty_sub FOREIGN KEY (faculty_id) REFERENCES faculty(faculty_id) ON DELETE RESTRICT
		);`,

		`CREATE TABLE IF NOT EXISTS student_subjects (
			student_id INT NOT NULL,
			subject_id INT NOT NULL,
			PRIMARY KEY (student_id, subject_id),
			CONSTRAINT fk_student_sub FOREIGN KEY (student_id) REFERENCES students(student_id) ON DELETE CASCADE,
			CONSTRAINT fk_subject_sub FOREIGN KEY (subject_id) REFERENCES subjects(subject_id) ON DELETE CASCADE
		);`,

		`CREATE TABLE IF NOT EXISTS attendance (
            attendance_id SERIAL PRIMARY KEY,
            usn VARCHAR(50) NOT NULL,
            subject_id INT NULL,
            date DATE NOT NULL,
            status VARCHAR(20) NOT NULL CHECK (status IN ('Present', 'Absent')),
            recorded_at TIMESTAMPTZ NOT NULL DEFAULT now(),
            created_at TIMESTAMPTZ DEFAULT now(),
            updated_at TIMESTAMPTZ DEFAULT now(),
            CONSTRAINT fk_attendance_student FOREIGN KEY (usn) REFERENCES students(usn) ON DELETE CASCADE,
            CONSTRAINT fk_attendance_subject FOREIGN KEY (subject_id) REFERENCES subjects(subject_id) ON DELETE SET NULL
       );
`,

		// 8. Attendance unique indexes for workflow
		`CREATE UNIQUE INDEX IF NOT EXISTS uniq_usn_date_null_subject
    ON attendance(usn, date)
    WHERE subject_id IS NULL;`,

		`CREATE UNIQUE INDEX IF NOT EXISTS uniq_usn_subject_date
    ON attendance(usn, subject_id, date);`,

		`CREATE INDEX IF NOT EXISTS idx_attendance_date_recorded
    ON attendance(date, recorded_at);`,

		`CREATE INDEX IF NOT EXISTS idx_attendance_subject_date
    ON attendance(subject_id, date);`,
	}

	for _, q := range queries {
		if _, err := p.db.Exec(q); err != nil {
			return fmt.Errorf("failed to exec init query: %w", err)
		}
	}

	return nil
}

//student 
func (p *PostgresRepo) StudentRegister(student domain.StudentRegisterPayload) (int64, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	fmt.Println("Department:",student.Department)

	var pwHash string
	if student.Password != "" {
		pwHash, err = utils.HashPassword(student.Password)
		if err != nil {
			return 0, fmt.Errorf("hash password: %w", err)
		}
	}

	var id int64
	query := `INSERT INTO students (usn, username, password_hash, department, sem)
	          VALUES ($1, $2, $3, $4, $5) RETURNING student_id;`
	err = tx.QueryRow(query, student.USN, student.Username, pwHash, student.Department, student.Sem).Scan(&id)
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

func (p *PostgresRepo) LoginStudent(usn, password string) (string, error) {
	var pwHash string
	var studentID int64
	q := `SELECT student_id, password_hash FROM students WHERE usn = $1;`
	if err := p.db.QueryRow(q, usn).Scan(&studentID, &pwHash); err != nil {
		return "", fmt.Errorf("query student: %w", err)
	}

	if err := utils.ComparePassword(pwHash, password); err != nil {
		return "", fmt.Errorf("invalid credentials: %w", err)
	}

	token, err := utils.GenerateTokenForStudent(studentID, usn)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	return token, nil
}

func (p *PostgresRepo) UpdateStudentInfo(studentID int, payload domain.StudentUpdatePayload) error {
	query := `UPDATE students SET username = $2, department = $3, sem = $4 WHERE student_id = $1;`
	if _, err := p.db.Exec(query, studentID, payload.Username, payload.Department, payload.Sem); err != nil {
		return fmt.Errorf("update student: %w", err)
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
func (p *PostgresRepo) GetSubjectsByFacultyID(facultyID int64) ([]domain.Subject, error) {
	fmt.Println("DEBUG: facultyID in repo:", facultyID)
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

// Register Faculty
func (p *PostgresRepo) CreateFaculty(req domain.FacultyRegisterPayload) (int64, error) {
    // 1. Hash password
    pwHash, err := utils.HashPassword(req.Password)
    if err != nil {
        return 0, fmt.Errorf("hash password: %w", err)
    }

    // 2. Insert into DB
    var id int64
    q := `
        INSERT INTO faculty (faculty_name, email, password_hash, department)
        VALUES ($1, $2, $3, $4)
        RETURNING faculty_id;
    `
    if err := p.db.QueryRow(q, req.Name, req.Email, pwHash, req.Department).Scan(&id); err != nil {
        return 0, fmt.Errorf("insert faculty: %w", err)
    }

    return id, nil
}

// Authenticate Faculty
func (p *PostgresRepo) AuthenticateFaculty(req domain.FacultyLoginPayload) (string, error) {
    var id int64
    var pwHash string

    // 1. Get faculty by email
    q := `SELECT faculty_id, password_hash FROM faculty WHERE email = $1;`
    if err := p.db.QueryRow(q, req.Email).Scan(&id, &pwHash); err != nil {
        if err == sql.ErrNoRows {
            return "", fmt.Errorf("faculty not found")
        }
        return "", fmt.Errorf("query faculty: %w", err)
    }

    // 2. Compare password
    if err := utils.ComparePassword(pwHash, req.Password); err != nil {
		fmt.Println("DEBUG: Stored hash:", pwHash)
fmt.Println("DEBUG: Login password:", req.Password)

if err := utils.ComparePassword(pwHash, req.Password); err != nil {
    fmt.Println("DEBUG: Compare failed:", err)
    return "", fmt.Errorf("invalid credentials")
}

        return "", fmt.Errorf("invalid credentials")
    }

    // 3. Generate JWT
    token, err := utils.GenerateTokenForFaculty(id, req.Email)
    if err != nil {
        return "", fmt.Errorf("generate token: %w", err)
    }

    return token, nil
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

func (p *PostgresRepo) GetFacultyByID(facultyID int64) (domain.Faculty, error) {
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


func (p *PostgresRepo) MarkAttendance(req *domain.AttendancePayload) (int64, error) {
	var attendanceID int64

	// Attendance date only (UTC, truncate to date)
	classDate := req.RecordedAt.UTC().Truncate(24 * time.Hour)

	// Insert attendance without subject
	query := `
	INSERT INTO attendance (usn, subject_id, date, status, recorded_at)
	VALUES ($1, NULL, $2, $3, $4)
	ON CONFLICT (usn, date)
	WHERE subject_id IS NULL
	DO UPDATE SET status = EXCLUDED.status,
	              recorded_at = EXCLUDED.recorded_at
	RETURNING attendance_id;
	`

	err := p.db.QueryRow(query, req.USN, classDate, req.Status, req.RecordedAt.UTC()).Scan(&attendanceID)
	if err != nil {
		return 0, fmt.Errorf("mark attendance: %w", err)
	}

	return attendanceID, nil
}

func (p *PostgresRepo) BulkMarkAttendance(attendances []domain.AttendancePayload) (int, error) {
    tx, err := p.db.Begin()
    if err != nil {
        return 0, fmt.Errorf("begin tx: %w", err)
    }
    defer tx.Rollback()

    query := `
    INSERT INTO attendance (usn, subject_id, date, status, recorded_at)
    VALUES ($1, NULL, $2, $3, $4)
    ON CONFLICT (usn, date)
    WHERE subject_id IS NULL
    DO UPDATE SET status = EXCLUDED.status,
                  recorded_at = EXCLUDED.recorded_at
    RETURNING attendance_id;
    `

    stmt, err := tx.Prepare(query)
    if err != nil {
        return 0, fmt.Errorf("prepare stmt: %w", err)
    }
    defer stmt.Close()

    count := 0
    for _, a := range attendances {
        classDate := a.RecordedAt.UTC().Truncate(24 * time.Hour)

        var id int64
        if err := stmt.QueryRow(a.USN, classDate, a.Status, a.RecordedAt.UTC()).Scan(&id); err != nil {
            return 0, fmt.Errorf("insert attendance (usn=%s): %w", a.USN, err)
        }
        count++
    }

    if err := tx.Commit(); err != nil {
        return 0, fmt.Errorf("commit tx: %w", err)
    }

    return count, nil
}

//here iam assigning a subject to time range of attendance marker in the attendance table 
func (p *PostgresRepo) AssignSubjectToTimeRange(
	facultyID int64,
	subjectCode string,
	classDate time.Time,
	startTime, endTime time.Time,
) (int64, int64, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return 0, 0, fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	// Lookup subject_id from subject_code
	var subjectID int64
	var ownerID int64
	err = tx.QueryRow(`SELECT subject_id, faculty_id FROM subjects WHERE subject_code = $1`, subjectCode).Scan(&subjectID, &ownerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, fmt.Errorf("subject not found")
		}
		return 0, 0, fmt.Errorf("query subject by code: %w", err)
	}

	// Verify faculty owns this subject
	if ownerID != facultyID {
		return 0, 0, fmt.Errorf("not authorized to assign this subject")
	}

	// Build UTC timestamps for the given date + time range (IST -> UTC)
	loc, _ := time.LoadLocation("Asia/Kolkata")
	startDT := time.Date(classDate.Year(), classDate.Month(), classDate.Day(),
		startTime.Hour(), startTime.Minute(), 0, 0, loc).UTC()
	endDT := time.Date(classDate.Year(), classDate.Month(), classDate.Day(),
		endTime.Hour(), endTime.Minute(), 59, 999999999, loc).UTC()

	// Update attendance rows that have subject_id=NULL in the given time range
	updateSQL := `
	UPDATE attendance a
	SET subject_id = $1, updated_at = NOW()
	WHERE a.subject_id IS NULL
	  AND a.date = $2
	  AND a.recorded_at BETWEEN $3 AND $4
	  AND NOT EXISTS (
	    SELECT 1 FROM attendance existing
	    WHERE existing.usn = a.usn
	      AND existing.subject_id = $1
	      AND existing.date = a.date
	  )
	RETURNING attendance_id;
	`

	rows, err := tx.Query(updateSQL, subjectID, classDate.Format("2006-01-02"), startDT, endDT)
	if err != nil {
		return 0, 0, fmt.Errorf("update attendance: %w", err)
	}
	defer rows.Close()

	var updatedCount int64
	for rows.Next() {
		updatedCount++
	}
	if err := rows.Err(); err != nil {
		return 0, 0, fmt.Errorf("rows error: %w", err)
	}

	// Count skipped (attendance rows that were NULL but already had subject)
	var totalCandidates int64
	if err := tx.QueryRow(`
		SELECT COUNT(*) FROM attendance
		WHERE subject_id IS NULL
		  AND date = $1
		  AND recorded_at BETWEEN $2 AND $3
	`, classDate.Format("2006-01-02"), startDT, endDT).Scan(&totalCandidates); err != nil {
		return 0, 0, fmt.Errorf("count candidates: %w", err)
	}
	skipped := totalCandidates - updatedCount
	if skipped < 0 {
		skipped = 0
	}

	if err := tx.Commit(); err != nil {
		return 0, 0, fmt.Errorf("commit tx: %w", err)
	}

	return updatedCount, skipped, nil
}


func (p *PostgresRepo) GetAttendanceByStudentAndSubject(usn string, subjectCode string) ([]domain.AttendanceWithNames, error) {
	var subjectID int64
	err := p.db.QueryRow(`SELECT subject_id FROM subjects WHERE subject_code = $1`, subjectCode).Scan(&subjectID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("subject not found for code: %s", subjectCode)
		}
		return nil, fmt.Errorf("lookup subject_id: %w", err)
	}

	
	q := `
	SELECT a.attendance_id, a.usn, st.username AS student_name,
	       a.subject_id, sub.subject_name,
	       a.date, a.status, a.recorded_at, a.created_at
	FROM attendance a
	JOIN students st ON a.usn = st.usn
	LEFT JOIN subjects sub ON a.subject_id = sub.subject_id
	WHERE a.usn = $1 AND a.subject_id = $2
	ORDER BY a.date ASC;`

	rows, err := p.db.Query(q, usn, subjectID)
	if err != nil {
		return nil, fmt.Errorf("query attendance: %w", err)
	}
	defer rows.Close()

	var list []domain.AttendanceWithNames
	for rows.Next() {
		var a domain.AttendanceWithNames
		if err := rows.Scan(&a.ID, &a.USN, &a.StudentName, &a.SubjectID, &a.SubjectName,
			&a.Date, &a.Status, &a.RecordedAt, &a.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan attendance: %w", err)
		}
		list = append(list, a)
	}
	return list, rows.Err()
}

func (p *PostgresRepo) GetAttendanceBySubjectAndDate(subjectCode string, date time.Time) ([]domain.AttendanceWithNames, error) {

	var subjectID int64
	err := p.db.QueryRow(`SELECT subject_id FROM subjects WHERE subject_code = $1`, subjectCode).Scan(&subjectID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("subject not found for code: %s", subjectCode)
		}
		return nil, fmt.Errorf("lookup subject_id: %w", err)
	}

	loc, _ := time.LoadLocation("Asia/Kolkata")
	startDT := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc).UTC()
	endDT := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, loc).UTC()

	q := `
	SELECT a.attendance_id, a.usn, st.username AS student_name,
	       a.subject_id, sub.subject_name,
	       a.date, a.status, a.recorded_at, a.created_at
	FROM attendance a
	JOIN students st ON a.usn = st.usn
	LEFT JOIN subjects sub ON a.subject_id = sub.subject_id
	WHERE a.subject_id = $1
	  AND a.recorded_at BETWEEN $2 AND $3
	ORDER BY a.recorded_at ASC;`

	rows, err := p.db.Query(q, subjectID, startDT, endDT)
	if err != nil {
		return nil, fmt.Errorf("query attendance by subject and date: %w", err)
	}
	defer rows.Close()

	var list []domain.AttendanceWithNames
	for rows.Next() {
		var a domain.AttendanceWithNames
		if err := rows.Scan(&a.ID, &a.USN, &a.StudentName, &a.SubjectID, &a.SubjectName,
			&a.Date, &a.Status, &a.RecordedAt, &a.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan attendance: %w", err)
		}
		list = append(list, a)
	}
	return list, rows.Err()
}
//i need to write the service and handler for this function 
func (p *PostgresRepo) GetAttendanceSummaryByStudent(usn string) ([]domain.SubjectSummary, error) {
	q := `
	SELECT subj.subject_id, subj.subject_name,
	       COUNT(*) AS total_classes,
	       SUM(CASE WHEN a.status = 'Present' THEN 1 ELSE 0 END) AS attended,
	       ROUND(100.0 * SUM(CASE WHEN a.status = 'Present' THEN 1 ELSE 0 END) / COUNT(*), 2) AS percentage
	FROM attendance a
	JOIN subjects subj ON a.subject_id = subj.subject_id
	JOIN student_subjects s ON s.subject_id = subj.subject_id AND s.student_id = (
	    SELECT student_id FROM students WHERE usn = a.usn
	)
	WHERE a.usn = $1 AND a.subject_id IS NOT NULL
	GROUP BY subj.subject_id, subj.subject_name;`

	rows, err := p.db.Query(q, usn)
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


func (p *PostgresRepo) GetAttendanceSummaryBySubject(subjectCode string) ([]domain.StudentSummary, error) {
	var subjectID int64
	err := p.db.QueryRow(`SELECT subject_id FROM subjects WHERE subject_code = $1`, subjectCode).Scan(&subjectID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("subject not found for code: %s", subjectCode)
		}
		return nil, fmt.Errorf("lookup subject_id: %w", err)
	}

	q := `
	SELECT a.usn, st.username AS student_name,
	       COUNT(*) AS total_classes,
	       SUM(CASE WHEN a.status = 'Present' THEN 1 ELSE 0 END) AS attended,
	       ROUND(100.0 * SUM(CASE WHEN a.status = 'Present' THEN 1 ELSE 0 END) / COUNT(*), 2) AS percentage
	FROM attendance a
	JOIN students st ON a.usn = st.usn
	WHERE a.subject_id = $1
	GROUP BY a.usn, st.username
	ORDER BY st.username;`

	rows, err := p.db.Query(q, subjectID)
	if err != nil {
		return nil, fmt.Errorf("get subject summary: %w", err)
	}
	defer rows.Close()

	var list []domain.StudentSummary
	for rows.Next() {
		var s domain.StudentSummary
		if err := rows.Scan(&s.USN, &s.StudentName, &s.TotalClasses, &s.Attended, &s.Percentage); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, nil
}


func (p *PostgresRepo) GetClassAttendance(subjectCode string, date time.Time) ([]domain.ClassAttendance, error) {

	var subjectID int64
	err := p.db.QueryRow(`SELECT subject_id FROM subjects WHERE subject_code = $1`, subjectCode).Scan(&subjectID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("subject not found for code: %s", subjectCode)
		}
		return nil, fmt.Errorf("lookup subject_id: %w", err)
	}

	loc, _ := time.LoadLocation("Asia/Kolkata")
	startDT := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc).UTC()
	endDT := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, loc).UTC()

	q := `
	SELECT a.usn, st.username AS student_name, a.date, a.status
	FROM attendance a
	JOIN students st ON a.usn = st.usn
	WHERE a.subject_id = $1
	  AND a.recorded_at BETWEEN $2 AND $3
	ORDER BY st.username;`

	rows, err := p.db.Query(q, subjectID, startDT, endDT)
	if err != nil {
		return nil, fmt.Errorf("get class attendance: %w", err)
	}
	defer rows.Close()

	var list []domain.ClassAttendance
	for rows.Next() {
		var ca domain.ClassAttendance
		if err := rows.Scan(&ca.USN, &ca.StudentName, &ca.Date, &ca.Status); err != nil {
			return nil, err
		}
		list = append(list, ca)
	}
	return list, nil
}

func (p *PostgresRepo) GetStudentAttendanceHistory(usn string, subjectCode string) ([]domain.StudentHistory, error) {
var subjectID int64
	err := p.db.QueryRow(`SELECT subject_id FROM subjects WHERE subject_code = $1`, subjectCode).Scan(&subjectID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("subject not found for code: %s", subjectCode)
		}
		return nil, fmt.Errorf("lookup subject_id: %w", err)
	}

	q := `
	SELECT a.attendance_id, a.date, a.status,
	       a.subject_id, sub.subject_name,
	       a.recorded_at
	FROM attendance a
	LEFT JOIN subjects sub ON a.subject_id = sub.subject_id
	WHERE a.usn = $1 AND a.subject_id = $2
	ORDER BY a.date ASC;`

	rows, err := p.db.Query(q, usn, subjectID)
	if err != nil {
		return nil, fmt.Errorf("get student history: %w", err)
	}
	defer rows.Close()

	var list []domain.StudentHistory
	for rows.Next() {
		var h domain.StudentHistory
		if err := rows.Scan(&h.ID, &h.Date, &h.Status, &h.SubjectID, &h.SubjectName, &h.RecordedAt); err != nil {
			return nil, err
		}
		list = append(list, h)
	}
	return list, nil
}