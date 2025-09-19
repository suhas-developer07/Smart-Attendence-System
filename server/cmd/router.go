package cmd

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	attendance_handler "github.com/suhas-developer07/Smart-Attendence-System/server/internals/handler/attendance"
	"github.com/suhas-developer07/Smart-Attendence-System/server/internals/handler/faculty"
	student_handler "github.com/suhas-developer07/Smart-Attendence-System/server/internals/handler/student"
	subject_handler "github.com/suhas-developer07/Smart-Attendence-System/server/internals/handler/subjects"
	facultymiddlerware "github.com/suhas-developer07/Smart-Attendence-System/server/internals/middlerwares/faculty_middlerware.go"
	//facultymiddlerware "github.com/suhas-developer07/Smart-Attendence-System/server/internals/middlerwares/faculty_middlerware.go"
	studentmiddlerwarego "github.com/suhas-developer07/Smart-Attendence-System/server/internals/middlerwares/student_middlerware.go"
	"github.com/suhas-developer07/Smart-Attendence-System/server/internals/repository"

	attendence_service "github.com/suhas-developer07/Smart-Attendence-System/server/internals/service/attendence"
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

	attendanceService := attendence_service.NewAttendanceService(repo)
	attendanceHandler := attendance_handler.NewAttendanceHandler(attendanceService)


	//  Student 
	student := e.Group("/students")
	{
		student.POST("/register", studentHandler.StudentRegisterHandler)
		student.POST("/login", studentHandler.LoginStudentHandler)       
		student.PUT("/:student_id", studentHandler.UpdateStudentInfoHandler) 
		//student.GET("/:student_id", studentHandler.GetStudentByIDHandler)    
		student.GET("/subjects", subjectHandler.GetSubjectsByStudentIDHandler, studentmiddlerwarego.JWTMiddleware) 

	//  Subject 
	subject := e.Group("/subjects")
	{
		//all routes are working fine
		subject.POST("", subjectHandler.AddSubjectHandler)                        
		subject.GET("",subjectHandler.GetSubjectsByDeptAndSemHandler)            
		subject.GET("/faculty", subjectHandler.GetSubjectsByFacultyIDHandler,facultymiddlerware.FacultyJWTMiddleware) 
	}

	// Faculty
	faculty := e.Group("/faculty")
	{
		//all routes are working fine
		faculty.POST("/register", facultyHandler.RegisterFacultyHandler) 
		faculty.POST("/login", facultyHandler.AuthenticateFacultyHandler) 
		faculty.GET("/getfaculty", facultyHandler.GetFacultyByIDHandler,facultymiddlerware.FacultyJWTMiddleware) 
		faculty.GET("", facultyHandler.GetAllFacultyHandler)              
		faculty.GET("/department/:dept", facultyHandler.GetFacultyByDepartmentHandler) 
	}

	attendance := e.Group("/attendance")
	{
		//attendance.POST("", attendanceHandler.MarkAttendanceHandler)
		attendance.POST("/bulk", attendanceHandler.BulkAttendanceHandler,)//payload:http://localhost:8080/attendance/bulk
		attendance.GET("", attendanceHandler.GetAttendanceByStudentAndSubjectHandler,studentmiddlerwarego.JWTMiddleware)//http://localhost:8080/attendance?usn=4AL23IS059&subjectCode=1
		attendance.GET("/subject", attendanceHandler.GetAttendanceBySubjectAndDateHandler)//payload:http://localhost:8080/attendance/subject?subject_id=1&date=2025-09-18
		attendance.GET("/summary/subject/:subjectCode", attendanceHandler.GetAttendanceSummaryBySubjectHandler)//http://localhost:8080/attendance/summary/subject/1
		attendance.GET("/class", attendanceHandler.GetClassAttendanceHandler)//http://localhost:8080/attendance/class?subject_id=1&date=2025-09-18
		attendance.GET("/student/history", attendanceHandler.GetStudentAttendanceHistoryHandler,studentmiddlerwarego.JWTMiddleware)//http://localhost:8080/attendance/student/history?usn=4AL23IS059&subject_id=1
		attendance.POST("/assingnsubject",attendanceHandler.AssignSubjectToTimeRangeHandler)
		attendance.GET("/summary/student", attendanceHandler.GetAttendanceSummaryByStudentHandler,studentmiddlerwarego.JWTMiddleware)//http://localhost:8080/attendance/summary/student/4AL23IS059
	}

	// Health
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "server is healthy")
	})
}
}