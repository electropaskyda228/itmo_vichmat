package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type matrix struct {
	n       int
	a       [][]float64
	b       []float64
	x       []float64
	epsilon float64
	delta   []float64
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func workWithUser(in *bufio.Reader, m *matrix) {
	fmt.Println("Do you want to upload parametres from some file? (y|n)")
	var answer string
	fmt.Fscan(in, &answer)
	switch answer {
	case "y":
		fmt.Println("bruh")
		getNumbersFromFile(in, m)
	case "n":
		getNumbersFromConsole(in, m)
	default:
		fmt.Println("Wrong format of answer. Canceled.")
	}
}

func getOut() {
	os.Exit(1)
}

func getNumbersFromFile(in *bufio.Reader, m *matrix) {
	fmt.Print("Enter path to file: ")
	var pathToFile string
	fmt.Fscan(in, &pathToFile)

	file, err := os.ReadFile(pathToFile)
	if err != nil {
		fmt.Println("Can't find file. Canceled.")
		getOut()
	}
	rows := strings.Split(string(file), "\n")
	if len(rows) == 0 {
		fmt.Println("Wrong format of file. Canceled.")
		getOut()
	}
	str_n_e := strings.Fields(rows[0])
	var err1, err2 error
	m.n, err1 = strconv.Atoi(str_n_e[0])
	m.epsilon, err2 = strconv.ParseFloat(strings.TrimSpace(str_n_e[1]), 64)
	if err1 != nil || err2 != nil {
		fmt.Println("Can't find n and epsilon. Canceld.")
		getOut()
	}
	if len(rows) != m.n+1 {
		fmt.Println("Wrong format of matrix. Canceled.")
		getOut()
	}
	m.a = make([][]float64, m.n)
	m.b = make([]float64, m.n)
	for i := 0; i < m.n; i++ {
		row := strings.Fields(rows[i+1])
		if len(row) != m.n+1 {
			fmt.Println("Wrong format of matrix. Canceled.")
			getOut()
		}
		m.a[i] = make([]float64, m.n)
		for j := 0; j < m.n; j++ {
			tmp_a, err := strconv.ParseFloat(row[j], 64)
			if err != nil {
				fmt.Println("Wrong format of matrix. Canceled.")
				getOut()
			}
			m.a[i][j] = tmp_a
		}
		tmp_b, err := strconv.ParseFloat(row[m.n], 64)
		if err != nil {
			fmt.Println("Wrong format of matrix. Canceled.")
			getOut()
		}
		m.b[i] = tmp_b
	}
	m.x = make([]float64, m.n)
	for i := 0; i < m.n; i++ {
		m.x[i] = m.b[i] / m.a[i][i]
	}
}

func getNumbersFromConsole(in *bufio.Reader, m *matrix) {
	fmt.Print("Enter n: ")
	fmt.Fscan(in, &m.n)
	fmt.Println()
	m.a = make([][]float64, m.n)
	m.b = make([]float64, m.n)
	fmt.Println("Enter matrix in format ")
	fmt.Println("a_1_1 a_1_2 ... a_1_n b_1")
	fmt.Println("a_2_1 a_2_2 ... a_2_n b_2")
	fmt.Println("------...-------")
	fmt.Println("a_n_1 a_n_2 ... a_n_n b_n")
	for i := 0; i < m.n; i++ {
		m.a[i] = make([]float64, m.n)
		for j := 0; j < m.n; j++ {
			var tmp float64
			fmt.Fscan(in, &tmp)
			m.a[i][j] = tmp
		}
		var tmp_b float64
		fmt.Fscan(in, &tmp_b)
		m.b[i] = tmp_b
	}

	m.x = make([]float64, m.n)
	for i := 0; i < m.n; i++ {
		m.x[i] = m.b[i] / m.a[i][i]
	}

	fmt.Print("Enter accuracy: ")
	fmt.Fscan(in, &m.epsilon)
	fmt.Println()
}

func calculate_norma(m *matrix) float64 {
	var answer float64 = -1
	for i := 0; i < m.n; i++ {
		var s float64 = 0
		for j := 0; j < m.n; j++ {
			if i != j {
				s += abs(m.a[i][j])
			}
		}
		answer = max(answer, s/m.a[i][i])
	}
	return answer
}

func diff_calculate_norma(m *matrix, perm []int) float64 {
	var answer float64 = -1
	for i := 0; i < m.n; i++ {
		var s float64 = 0
		for j := 0; j < m.n; j++ {
			if j != i {
				s += abs(m.a[perm[i]][j])
			}
		}
		answer = max(answer, s/m.a[perm[i]][i])
	}
	return answer
}

func check(m *matrix) (int, float64) {
	var c float64 = calculate_norma(m)
	if c < 1 {
		return 1, c
	}
	perm := make([]int, m.n)
	for i := 0; i < m.n; i++ {
		perm[i] = i
	}
	for true {
		i := m.n - 2
		for i >= 0 && perm[i] >= perm[i+1] {
			i--
		}
		if i < 0 {
			break
		}
		j := m.n - 1
		for perm[j] <= perm[i] {
			j--
		}
		perm[i], perm[j] = perm[j], perm[i]
		left := i + 1
		right := m.n - 1
		for left < right {
			perm[left], perm[right] = perm[right], perm[left]
			left++
			right--
		}
		var tmp float64 = diff_calculate_norma(m, perm)
		if tmp < 1 {
			new_a := make([][]float64, m.n)
			new_b := make([]float64, m.n)
			for g := 0; g < m.n; g++ {
				new_a[g] = make([]float64, m.n)
				for l := 0; l < m.n; l++ {
					new_a[g][l] = m.a[perm[g]][l]
				}
				new_b[g] = m.b[perm[g]]
			}
			m.a = new_a
			m.b = new_b
			return 0, tmp
		}
	}
	return -1, c
}

func calculate(m *matrix, M int) (int, float64, bool) {
	result, norma := check(m)
	if result == -1 {
		return -1, norma, false
	}

	m.delta = make([]float64, m.n)

	var k int = 1
	var success bool = true
	for true {
		var delta float64 = 0
		for i := 0; i < m.n; i++ {
			var s float64 = 0
			for j := 0; j < m.n; j++ {
				if j != i {
					s += m.a[i][j] * m.x[j]
				}
			}
			var tmp_x float64 = (m.b[i] - s) / m.a[i][i]
			var d float64 = abs(m.x[i] - tmp_x)
			m.delta[i] = d
			delta = max(delta, d)
			m.x[i] = tmp_x
		}

		if delta < m.epsilon {
			break
		}
		if k == M {
			success = false
			break
		}
		k += 1
	}

	return k, norma, success

}

func showGoodAnswer(out *bufio.Writer, m *matrix, k int, c float64) {
	fmt.Fprintln(out, "Good news! Linear system has an answer.")
	fmt.Fprintln(out, "Vector x:", m.x)
	fmt.Fprintln(out, "Count of iterations:", k)
	fmt.Fprintln(out, "matrix's norm:", c)
	fmt.Fprintln(out, "Error vector:", m.delta)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	// ввод данных
	var m matrix = matrix{0, nil, nil, nil, 0, nil}
	workWithUser(in, &m)

	// рассчет
	var M int = 10000
	k, c, success := calculate(&m, M)

	if success {
		showGoodAnswer(out, &m, k, c)
	} else if k == -1 {
		fmt.Fprintln(out, "plaki plaki: there is no permutation of rows to get the suitable matrix")
	} else {
		fmt.Fprintln(out, "Not enough iterations", M, "to get the answer")
	}
}
