package funcs // HashPassword ...
import (
	"golang.org/x/crypto/bcrypt"
)

//HashPassword ...
func HashPassword(pwd string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	return string(hash)
}

// ComparePassword --> Decrypt Password
func ComparePassword(hashpwd, pwd string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashpwd), []byte(pwd)); err != nil {
		return false
	}
	return true
}
