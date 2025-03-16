package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	fmt.Println("Выберите, какой интеграл хотите вычислить: ")
	fmt.Print(" 1) Определенный интеграл\n 2) Несобственный интеграл\n Enter: ")

	var option int
	ReadInt(in, &option, true)

	if option == 1 {
		SolveIntegral(in, out)
	} else if option == 2 {
		SolveInftyIntegral(in, out)
	} else {
		GetOut(OptionError{})
	}

}
