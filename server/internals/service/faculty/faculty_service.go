package faculty_service

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/suhas-developer07/Smart-Attendence-System/server/internals/domain"
	"github.com/suhas-developer07/Smart-Attendence-System/server/pkg/utils"
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

func (s *FacultyService) RegisterFaculty(req domain.FacultyRegisterPayload) (int64,error) {
	if err := s.validate.Struct(req); err != nil {
		return 0, fmt.Errorf("validation error: %w", err)
	}

	hashedPassword, err := utils.HashPassword(req.Password)

	if err != nil {
		return 0, fmt.Errorf("error while hashing password: %w", err)
	}

	req.Password = hashedPassword

	id, err := s.facultyRepo.CreateFaculty(req)

	if err != nil {
		return 0, fmt.Errorf("error while creating faculty: %w", err)
	}

	return id, nil	
}


func (s *FacultyService) AuthenticateFaculty(req domain.FacultyLoginPayload) (int64, error) {
	if err := s.validate.Struct(req); err != nil {
		return 0, fmt.Errorf("validation error: %w", err)
	}

	facultyID, err := s.facultyRepo.AuthenticateFaculty(req)
	if err != nil {
		return 0, fmt.Errorf("authentication failed: %w", err)
	}

	return facultyID, nil
}

func (s *FacultyService) GetFacultyByID(facultyID int) (domain.Faculty, error) {
	faculty, err := s.facultyRepo.GetFacultyByID(facultyID)
	if err != nil {
		return domain.Faculty{}, err
	}
	return faculty, nil
}


func (s *FacultyService) GetAllFaculty() ([]domain.Faculty, error) {
	faculties, err := s.facultyRepo.GetAllFaculty()
	if err != nil {
		return nil, err
	}
	return faculties, nil
}

func (s *FacultyService) GetFacultyByDepartment(department string) ([]domain.Faculty, error) {
	faculties, err := s.facultyRepo.GetAllFaculty()
	if err != nil {
		return nil, err
	}

	var filtered []domain.Faculty
	for _, f := range faculties {
		if f.Department == department {
			filtered = append(filtered, f)
		}
	}

	return filtered, nil
}

