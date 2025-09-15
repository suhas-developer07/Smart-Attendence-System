package cmd

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	faculty_handler "github.com/suhas-developer07/Smart-Attendence-System/server/internals/handler/faculty"
	student_handler "github.com/suhas-developer07/Smart-Attendence-System/server/internals/handler/student"
	"github.com/suhas-developer07/Smart-Attendence-System/server/internals/repository"
	faculty_service "github.com/suhas-developer07/Smart-Attendence-System/server/internals/service/faculty"
	student_service "github.com/suhas-developer07/Smart-Attendence-System/server/internals/service/student"
)

func SetupRoutes(e *echo.Echo, db *sql.DB) {
	repo := repository.NewPostgresRepo(db)

	if err := repo.InitTables(); err != nil {
		log.Fatalf("Error Initializing the tables")
	}

	StudentService := student_service.NewStudentService(repo)

	studentHandler := student_handler.NewStudentHandler(StudentService)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "server is healthy")
	})

	e.POST("/students/register", studentHandler.StudentRegisterHandler)

	facultyService := faculty_service.NewFacultyService(repo)
	facultyHandler := faculty_handler.NewFacultyHandler(facultyService)

	e.POST("/faculty/addsubjects", facultyHandler.AddSubject)	
}
