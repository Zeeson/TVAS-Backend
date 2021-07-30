package customErrorFormat

import (
	"errors"
	"strings"
)

func FormatError(err string) error {

	if strings.Contains(err, "user_name") {
		return errors.New("UserName Already Taken")
	}

	if strings.Contains(err, "email") {
		return errors.New("Email Already Taken")
	}

	if strings.Contains(err, "name") {
		return errors.New("Name Already Taken")
	}

	if strings.Contains(err, "Account") {
		return errors.New("Account is locked!!")
	}

	if strings.Contains(err, "hashedPassword") {
		return errors.New("Incorrect Password")
	}
	return errors.New("Incorrect Details")
}