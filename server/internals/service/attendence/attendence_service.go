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
		return 0, fmt.Errorf("validation error: %w", err)
	}
	id, err := s.attendanceRepo.MarkAttendance(attendance)
	if err != nil {
		return 0, fmt.Errorf("error marking attendance: %w", err)
	}
	return id, nil
	
}

func (s *AttendanceService) GetAttendanceByStudentAndSubject(studentID int64, subjectID int64) ([]domain.Attendance, error) {
	if err := s.validate.Var(studentID, "required"); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}
	if err := s.validate.Var(subjectID, "required"); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	attendances, err := s.attendanceRepo.GetAttendanceByStudentAndSubject(studentID, subjectID)
	if err != nil {
		return nil, fmt.Errorf("error fetching attendance: %w", err)
	}
	return attendances, nil
}

func (s *AttendanceService) GetAttendanceBySubject(subjectID int64, fromDate, toDate time.Time) ([]domain.Attendance, error) {
	if err := s.validate.Var(subjectID, "required"); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}
	if err := s.validate.Var(fromDate, "required,datetime=2006-01-02"); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}
	if err := s.validate.Var(toDate, "required,datetime=2006-01-02"); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	attendances, err := s.attendanceRepo.GetAttendanceBySubject(subjectID, fromDate, toDate)
	if err != nil {
		return nil, fmt.Errorf("error fetching attendance: %w", err)
	}
	return attendances, nil
}

func (s *AttendanceService) AssignSubjectToTimeRange(facultyID int, subjectID int64, start, end time.Time) (int64, int64, error) {

	if err := s.validate.Var(facultyID, "required"); err != nil {
		return 0, 0, fmt.Errorf("validation error: %w", err)
	}
	if err := s.validate.Var(subjectID, "required"); err != nil {
		return 0, 0, fmt.Errorf("validation error: %w", err)
	}
	if err := s.validate.Var(start, "required,datetime=2006-01-02"); err != nil {
		return 0, 0, fmt.Errorf("validation error: %w", err)
	}
	if err := s.validate.Var(end, "required,datetime=2006-01-02"); err != nil {
		return 0, 0, fmt.Errorf("validation error: %w", err)
	}

	updatedCount, skipped, err := s.attendanceRepo.AssignSubjectToTimeRange(facultyID, subjectID, start, end)
	if err != nil {
		return 0, 0, fmt.Errorf("error assigning subject to time range: %w", err)
	}
	return updatedCount, skipped, nil
}

func (s *AttendanceService) GetAttendanceSummaryBySubject(subjectID int64) ([]domain.StudentSummary, error) {

	if err := s.validate.Var(subjectID, "required"); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	summaries, err := s.attendanceRepo.GetAttendanceSummaryBySubject(subjectID)
	if err != nil {
		return nil, fmt.Errorf("error fetching attendance summary: %w", err)
	}
	return summaries, nil
}

func (s *AttendanceService) GetClassAttendance(subjectID int64,date time.Time)([]domain.ClassAttendance,error){
	if err := s.validate.Var(subjectID, "required"); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}
	if err := s.validate.Var(date, "required,datetime=2006-01-02"); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	attendances, err := s.attendanceRepo.GetClassAttendance(subjectID,date)
	if err != nil {
		return nil, fmt.Errorf("error fetching class attendance: %w", err)
	}
	return attendances, nil
}

func (s *AttendanceService) GetStudentAttendanceHistory(studentID int64, subjectID int64)([]domain.Attendance,error){
	if err := s.validate.Var(studentID, "required"); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}
	if err := s.validate.Var(subjectID, "required"); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	history, err := s.attendanceRepo.GetStudentAttendanceHistory(studentID,subjectID)
	if err != nil {
		return nil, fmt.Errorf("error fetching student attendance history: %w", err)
	}
	return history, nil
}



