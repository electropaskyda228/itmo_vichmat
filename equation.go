package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Определение структур уравнения и ошибок
type equation struct {
	koeff       []float64
	hasSolution bool
	a           float64
	b           float64
	accuracy    float64
	itera       int
}

func checkMultipleError(eq equation) bool {
	f := getFirstDerivative(eq)
	may := f(eq.a) > 0
	for x := eq.a; x <= eq.b; x += eq.accuracy {
		if (f(x) > 0 && !may) || (f(x) < 0 && may) {
			return false
		}
	}
	f = getSourceFunction(eq)
	return f(eq.a)*f(eq.b) < 0
}

// Функция ввода коеффициентов уравнения
func getKoeff(in *bufio.Reader) ([]float64, error) {
	fmt.Print("Введите коэффициенты уравнения в порядке возрастания степеней: ")
	rowKoeff, prefix, err := in.ReadLine()
	if err != nil {
		return nil, ReadError{"Can not read the line of koeffs"}
	}
	if prefix {
		return nil, nil
	}
	stringKoeff := strings.Split(string(rowKoeff), " ")
	answer := make([]float64, len(stringKoeff))
	for index, number := range stringKoeff {
		value, err := strconv.ParseFloat(number, 64)
		if err != nil {
			return nil, ParseError{number}
		}
		answer[index] = value
	}
	return answer, nil
}

// Функция чтения границ изоляции корня
func getBorder(in *bufio.Reader) (float64, float64) {
	fmt.Print("Введите два числа: левую и правую границу изоляции корня: ")
	var a, b float64
	ReadFloat(in, &a, false, "left border izolation of x")
	ReadFloat(in, &b, true, "right border izolation of y")
	return a, b
}

// Функция чтения точности вычислений
func getAccuracy(in *bufio.Reader) float64 {
	fmt.Print("Введите точность вычислений: ")
	var epsilon float64
	ReadFloat(in, &epsilon, true, "accuracy")
	return epsilon
}

// Получение второй производной исходной функции
func getSecondDerivative(eq equation) func(x float64) float64 {
	if len(eq.koeff) < 3 {
		return nil
	}
	newKoeff := make([]float64, len(eq.koeff)-2)
	for i := 2; i < len(eq.koeff); i++ {
		newKoeff[i-2] = float64((i-1)*i) * eq.koeff[i]
	}
	return func(x float64) float64 {
		var answer float64 = 0
		for index, koeff := range newKoeff {
			answer += koeff * FastPow(x, index)
		}
		return answer
	}
}

// Получение первой производной исходной функции
func getFirstDerivative(eq equation) func(x float64) float64 {
	if len(eq.koeff) < 2 {
		return nil
	}
	newKoeff := make([]float64, len(eq.koeff)-1)
	for i := 1; i < len(eq.koeff); i++ {
		newKoeff[i-1] = float64(i) * eq.koeff[i]
	}
	return func(x float64) float64 {
		var answer float64 = 0
		for index, koeff := range newKoeff {
			answer += koeff * FastPow(x, index)
		}
		return answer
	}
}

// Получение исходной функции
func getSourceFunction(eq equation) func(x float64) float64 {
	return func(x float64) float64 {
		var answer float64 = 0
		for index, koeff := range eq.koeff {
			answer += koeff * FastPow(x, index)
		}
		return answer
	}
}

// Главный метод хорд
func methodChord(eq equation) (float64, int, error) {
	secondDerivative := getSecondDerivative(eq)
	f := getSourceFunction(eq)
	var x float64
	var k int
	if secondDerivative == nil {
		x, k = methodChordDefault(eq)
	} else if f(eq.a)*secondDerivative(eq.a) > 0 {
		x, k = methodChordFixLeftBorder(eq)
	} else if f(eq.b)*secondDerivative(eq.b) > 0 {
		x, k = methodChordFixRightBorder(eq)
	} else {
		x, k = methodChordDefault(eq)
	}
	if k == eq.itera {
		return 0, 0, IterationError{}
	}
	return x, k, nil
}

// Обычный метод хорд (без фиксации границ)
func methodChordDefault(eq equation) (float64, int) {
	f := getSourceFunction(eq)
	var x float64 = eq.a - (eq.b-eq.a)/(f(eq.b)-f(eq.a))*f(eq.a)
	var k int = 0
	for Abs(f(x)) >= eq.accuracy && k < eq.itera {
		var f_a, f_x float64 = f(eq.a), f(x)
		if f_a*f_x <= 0 {
			eq.b = x
		} else {
			eq.a = x
		}
		x = eq.a - (eq.b-eq.a)/(f(eq.b)-f(eq.a))*f(eq.a)
		k += 1
	}
	return x, k
}

// Метод хорд с фиксацией левой границы
func methodChordFixLeftBorder(eq equation) (float64, int) {
	f := getSourceFunction(eq)
	var x float64 = eq.b
	var k int = 0
	for Abs(f(x)) >= eq.accuracy && k < eq.itera {
		x = x - (eq.a-x)/(f(eq.a)-f(x))*f(x)
		k++
	}
	return x, k
}

// Метод хорд с фиксацией правой границы
func methodChordFixRightBorder(eq equation) (float64, int) {
	f := getSourceFunction(eq)
	var x float64 = eq.a
	var k int = 0
	for Abs(f(x)) >= eq.accuracy && k < eq.itera {
		x = x - (eq.b-x)/(f(eq.b)-f(x))*f(x)
	}
	return x, k
}

// Главный метод секущих
func methodSecant(eq equation) (float64, int, error) {
	f := getSourceFunction(eq)
	secondDerivative := getSecondDerivative(eq)
	var tmp float64
	if f(eq.a)*secondDerivative(eq.a) > 0 {
		tmp = eq.a
	} else if f(eq.b)*secondDerivative(eq.b) > 0 {
		tmp = eq.b
	} else {
		return 0, 0, SecantError{eq.koeff}
	}
	answer, itera, err := methodSecantCalculate(eq, tmp)
	if err != nil {
		return 0, 0, err
	}
	return answer, itera, nil
}

// Функция по вычислению по методу секущих
func methodSecantCalculate(eq equation, start float64) (float64, int, error) {
	f := getSourceFunction(eq)
	var xPrev, xPrevPrev float64 = start + eq.accuracy, start
	var x float64 = xPrev - (xPrev-xPrevPrev)/(f(xPrev)-f(xPrevPrev))*f(xPrev)
	var k int = 0
	for Abs(f(x)) >= eq.accuracy && k < eq.itera {
		xPrevPrev = xPrev
		xPrev = x
		x = xPrev - (xPrev-xPrevPrev)/(f(xPrev)-f(xPrevPrev))*f(xPrev)
		k++
	}
	if k == eq.itera {
		return 0, 0, IterationError{}
	}
	return x, k, nil
}

// Главный метод простых итераций
func methodSimpleItaration(eq equation) (float64, int, error) {
	lambda, err := getLambda(eq)
	if err != nil {
		return 0, 0, err
	}
	phiEq := getNewFunction(eq, lambda)
	if len(phiEq.koeff) > 2 {
		phiEq.koeff[1] += 1
	} else {
		phiEq.koeff = append(phiEq.koeff, 1)
	}
	phiFirstDerivative := getFirstDerivative(phiEq)
	if Abs(phiFirstDerivative(eq.a)) >= 1 || Abs(phiFirstDerivative(eq.b)) >= 1 {
		return 0, 0, SimpleIterationError{"There is no permition to use method of simple iterations (problem with phi)"}
	}
	var x1 float64 = (eq.a + eq.b) / 2
	var x2 float64 = x1 + eq.accuracy + 1
	phi := getSourceFunction(phiEq)
	var k int = 0
	for Abs(x2-x1) >= eq.accuracy && k < eq.itera {
		x2 = x1
		x1 = phi(x1)
		k++
	}
	if k == eq.itera {
		return 0, 0, IterationError{}
	}
	return x1, k, nil
}

// Получение коэффициента лямбда
func getLambda(eq equation) (float64, error) {
	firstDerivative := getFirstDerivative(eq)
	if firstDerivative == nil {
		return 0, SimpleIterationError{"There is no permition to use method of simple iterations (first derivative is 0)"}
	}
	i := eq.a
	var maximum float64 = -1
	var isPositive bool = firstDerivative(i) > 0 || firstDerivative(i+eq.accuracy) > 0
	for i < eq.b {
		if (firstDerivative(i) > 0 && !isPositive) || (firstDerivative(i) < 0 && isPositive) {
			return 0, SimpleIterationError{"There is no permition to use method of simple itaerations (different signs of first derivative)"}
		}
		maximum = max(Abs(firstDerivative(i)), maximum)
		i += eq.accuracy / 10
	}
	if isPositive {
		return -1 / maximum, nil
	}
	return 1 / maximum, nil
}

// Получение новой функции умножением ее на коэффициент
func getNewFunction(eq equation, lambda float64) equation {
	newEq := eq
	newEq.koeff = make([]float64, len(eq.koeff))
	for index, number := range eq.koeff {
		newEq.koeff[index] = number * lambda
	}
	return newEq
}

// Взять данные для уравнения с консоли
func getInfoFromConsole(in *bufio.Reader, eq *equation) {
	eq.a, eq.b = getBorder(in)
	eq.accuracy = getAccuracy(in)
}

// Взять данные для уравнения из файла
func getInfoFromFile(in *bufio.Reader, eq *equation) {
	fmt.Print("Enter path to file: ")
	var pathToFile string
	fmt.Fscan(in, &pathToFile)
	in.ReadLine()

	file, err := os.ReadFile(pathToFile)
	if err != nil {
		GetOut(ReadFileError{"Can't find file. Canceled."})
	}
	rows := strings.Split(string(file), "\n")
	if len(rows) == 0 {
		GetOut(ReadFileError{"Wrong format of file. Canceled."})
	}
	str_n_e := strings.Fields(rows[0])
	var err1, err2 error
	eq.a, err1 = strconv.ParseFloat(strings.TrimSpace(str_n_e[0]), 64)
	eq.b, err2 = strconv.ParseFloat(strings.TrimSpace(str_n_e[1]), 64)
	if err1 != nil || err2 != nil {
		GetOut(ReadFileError{"Can't find border. Canceld."})
	}
	if len(rows) != 2 {
		GetOut(ReadFileError{"Wrong format of file. Canceled."})
	}
	str_e := strings.Fields(rows[1])
	eq.accuracy, err1 = strconv.ParseFloat(strings.TrimSpace(str_e[0]), 64)
	if err1 != nil {
		GetOut(ReadFileError{"Can't find accuracy. Canceld."})
	}
}

// Запуск программы по решению нелинейных уравнений
func LinearEquation(in *bufio.Reader, out *bufio.Writer) {
	koeff, err := getKoeff(in)
	if err != nil {
		GetOut(err)
	}
	var eq equation = equation{koeff, false, 0, 0, 0, 1000000}
	fmt.Print("Выберете, как ввести данные\n 1) Файл\n 2) Вручную\n Enter: ")
	var option int
	ReadInt(in, &option, true)
	if option == 1 {
		getInfoFromFile(in, &eq)
	} else if option == 2 {
		getInfoFromConsole(in, &eq)
	} else {
		GetOut(OptionError{})
	}

	if !checkMultipleError(eq) {
		GetOut(MultipleRootsError{})
	}

	DrawSingleFunction(getSourceFunction(eq), eq.a-(eq.b-eq.a)/10, eq.b+(eq.b-eq.a)/10, eq.accuracy)
	f := getSourceFunction(eq)

	answer, itera, err := methodChord(eq)
	if err != nil {
		GetOut(err)
	}
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "Корень, полученный методом хорд: ", answer)
	fmt.Fprintln(out, "Значение функции в данной точке: ", f(answer))
	fmt.Fprintln(out, "Количество итераций: ", itera)
	fmt.Fprintln(out, "")

	anotherAnswer, itera, err := methodSecant(eq)
	if err != nil {
		GetOut(err)
	}
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "Корень, полученный методом секущих: ", anotherAnswer)
	fmt.Fprintln(out, "Значение функции в данной точке: ", f(anotherAnswer))
	fmt.Fprintln(out, "Количество итераций: ", itera)
	fmt.Fprintln(out, "")

	aAnotherAnswer, itera, err := methodSimpleItaration(eq)
	if err != nil {
		GetOut(err)
	}
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "Корень, полученный методом простых итераций", aAnotherAnswer)
	fmt.Fprintln(out, "Значение функции в данной точке: ", f(anotherAnswer))
	fmt.Fprintln(out, "Количество итераций: ", itera)
	fmt.Fprintln(out, "")
}
