package main

type OptionError struct{}

func (oe OptionError) Error() string {
	return "Ну ты внытуре фрик"
}

type ParseError struct {
	textError string
}

func (pe ParseError) Error() string {
	return pe.textError
}

type IntegralError struct{}

func (ie IntegralError) Error() string {
	return "integral has no solution to calculate it"
}
