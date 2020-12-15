package helper

import "testing"

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

// func TestInArray(t *testing.T) {

// }
// func TestInc(t *testing.T) {

// }
// func TestStrToSlice(t *testing.T) {

// }
