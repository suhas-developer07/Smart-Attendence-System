package facultyhandler

import facultyservice "github.com/suhas-developer07/Smart-Attendence-System/server/internals/service/faculty_service"

type FacultyHandler struct{
     facultyRepo *facultyservice.Facultyservice

}

func NewfacultyHandler(fr *facultyservice.Facultyservice)*FacultyHandler{
	return &FacultyHandler{
		facultyRepo: fr,
	}
}


