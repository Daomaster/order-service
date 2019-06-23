package distance

var calc Calculator

type Calculator interface {
	Calculate(src []string, des []string) (int, error)
}

// getter for the calculator
func GetCalculator() Calculator {
	return calc
}