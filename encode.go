package cube

import (
	"errors"
	"strings"
)

func parseScrambe(s string) (cube, error) {
	movesMap := map[string]int{
		"U":  moveU,
		"U2": moveU2,
		"U'": moveU3,
		"L":  moveL,
		"L2": moveL2,
		"L'": moveL3,
		"F":  moveF,
		"F2": moveF2,
		"F'": moveF3,
		"R":  moveR,
		"R2": moveR2,
		"R'": moveR3,
		"B":  moveB,
		"B2": moveB2,
		"B'": moveB3,
		"D":  moveD,
		"D2": moveD2,
		"D'": moveD3,
	}
	s = strings.ReplaceAll(s, " ", "")
	result := cubeSolved

	for i := 0; i < len(s); i++ {
		if i+1 < len(s) && (s[i+1] == '2' || s[i+1] == '\'') {
			if move, exists := movesMap[s[i:i+2]]; exists {
				result = transform(result, moves[move])
				i++
			}
		} else if move, exists := movesMap[string(s[i])]; exists {
			result = transform(result, moves[move])

		} else {
			return cubeSolved, errors.New("invalid scamble")
		}
	}
	return result, nil
}

func toString(moves []int) string {
	moveStrings := [moveCount]string{
		"U",
		"U2",
		"U'",
		"L",
		"L2",
		"L'",
		"F",
		"F2",
		"F'",
		"R",
		"R2",
		"R'",
		"B",
		"B2",
		"B'",
		"D",
		"D2",
		"D'",
	}
	var b strings.Builder
	for i := 0; i < len(moves); i++ {
		b.WriteString(moveStrings[moves[i]])
		b.WriteString(" ")
	}
	return b.String()
}
