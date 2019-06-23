package distance

type mockCalculator struct {
	distance int
	err *error
}

// mock the google distance calculator
func (m *mockCalculator)Calculate(src []string, des []string) (int, error) {
	if m.err == nil {
		return m.distance, nil
	}

	return 0, *m.err
}

// initialize the mock calculator
func InitMockCalculator(d int, err error) {
	var mock mockCalculator

	if err == nil {
		mock.distance = d
	} else {
		mock.err = &err
	}

	calc = &mock
}

