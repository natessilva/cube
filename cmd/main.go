package main

import (
	"cube"
	"flag"
	"fmt"
)

func main() {
	var scramble string

	flag.StringVar(&scramble, "scramble", "", "Supply a scramble to solve")
	flag.Parse()

	solution, err := cube.Solve(scramble)
	if err != nil {
		panic(err)
	}
	fmt.Println(solution)
}
