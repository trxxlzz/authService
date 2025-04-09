package utils

import "golang.org/x/crypto/bcrypt"

func VerifyPassword(hashedPassword string, candidatePassword string) bool {
	//проверяем хешированный пароль из бд с тем, который пришел
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
	return err == nil
}
