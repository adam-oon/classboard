/*
	Package summary provides summary model that uses for reporting.
*/
package summary

import "math"

type Summary struct {
	StudentTotal  int
	QuestionTotal int
	Participation float64
	Correctness   float64
}

// CalculateRatio takes in 2 integers and calculate percentage ratio as 2 decimal float
func CalculateRatio(question_nums, answer_nums int) float64 {
	var participation float64
	qn := float64(question_nums)
	an := float64(answer_nums)

	if qn == 0 {
		return participation
	}

	participation = an / qn * 100

	// to 2 decimal place
	participation = math.Round(participation*100) / 100
	return participation
}
