package service

import (
	"fmt"

	"github.com/suhas-developer07/Smart-Attendence-System/server/internals/domain"
)

type StudentService struct {
	repo domain.StudentRepo
}

func NewStudentService(r domain.StudentRepo) *StudentService {
	return &StudentService{repo: r}
}

func (s *StudentService) RegisterStudentService(req domain.StudentRegisterPayload) (int64, error) {

	id, err := s.repo.StudentRegister(req)

	if err != nil {
		return 0, fmt.Errorf("Error Registering student : %v", err)
	}

	return id, nil
}
