package student_service

import (
	"fmt"

	"github.com/suhas-developer07/Smart-Attendence-System/server/internals/domain"
)

type StudentService struct {
	studentRepo domain.StudentRepo
}

func NewStudentService(studentRepo domain.StudentRepo) *StudentService {
	return &StudentService{
		studentRepo: studentRepo,
	}
}

func (s *StudentService) RegisterStudent(req domain.StudentRegisterPayload) (int64, error) {

	// TODO : here i need call rabbit mq service and  drop a message "new entry with usn req.usn"
	id, err := s.studentRepo.StudentRegister(req)

	if err != nil {
		return 0,fmt.Errorf("Failed to register the student : %v",err)
	}

	return id, nil
}
