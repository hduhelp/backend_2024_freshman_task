package user

import "regexp"

// 2 <= length <= 20
// ^[a-zA-Z0-9]+$
func validateUsername(username string) bool {
	if len(username) < 2 || len(username) > 20 {
		return false
	}

	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return re.MatchString(username)
}

// 8 <= length <= 30
// ^[a-zA-Z0-9\W_]+$
func validatePassword(password string) bool {
	if len(password) < 8 || len(password) > 30 {
		return false
	}

	re := regexp.MustCompile(`^[a-zA-Z0-9\W_]+$`)
	return re.MatchString(password)
}
