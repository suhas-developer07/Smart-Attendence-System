package attendence_service

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/suhas-developer07/Smart-Attendence-System/server/internals/domain"
)

type AttendanceService struct {
	attendanceRepo domain.AttendanceRepository
	validate   *validator.Validate
}

func NewAttendanceService(attendanceRepo domain.AttendanceRepository) *AttendanceService {
	v := validator.New()
	return &AttendanceService{
		attendanceRepo: attendanceRepo,
		validate:       v,
	}
}

func (s *AttendanceService) MarkAttendance(attendance *domain.AttendancePayload) (int64, error) {
	if err := s.validate.Struct(attendance); err != nil {
    fmt.Printf("%#v\n", attendance) 
}

	id, err := s.attendanceRepo.MarkAttendance(attendance)
	if err != nil {
		return 0, fmt.Errorf("error marking attendance: %w", err)
	}
	return id, nil
}

func (s *AttendanceService) BulkMarkAttendance(attendances []domain.AttendancePayload) (int, error) {
    // Validate each attendance
    for _, a := range attendances {
        if err := s.validate.Struct(a); err != nil {
            return 0, fmt.Errorf("validation failed for USN=%s: %w", a.USN, err)
        }
    }

    // Call repo
    return s.attendanceRepo.BulkMarkAttendance(attendances)
}


func (s *AttendanceService) GetAttendanceByStudentAndSubject(usn string, subjectCode string) ([]domain.AttendanceWithNames, error) {
	if err := s.validate.Var(usn, "required"); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}
	if err := s.validate.Var(subjectCode, "required"); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	attendances, err := s.attendanceRepo.GetAttendanceByStudentAndSubject(usn, subjectCode)
	if err != nil {
		return nil, fmt.Errorf("error fetching attendance: %w", err)
	}
	return attendances, nil
}
func (s *AttendanceService) GetAttendanceBySubjectAndDate(subjectCode string, date time.Time) ([]domain.AttendanceWithNames, error) {

	if subjectCode == "" {
		return nil, fmt.Errorf("validation error: subject_id is required")
	}
	if date.IsZero() {
		return nil, fmt.Errorf("validation error: date is required")
	}

	attendances, err := s.attendanceRepo.GetAttendanceBySubjectAndDate(subjectCode, date)
	if err != nil {
		return nil, fmt.Errorf("error fetching attendance: %w", err)
	}
	return attendances, nil
}


func (s *AttendanceService) AssignSubjectToTimeRange(
	facultyID int64,
	subjectCode string,
	classDate time.Time, // only the date of class
	start time.Time,     // class start time
	end time.Time,       // class end time
) (int64, int64, error) {


	if err := s.validate.Var(facultyID, "required"); err != nil {
		return 0, 0, fmt.Errorf("validation error: %w", err)
	}
	if err := s.validate.Var(subjectCode, "required"); err != nil {
		return 0, 0, fmt.Errorf("validation error: %w", err)
	}
	if classDate.IsZero() || start.IsZero() || end.IsZero() {
		return 0, 0, fmt.Errorf("start, end, and class date are required")
	}

	updatedCount, skipped, err := s.attendanceRepo.AssignSubjectToTimeRange(facultyID, subjectCode, classDate, start, end)
	if err != nil {
		return 0, 0, fmt.Errorf("error assigning subject to time range: %w", err)
	}
	return updatedCount, skipped, nil
}


func (s *AttendanceService) GetAttendanceSummaryBySubject(subjectCode string) ([]domain.StudentSummary, error) {

	if err := s.validate.Var(subjectCode, "required"); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	summaries, err := s.attendanceRepo.GetAttendanceSummaryBySubject(subjectCode)
	if err != nil {
		return nil, fmt.Errorf("error fetching attendance summary: %w", err)
	}
	return summaries, nil
}

func (s *AttendanceService) GetClassAttendance(subjectCode string, date time.Time) ([]domain.ClassAttendance, error) {
	if err := s.validate.Var(subjectCode, "required"); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}
	if date.IsZero() {
		return nil, fmt.Errorf("validation error: date is required")
	}

	attendances, err := s.attendanceRepo.GetClassAttendance(subjectCode, date)
	if err != nil {
		return nil, fmt.Errorf("error fetching class attendance: %w", err)
	}
	return attendances, nil
}

func (s *AttendanceService) GetStudentAttendanceHistory(usn string, subjectCode string)([]domain.StudentHistory,error){
	if err := s.validate.Var(usn, "required"); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}
	if err := s.validate.Var(subjectCode, "required"); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	history, err := s.attendanceRepo.GetStudentAttendanceHistory(usn, subjectCode)
	if err != nil {
		return nil, fmt.Errorf("error fetching student attendance history: %w", err)
	}
	return history, nil
}

func (s *AttendanceService) GetAttendanceSummaryByStudent(usn string) ([]domain.SubjectSummary, error) {
	if err := s.validate.Var(usn, "required"); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	summaries, err := s.attendanceRepo.GetAttendanceSummaryByStudent(usn)
	if err != nil {
		return nil, fmt.Errorf("error fetching attendance summary: %w", err)
	}
	return summaries, nil
}


