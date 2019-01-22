package util

const (
	PasswordPrefix = "gotest"
	PasswordSuffix = "test"
)

func GeneratePassword(text string) string {
	password := MD5(PasswordPrefix + text + PasswordSuffix)
	return password
}
