package faculty_service

import (
	"github.com/go-playground/validator/v10"
	"github.com/suhas-developer07/Smart-Attendence-System/server/internals/domain"
)

type FacultyService struct {
	facultyRepo domain.FacultyRepo
	validate    *validator.Validate
}

func NewFacultyService(facultyRepo domain.FacultyRepo) *FacultyService {
	v := validator.New()
	return &FacultyService{
		facultyRepo: facultyRepo,
		validate:    v,
	}
}

func (s *FacultyService) GetFacultyByID(facultyID int) (domain.Faculty, error) {
	faculty, err := s.facultyRepo.GetFacultyByID(facultyID)
	if err != nil {
		return domain.Faculty{}, err
	}
	return faculty, nil
}
