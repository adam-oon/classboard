/*
	Package helper provides general functions that simplify process and/or perform validation check.
*/
package helper

import (
	"errors"
	"log"
	"regexp"
	"strings"
)

// CheckPanic simplifies panic recover process.
func CheckPanic() {
	if r := recover(); r != nil {
		log.Println("Something went wrong. Please report to administrator")
		log.Println(r.(error))
		return
	}
}

// CheckPasswordStrength checks against string with following condition: at least 8 chars, 1 digit, 1 lower, 1 upper and 1 special char.
func CheckPasswordStrength(password string) error {
	whitespace := `[\t\n\f\r ]`
	len := `.{8,}`
	digit := `[0-9]{1}`
	charLower := `[a-z]{1}`
	charTitle := `[A-Z]{1}`
	symbol := `[!@#~$%^&*()+|_]{1}`

	// prevent whitespace
	if matched, err := regexp.MatchString(whitespace, password); matched || err != nil {
		return errors.New("Password cannot contains spaces!")
	}
	if matched, err := regexp.MatchString(len, password); !matched || err != nil {
		return errors.New("Password requires eight or more characters")
	}
	if matched, err := regexp.MatchString(digit, password); !matched || err != nil {
		return errors.New("Password requires at least one digit")
	}
	if matched, err := regexp.MatchString(charLower, password); !matched || err != nil {
		return errors.New("Password requires at least one lowercase character")
	}
	if matched, err := regexp.MatchString(charTitle, password); !matched || err != nil {
		return errors.New("Password requires at least one uppercase character")
	}
	if matched, err := regexp.MatchString(symbol, password); !matched || err != nil {
		return errors.New("Password requires at least one special character")
	}
	return nil
}

// InArray checks whether the element(needle) is inside array(haystack).
func InArray(needle string, haystack []string) bool {
	for _, value := range haystack {
		if value == needle {
			return true
		}
	}
	return false
}

// Inc increments integer by interval set.
func Inc(target, interval int) int {
	return target + interval
}

// StrToSlice convert string sentence to slice.
func StrToSlice(src string, delimiter string) []string {
	choiceSlice := strings.Split(src, delimiter)
	return choiceSlice
}
