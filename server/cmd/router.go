package cmd

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	student_handler "github.com/suhas-developer07/Smart-Attendence-System/server/internals/handler/student"
	subject_handler "github.com/suhas-developer07/Smart-Attendence-System/server/internals/handler/subjects"
	"github.com/suhas-developer07/Smart-Attendence-System/server/internals/repository"

	student_service "github.com/suhas-developer07/Smart-Attendence-System/server/internals/service/student"
	subject_service "github.com/suhas-developer07/Smart-Attendence-System/server/internals/service/subject"
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


	subjectService := subject_service.NewSubjectService(repo)
	subjectHandler := subject_handler.NewSubjectHandler(subjectService)

	e.POST("/faculty/addsubjects", subjectHandler.AddSubjectHandler)
	e.GET("/students/listsubjects", subjectHandler.ListSubjectsHandler)
	//e.POST("/faculty/getstudentswithsubjects", facultyHandler.GetStudentsWithSubjectsHandler)
}
