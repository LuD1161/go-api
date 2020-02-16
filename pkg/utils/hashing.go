package utils

import "golang.org/x/crypto/bcrypt"

// Hash : Returns bcrypt hash of the plaintext string passed
func Hash(plaintext string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
}

// VerifyHash : Verifies string and their hash
func VerifyHash(hashedString, plaintext string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedString), []byte(plaintext))
}
