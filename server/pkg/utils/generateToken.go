package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("KALI_LINUX")
type Claims struct {
	StudentID int64  `json:"student_id"`
	USN       string `json:"usn"`
	jwt.RegisteredClaims
}
func GenerateTokenForStudent(student_id int64 ,usn string) (string, error) {

	Claims := &Claims{
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

func GenerateTokenForFaculty(faculty_id int64 ,email string) (string, error) {
	Claims := &Claims{
		StudentID: faculty_id,
		USN:       email,
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