package main

import (
	"fmt"
	"log"
	"os"
	 "github.com/labstack/echo/v4/middleware"


	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/suhas-developer07/Smart-Attendence-System/server/cmd"
)

func main() {
	fmt.Println("Server is running...")

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error Loading env file: %v", err)
	}

	DatabaseUrl := os.Getenv("DatabaseURL")

	if DatabaseUrl == "" {
		log.Fatalf("DatabaseUrl not found in env file ")
	}

	Database, err := cmd.ConnectToDB(DatabaseUrl)

	if err != nil {
		log.Fatalf("Failed to innitialize database.:%v", err)
	}

	defer Database.Close()

	

e := echo.New()

// Enable CORS for all requests
e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins: []string{"*"}, // allow all origins
    AllowMethods: []string{
        echo.GET,
        echo.POST,
        echo.PUT,
        echo.DELETE,
        echo.PATCH,
        echo.OPTIONS,
    },
    AllowHeaders: []string{
        echo.HeaderOrigin,
        echo.HeaderContentType,
        echo.HeaderAccept,
        echo.HeaderAuthorization,
    },
}))

cmd.SetupRoutes(e, Database)

e.Logger.Fatal(e.Start(":8080"))


	e.Logger.Fatal(e.Start(":8080"))
}
