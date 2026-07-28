package users

import "golang.org/x/crypto/bcrypt"

// hashPassword converts a plaintext password into a bcrypt hash suitable for
// storage. The plaintext password must never be persisted.
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// comparePassword reports whether password matches the stored bcrypt hash.
func comparePassword(passwordHash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
}
