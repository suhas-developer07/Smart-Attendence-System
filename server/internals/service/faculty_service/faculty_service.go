package facultyservice

import "github.com/suhas-developer07/Smart-Attendence-System/server/internals/domain"


type Facultyservice struct{
	subjectRepo domain.SubjectRepo
}

func (s *Facultyservice) AddSubjectService(req domain.SubjectPayload) (int64, error){

	id ,err := s.subjectRepo.AddSubject(req)

	if err != nil{
		return 0,err
	}

	return id,nil

} 

