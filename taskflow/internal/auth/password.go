// Package auth handles password hashing and JWT tokens.
package auth

import "golang.org/x/crypto/bcrypt"

// HashPassword returns a bcrypt hash of the password. bcrypt is deliberately
// slow and salts automatically, which is exactly what you want for passwords.
func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

// CheckPassword reports whether password matches the stored hash.
func CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
