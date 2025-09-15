package student_service

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/suhas-developer07/Smart-Attendence-System/server/internals/domain"
)

type StudentService struct {
	studentRepo domain.StudentRepo
	validate *validator.Validate
}

func NewStudentService(studentRepo domain.StudentRepo) *StudentService {
	v := validator.New()
	return &StudentService{
		studentRepo: studentRepo,
		validate: v,
	}
}

func (s *StudentService) RegisterStudent(req domain.StudentRegisterPayload) (int64, error) {

	if err := s.validate.Struct(req); err != nil {
		return 0, fmt.Errorf("validation error: %w", err)
	}

	// TODO : here i need call rabbit mq service and  drop a message "new entry with usn req.usn"
	id, err := s.studentRepo.StudentRegister(req)

	if err != nil {
		return 0,err
	}

	return id, nil
}

func (s *StudentService) UpdateStudentInfo(studentID int, payload domain.StudentUpdatePayload) error {

	if err := s.validate.Struct(payload); err != nil {
		return err
	}
	err := s.studentRepo.UpdateStudentInfo(studentID, payload)

	if err != nil {
		return err
	}
	return nil
}
