package main

import (
	"fmt"
	"log"
	"os"

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
	cmd.SetupRoutes(e, Database)

	e.Logger.Fatal(e.Start(":8080"))
}
