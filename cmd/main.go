package main

import (
	"cube"
	"fmt"
	"syscall/js"
	"time"
)

func solve(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return "Error: No argument provided"
	}

	message := args[0].String()
	start := time.Now()
	solution, err := cube.Solve(message)
	duration := time.Since(start)
	if err != nil {
		return err.Error()
	}
	return js.ValueOf(map[string]interface{}{
		"solution": solution,
		"duration": fmt.Sprintf("Solved in %s!", duration),
	})
}

func main() {
	js.Global().Set("solve", js.FuncOf(solve))

	select {}
}
