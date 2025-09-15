package faculty_service

import "github.com/suhas-developer07/Smart-Attendence-System/server/internals/domain"


type FacultyService struct{
	subjectRepo domain.SubjectRepo
}
func NewFacultyService(subjectRepo domain.SubjectRepo) *FacultyService {
	return &FacultyService{
		subjectRepo:subjectRepo,
	}
}

func (s *FacultyService) AddSubject(req domain.SubjectPayload) (int64, error){

	id ,err := s.subjectRepo.AddSubject(req)

	if err != nil{
		return 0,err
	}

	return id,nil

} 

