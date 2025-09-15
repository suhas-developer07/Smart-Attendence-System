package stundentservice

import (
	"fmt"

	"github.com/suhas-developer07/Smart-Attendence-System/server/internals/domain"
)

type StudentService struct {
	studentRepo domain.StudentRepo
	subjectRepo domain.SubjectRepo
}

func NewStudentService(studentRepo domain.StudentRepo,subjectRepo domain.SubjectRepo) *StudentService {
	return &StudentService{
		studentRepo: studentRepo,
	    subjectRepo:subjectRepo,
	}
}

func (s *StudentService) RegisterStudentService(req domain.StudentRegisterPayload) (int64, error) {

	// TODO : here i need call rabbit mq service and  drop a message "new entry with usn req.usn"
	id, err := s.studentRepo.StudentRegister(req)

	if err != nil {
		return 0,fmt.Errorf("Failed to register the student : %v",err)
	}

	return id, nil
}


func (s *StudentService) AddSubjectService(req domain.SubjectPayload) (int64, error){

	id ,err := s.subjectRepo.AddSubject(req)

	if err != nil{
		return 0,err
	}

	return id,nil

} 

