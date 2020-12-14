package summary

import (
	"testing"
)

func TestCalculateRatio(t *testing.T) {
	var multipleTest = []struct {
		a, b int
		want float64
	}{
		{200, 100, 50},
		{10, 10, 100},
		{11, 3, 27.27},
		{0, 2, 0},
		{2, 0, 0},
		{0, 0, 0},
	}
	for _, td := range multipleTest {
		t.Run("TestCalculateRatio", func(t *testing.T) {
			ans := CalculateRatio(td.a, td.b)
			if ans != td.want {
				t.Errorf("CalculateRatio(%d,%d) result is not %f. %f is returned\n",
					td.a, td.b, td.want, ans)
			}
		})
	}
}
