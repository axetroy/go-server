package util

const (
	passwordPrefix = "gotest"
	passwordSuffix = "test"
)

func GeneratePassword(text string) string {
	password := MD5(passwordPrefix + text + passwordSuffix)
	return password
}
