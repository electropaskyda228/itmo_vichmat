package main

import (
	"bufio"
	"fmt"
	"math"
)

type integral struct {
	f        func(x float64) float64
	a        float64
	b        float64
	accuracy float64
}

func solveWithRunge(I integral, f func(I integral, n int) float64) (float64, int) {
	I0 := f(I, 4)
	I1 := f(I, 8)
	var n int = 8
	for Abs(I0-I1) >= I.accuracy {
		n *= 2
		I0 = I1
		I1 = f(I, n)
	}
	return I1, n
}

func methodSquareMiddle(I integral, n int) float64 {
	var h, summa float64 = (I.b - I.a) / float64(n), 0
	var x float64 = I.a + h/2
	for ; x <= I.b; x += h {
		summa += I.f(x)
	}
	return summa * h
}

func methodSquareLeft(I integral, n int) float64 {
	var h, summa float64 = (I.b - I.a) / float64(n), 0
	var x float64 = I.a
	for ; x < I.b; x += h {
		summa += I.f(x)
	}
	return summa * h
}

func methodSquareRight(I integral, n int) float64 {
	var h, summa float64 = (I.b - I.a) / float64(n), 0
	var x float64 = I.a + h
	for ; x <= I.b; x += h {
		summa += I.f(x)
	}
	return summa * h
}

func methodTrapezoid(I integral, n int) float64 {
	var h, summa float64 = (I.b - I.a) / float64(n), 0
	var x float64 = I.a + h
	for ; x < I.b; x += h {
		summa += I.f(x)
	}
	return (summa*2 + I.f(I.a) + I.f(I.b)) * h / 2
}

func methodSimpson(I integral, n int) float64 {
	var h, summaA, summaB float64 = (I.b - I.a) / float64(n), 0, 0
	var x float64 = I.a + h
	for ; x < I.b; x += 2 * h {
		summaA += I.f(x)
	}
	x = I.a + 2*h
	for ; x < I.b; x += 2 * h {
		summaB += I.f(x)
	}
	var summa float64 = 4*summaA + 2*summaB + I.f(I.a) + I.f(I.b)
	return summa * h / 3
}

func getIntegralData(in *bufio.Reader, I *integral) {
	fmt.Print("Введите границы интегрирования (start end): ")
	var a, b float64
	ReadFloat(in, &a, false, "left border is incorrect")
	ReadFloat(in, &b, true, "right border is incorrect")
	I.a, I.b = a, b

	fmt.Print("Введите точность вычислений: ")
	var epsilon float64
	ReadFloat(in, &epsilon, true, "accuracy is incorrect")
	I.accuracy = epsilon
}

func chooseIntegral(in *bufio.Reader, I *integral) {
	fmt.Println("Выберите функцию, определенный интеграл которой вы хотите вычислить:")
	fmt.Print(" 1) 2x^3 - 2x^2 + 7x - 14\n 2) ln(x)\n 3) arctg(x)\n Enter: ")
	var option int
	ReadInt(in, &option, true)

	if option == 1 {
		I.f = func(x float64) float64 {
			return 2*x*x*x - 2*x*x + 7*x - 14
		}
	} else if option == 2 {
		I.f = func(x float64) float64 {
			return math.Log(x)
		}
	} else if option == 3 {
		I.f = func(x float64) float64 {
			return math.Atan(x)
		}
	} else {
		GetOut(OptionError{})
	}
}

func chooseInftyIntegral(in *bufio.Reader, I *integral) float64 {
	fmt.Println("Выберите функцию, несобственный интеграл которой вы хотите вычислить:")
	fmt.Print(" 1) 1/sqrt(x) на [0; 2]\n 2) 1/(1-x) на [0; 1]\n 3) 1/x^3 на [-2; 3]\n Enter: ")
	var option int
	ReadInt(in, &option, true)

	if option == 1 {
		I.f = func(x float64) float64 {
			return 1 / math.Sqrt(x)
		}
		I.a = 1
		I.b = 2
		return 2
	} else if option == 2 {
		I.f = nil
		return 0
	} else if option == 3 {
		I.f = func(x float64) float64 {
			return 1 / (x * x * x)
		}
		I.a = 2
		I.b = 3
		return 0
	} else {
		GetOut(OptionError{})
	}
	return 0
}

func SolveIntegral(in *bufio.Reader, out *bufio.Writer) {
	I := integral{nil, 2, 4, 0.001}

	chooseIntegral(in, &I)
	getIntegralData(in, &I)

	fmt.Println("Выберите метод, которым хотите вычислить интеграл:")
	fmt.Println(" 1) Метод левых прямоугольников\n 2) Метод правых прямоугольников\n 3) Метод серединных прямоугольников\n 4) Метод трапеций\n 5) Метод Симпсона")
	fmt.Print("Enter: ")

	var option int
	ReadInt(in, &option, true)

	var answer float64
	var itera int
	if option == 1 {
		answer, itera = solveWithRunge(I, methodSquareLeft)
	} else if option == 2 {
		answer, itera = solveWithRunge(I, methodSquareRight)
	} else if option == 3 {
		answer, itera = solveWithRunge(I, methodSquareMiddle)
	} else if option == 4 {
		answer, itera = solveWithRunge(I, methodTrapezoid)
	} else if option == 5 {
		answer, itera = solveWithRunge(I, methodSimpson)
	} else {
		GetOut(OptionError{})
	}

	fmt.Fprintln(out, "Ответ:", answer, "  количество разбиений отрезка:", itera)
}

func SolveInftyIntegral(in *bufio.Reader, out *bufio.Writer) {
	I := integral{nil, 2, 4, 0.001}

	var toAdd float64 = chooseInftyIntegral(in, &I)
	if I.f == nil {
		GetOut(IntegralError{})
	}

	fmt.Print("Введите точность вычислений: ")
	var epsilon float64
	ReadFloat(in, &epsilon, true, "accuracy is incorrect")
	I.accuracy = epsilon

	answer, itera := solveWithRunge(I, methodSimpson)
	fmt.Fprintln(out, "Ответ:", answer+toAdd, "  количество разбиений отрезка:", itera)

}
