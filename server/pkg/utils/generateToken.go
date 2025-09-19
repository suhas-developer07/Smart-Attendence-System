package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("KALI_LINUX")
type StudentClaims struct {
	StudentID int64  `json:"student_id"`
	USN       string `json:"usn"`
	jwt.RegisteredClaims
}
func GenerateTokenForStudent(student_id int64 ,usn string) (string, error) {

	Claims := &StudentClaims{
		StudentID: student_id,
		USN:       usn,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "smart-attendence-system",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

type FacultyClaims struct {
	FacultyID int64  `json:"faculty_id"`
	Email     string `json:"email"`
	jwt.RegisteredClaims
}
func GenerateTokenForFaculty(faculty_id int64 ,email string) (string, error) {
	Claims := &FacultyClaims{
		FacultyID: faculty_id,
		Email:     email,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "smart-attendence-system",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}