package helper

import (
	"reflect"
	"testing"
)

func TestCheckPasswordStrength(t *testing.T) {
	testData := "1@2b3C4d"
	t.Run("TestCheckPasswordStrength", func(t *testing.T) {
		err := CheckPasswordStrength(testData)
		if err != nil {
			t.Errorf(err.Error())
		}
	})

	var multipleTest = []struct {
		pw, want string
	}{
		{"ABCDE FG", "Password cannot contains spaces!"},
		{"ABCDEFG", "Password requires eight or more characters"},
		{"ABCDEFGH", "Password requires at least one digit"},
		{"ABCDEFG8", "Password requires at least one lowercase character"},
		{"ABCDEFg8", "Password requires at least one special character"},
		{"@bcdefg8", "Password requires at least one uppercase character"},
	}
	for _, td := range multipleTest {
		t.Run("TestCheckPasswordStrengthError", func(t *testing.T) {
			err := CheckPasswordStrength(td.pw)
			if err != nil && err.Error() != td.want {
				t.Errorf("Password Check on \"%s\" return Incorrect Error \"%s\" instead of \"%s\"\n",
					td.pw, td.want, err.Error())
			}
		})
	}
}

func TestInArray(t *testing.T) {
	var collection []string = make([]string, 5)
	collection[0] = "Apple"
	collection[1] = "Boy"
	collection[2] = "Cat & Dog"
	collection[3] = "Elephant"
	collection[4] = "Flower"

	var multipleTest = []struct {
		target string
		found  bool
	}{
		{"apple", false},
		{"Boy", true},
		{"Cat & Dog", true},
		{"Cat", false},
		{"Gundam", false},
	}
	for _, td := range multipleTest {
		t.Run("TestInArray", func(t *testing.T) {
			isFound := InArray(td.target, collection)
			if isFound != td.found {
				t.Errorf("Incorrect return \"%t\" in InArray(\"%s\",%+v)\n",
					isFound, td.target, collection)
			}
		})
	}
}

func TestInc(t *testing.T) {
	var multipleTest = []struct {
		target, interval, want int
	}{
		{0, 1, 1},
		{0, 4, 4},
		{4, 4, 8},
		{16, 1, 17},
		{17, 3, 20},
	}

	for _, td := range multipleTest {
		t.Run("TestInc", func(t *testing.T) {
			result := Inc(td.target, td.interval)
			if result != td.want {
				t.Errorf("Incorrect return %d in Inc(%d,%d). Expecting %d\n",
					result, td.target, td.interval, td.want)
			}
		})
	}
}

func TestStrToSlice(t *testing.T) {
	var multipleTest = []struct {
		src       string
		delimiter string
		want      []string
	}{
		{"1,2,3,4,5", ",", []string{"1", "2", "3", "4", "5"}},
		{"apple|boy|cat|dog|elephant", "|", []string{"apple", "boy", "cat", "dog", "elephant"}},
	}
	for _, td := range multipleTest {
		t.Run("TestStrToSlice", func(t *testing.T) {
			result := StrToSlice(td.src, td.delimiter)
			//reflect.DeepEqual() return bool, poor performance but suitable for test case
			if !reflect.DeepEqual(result, td.want) {
				t.Errorf("Incorrect return \"%+v\" in StrToSlice(\"%s\",\"%s\"). Expecting \"%+v\"\n",
					result, td.src, td.delimiter, td.want)
			}
		})
	}
}
