package subject_service

import (
	"github.com/go-playground/validator/v10"
	"github.com/suhas-developer07/Smart-Attendence-System/server/internals/domain"
)


type SubjectService struct{
	subjectRepo domain.SubjectRepo
	validator *validator.Validate
}

func NewSubjectService(subjectRepo domain.SubjectRepo) *SubjectService {
	v := validator.New()
	return &SubjectService{
		subjectRepo: subjectRepo,
		validator:   v,
	}
}

func (s *SubjectService) AddSubject(req domain.SubjectPayload) (int64, error){

	id ,err := s.subjectRepo.AddSubject(req)

	if err != nil{
		return 0,err
	}

	return id,nil
}

func (s *SubjectService) GetSubjectsByDeptAndSem(department string, sem int) ([]domain.Subject, error){

	subjects, err := s.subjectRepo.GetSubjectsByDeptAndSem(department, sem)
	if err != nil {
		return nil, err
	}
	return subjects, nil
}

func (s *SubjectService) GetSubjectsByFacultyID(facultyID int) ([]domain.Subject, error){

	subjects,err := s.subjectRepo.GetSubjectsByFacultyID(facultyID)

	if err != nil {
		return nil,err
	}

	return subjects,nil
}
//subjects of a particular student
func (s *SubjectService) GetSubjectsByStudentID(studentID int64) ([]domain.SubjectPayload, error){

	subjects, err := s.subjectRepo.GetSubjectsByStudentID(studentID)
	if err != nil {
		return nil, err
	}
	return subjects, nil
}


 

// func (s *FacultyService) GetStudentsWithSubjects(req domain.GetStudentsWithSubjectsPayload) ([]domain.StudentWithSubjects, error){

// 	students, err := s.subjectRepo.GetStudentsWithSubjects(req)

// 	if err != nil{
// 		return nil,err
// 	}

// 	return students,nil
// }	
