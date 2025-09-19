package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword hashes a plain text password
func HashPassword(password string) (string, error) {
    b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(b), nil
}

// ComparePassword compares plain text with a hashed password
func ComparePassword(hash, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
