package main

import (
	"strconv"
	"strings"
)

type OptionError struct {
}

func (oe OptionError) Error() string {
	return "Ну ты внатуре фрик"
}

type ReadFileError struct {
	text string
}

func (rfe ReadFileError) Error() string {
	return rfe.text
}

type ReadError struct {
	text string
}

func (re ReadError) Error() string {
	return re.text
}

type SystemError struct {
}

func (se SystemError) Error() string {
	return "Нарушено условие использования метода простых итераций для решения системы нелинейных уравнений"
}

type ParseError struct {
	value string
}

func (pe ParseError) Error() string {
	return "Can not parse " + pe.value
}

type SecantError struct {
	koeff []float64
}

func (se SecantError) Error() string {
	return "There is no permition to enforce secant method to solve the equation with koeffs: " + massiveFloatToString(se.koeff)
}

type MultipleRootsError struct{}

func (mre MultipleRootsError) Error() string {
	return "In this izolation scope multiple roots of equation"
}

type SimpleIterationError struct {
	text string
}

func (sie SimpleIterationError) Error() string {
	return sie.text
}

type IterationError struct{}

func (ie IterationError) Error() string {
	return "Превышен лимит итераций"
}

func massiveFloatToString(numbers []float64) string {
	var strNumbers []string
	for _, num := range numbers {
		strNumbers = append(strNumbers, strconv.FormatFloat(num, 'f', -1, 64))
	}

	return strings.Join(strNumbers, ", ")
}
