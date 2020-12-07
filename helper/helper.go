/*
	Package helper provides general functions that simplify process and/or perform validation check.
*/
package helper

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"time"
)

// Constants of date/time format.
const (
	DATELONGFORMAT   string = "2006-01-02"
	DATESHORTFORMAT  string = "20060102"
	TIME24FORMAT     string = "15:04"
	TIME24LONGFORMAT string = "15:04:05"
)

// CheckPanic simplifies panic recover process.
func CheckPanic() {
	if r := recover(); r != nil {
		fmt.Println("Something went wrong. Please report to administrator")
		fmt.Println(r.(error))
		fmt.Scanln()
		return
	}
}

// IsOverlapBetween checks whether the selected time is within start time and end time.
// Use time long format due the time comparison will not cater when both times are same e.g. 16:24.After(16:24) = false.
func IsOverlapBetween(selected string, start string, end string, check chan bool) {

	selectedTime, errSelected := time.Parse(TIME24LONGFORMAT, selected+":30")
	startTime, errStart := time.Parse(TIME24LONGFORMAT, start+":00")
	endTime, errEnd := time.Parse(TIME24LONGFORMAT, end+":59")

	if errSelected != nil || errStart != nil || errEnd != nil {
		log.Println("Invalid Time!!")
	}

	if selectedTime.After(startTime) && selectedTime.Before(endTime) {
		check <- true
	} else {
		check <- false
	}
}

// CheckValidDate checks whether the string is valid DATELONGFORMAT date.
func CheckValidDate(date string) bool {
	_, err := time.Parse(DATELONGFORMAT, date)
	if err != nil {
		return false
	}
	return true
}

// CheckValidTime checks whether the string is valid TIME24FORMAT time.
func CheckValidTime(timing string) bool {
	_, err := time.Parse(TIME24FORMAT, timing)
	if err != nil {
		return false
	}
	return true
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
