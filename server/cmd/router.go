package cmd

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/suhas-developer07/Smart-Attendence-System/server/internals/handler/faculty"
	student_handler "github.com/suhas-developer07/Smart-Attendence-System/server/internals/handler/student"
	subject_handler "github.com/suhas-developer07/Smart-Attendence-System/server/internals/handler/subjects"
	"github.com/suhas-developer07/Smart-Attendence-System/server/internals/repository"

	faculty_service "github.com/suhas-developer07/Smart-Attendence-System/server/internals/service/faculty"
	student_service "github.com/suhas-developer07/Smart-Attendence-System/server/internals/service/student"
	subject_service "github.com/suhas-developer07/Smart-Attendence-System/server/internals/service/subject"
)
func SetupRoutes(e *echo.Echo, db *sql.DB) {
	repo := repository.NewPostgresRepo(db)

	if err := repo.InitTables(); err != nil {
		log.Fatalf("Error Initializing the tables")
	}

	
	studentService := student_service.NewStudentService(repo)
	studentHandler := student_handler.NewStudentHandler(studentService)

	subjectService := subject_service.NewSubjectService(repo)
	subjectHandler := subject_handler.NewSubjectHandler(subjectService)

	facultyService := faculty_service.NewFacultyService(repo)
	facultyHandler := faculty.NewFacultyHandler(facultyService)


	//  Student 
	student := e.Group("/students")
	{
		student.POST("", studentHandler.StudentRegisterHandler)       
		student.PUT("/:student_id", studentHandler.UpdateStudentInfoHandler) 
		//student.GET("/:student_id", studentHandler.GetStudentByIDHandler)    
		student.GET("/:student_id/subjects", subjectHandler.GetSubjectsByStudentIDHandler) 

	//  Subject 
	subject := e.Group("/subjects")
	{
		subject.POST("", subjectHandler.AddSubjectHandler)                        
		subject.GET("", subjectHandler.GetSubjectsByDeptAndSemHandler)            
		subject.GET("/faculty/:faculty_id", subjectHandler.GetSubjectsByFacultyIDHandler) 
	}

	// Faculty
	faculty := e.Group("/faculty")
	{
		faculty.POST("/register", facultyHandler.RegisterFacultyHandler) 
		faculty.POST("/login", facultyHandler.AuthenticateFacultyHandler) 
		faculty.GET("/:faculty_id", facultyHandler.GetFacultyByIDHandler) 
		faculty.GET("", facultyHandler.GetAllFacultyHandler)              
		faculty.GET("/department/:dept", facultyHandler.GetFacultyByDepartmentHandler) 
	}

	// ---------- Health ----------
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "server is healthy")
	})
}
}