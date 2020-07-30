package utils

import "golang.org/x/crypto/bcrypt"

func GenerateFromPassword(psssword string) (bcryptpwd string, err error) {
	bcrypeBytes, err := bcrypt.GenerateFromPassword([]byte(psssword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bcrypeBytes), nil
}

func CheckPassword(password, bcryptPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(bcryptPassword), []byte(password))
	if err != nil {
		return false
	} else {
		return true
	}
}
