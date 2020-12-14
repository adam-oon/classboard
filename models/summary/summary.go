package summary

type Summary struct {
	StudentTotal  int
	QuestionTotal int
	Participation float64
	Correctness   float64
}

func CalculateRatio(question_nums, answer_nums int) float64 {
	var participation float64
	qn := float64(question_nums)
	an := float64(answer_nums)

	if qn == 0 {
		return participation
	}

	participation = an / qn * 100
	return participation
}
