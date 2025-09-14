package cmd

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/suhas-developer07/Smart-Attendence-System/server/internals/handler"
	"github.com/suhas-developer07/Smart-Attendence-System/server/internals/repository"
	"github.com/suhas-developer07/Smart-Attendence-System/server/internals/service"
)

func SetupRoutes(e *echo.Echo, db *sql.DB) {
	repo := repository.NewPostgresRepo(db)

	StudentService := service.NewStudentService(repo)

	studentHandler := handler.NewStudentHandler(StudentService)

	e.POST("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "server is healthy")
	})

	e.POST("/students/register", studentHandler.StudentRegisterHandler)
}

