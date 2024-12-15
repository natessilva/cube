package cube

import "math"

// a bunch of lookup tables for each of the coordinate
// spaces. We have a lookup table of each coordinate in
// the space and the coordinate after applying every possible
// move. These tables represet a graph and using a DFS
// we are able to also fill in a lookup tables of exactly
// the minimun number of moves required to get to coordinate
// 0 for every other coordinate.
// These lookup tables are used during searching to "apply" moves
// rather than converting cubes to and from the various coordinate
// spaces, we just use the lookup tables for moves to navigate around
// the graph during search. We use the min moves tables as a hueristic
// to guide the actual search algorithm towards good solutions.
var eoLookup [2048][moveCount]int
var eoMinMoves [2048]int
var coLookup [2187][moveCount]int
var coMinMoves [2187]int
var eSliceP1Lookup [495][moveCount]int
var eSliceP1MinMoves [495]int
var cpLookup [40320][moveCount]int
var cpMinMoves [40320]int
var udLookup [40320][moveCount]int
var udMinMoves [40320]int
var eSliceP2Lookup [24][moveCount]int
var eSliceP2MinMoves [24]int

func init() {
	initEO()
	initCO()
	initESliceP1()
	initCP()
	initUD()
	initESliceP2()
}

func initEO() {
	for i := 0; i < 2048; i++ {
		c := fromEOCoordinate(i)
		for j := 0; j < moveCount; j++ {
			movedCube := transform(c, moves[j])
			eoLookup[i][j] = toEOCoordinate(movedCube)
		}
	}
	queue := []int{0}
	max := 0
	for len(queue) > 0 {
		current := queue[0]
		depth := eoMinMoves[current]
		queue = queue[1:]
		for i := 0; i < moveCount; i++ {
			next := eoLookup[current][i]
			if eoMinMoves[next] == 0 && next != 0 {
				eoMinMoves[next] = 1 + depth
				if eoMinMoves[next] > max {
					max = eoMinMoves[next]
				}
				queue = append(queue, next)
			}
		}
	}
}

func initCO() {
	for i := 0; i < 2187; i++ {
		c := fromCOCoordinate(i)
		for j := 0; j < moveCount; j++ {
			movedCube := transform(c, moves[j])
			coLookup[i][j] = toCOCoordinate(movedCube)
		}
	}
	queue := []int{0}
	max := 0
	for len(queue) > 0 {
		current := queue[0]
		depth := coMinMoves[current]
		queue = queue[1:]
		for i := 0; i < moveCount; i++ {
			next := coLookup[current][i]
			if coMinMoves[next] == 0 && next != 0 {
				coMinMoves[next] = 1 + depth
				if coMinMoves[next] > max {
					max = coMinMoves[next]
				}
				queue = append(queue, next)
			}
		}
	}
}

func initESliceP1() {
	for i := 0; i < 495; i++ {
		c := fromESliceP1Coordinate(i)
		for j := 0; j < moveCount; j++ {
			movedCube := transform(c, moves[j])
			eSliceP1Lookup[i][j] = toESliceP1Coordinate(movedCube)
		}
	}
	queue := []int{0}
	max := 0
	for len(queue) > 0 {
		current := queue[0]
		depth := eSliceP1MinMoves[current]
		queue = queue[1:]
		for i := 0; i < moveCount; i++ {
			next := eSliceP1Lookup[current][i]
			if eSliceP1MinMoves[next] == 0 && next != 0 {
				eSliceP1MinMoves[next] = 1 + depth
				if eSliceP1MinMoves[next] > max {
					max = eSliceP1MinMoves[next]
				}
				queue = append(queue, next)
			}
		}
	}
}

func initCP() {
	for i := 0; i < 40320; i++ {
		c := fromCPCoordinate(i)
		for j := 0; j < moveCount; j++ {
			movedCube := transform(c, moves[j])
			cpLookup[i][j] = toCPCoordinate(movedCube)
		}
	}
	queue := []int{0}
	max := 0
	for len(queue) > 0 {
		current := queue[0]
		depth := cpMinMoves[current]
		queue = queue[1:]
		for i := 0; i < moveCount; i++ {
			next := cpLookup[current][i]
			if cpMinMoves[next] == 0 && next != 0 {
				cpMinMoves[next] = 1 + depth
				if cpMinMoves[next] > max {
					max = cpMinMoves[next]
				}
				queue = append(queue, next)
			}
		}
	}
}

func initUD() {
	for i := 0; i < 40320; i++ {
		c := fromUDCoordinate(i)
		for j := 0; j < moveCount; j++ {
			movedCube := transform(c, moves[j])
			udLookup[i][j] = toUDCoordinate(movedCube)
		}
	}
	queue := []int{0}
	max := 0
	for len(queue) > 0 {
		current := queue[0]
		depth := udMinMoves[current]
		queue = queue[1:]
		for i := 0; i < moveCount; i++ {
			next := udLookup[current][i]
			if udMinMoves[next] == 0 && next != 0 {
				udMinMoves[next] = 1 + depth
				if udMinMoves[next] > max {
					max = udMinMoves[next]
				}
				queue = append(queue, next)
			}
		}
	}
}

func initESliceP2() {
	for i := 0; i < 24; i++ {
		c := fromESliceP2Coordinate(i)
		for j := 0; j < moveCount; j++ {
			movedCube := transform(c, moves[j])
			eSliceP2Lookup[i][j] = toESliceP2Coordinate(movedCube)
		}
	}
	queue := []int{0}
	max := 0
	for len(queue) > 0 {
		current := queue[0]
		depth := eSliceP2MinMoves[current]
		queue = queue[1:]
		for i := 0; i < moveCount; i++ {
			next := eSliceP2Lookup[current][i]
			if eSliceP2MinMoves[next] == 0 && next != 0 {
				eSliceP2MinMoves[next] = 1 + depth
				if eSliceP2MinMoves[next] > max {
					max = eSliceP2MinMoves[next]
				}
				queue = append(queue, next)
			}
		}
	}
}

// phase1Hueristic combines the min moves for the EO space
// the CO space and the E Slice Phase 1 spaces. We are attempting
// to solve all three spaces simultaneously and we know the optimatl
// number of moves required to solve each indepentantly. It follow from
// there that the max of the three mins represents a lower bound for
// possible solves - any branch in the search space that would cost
// more than that lower bound is a branch not worth searching and can
// be ignored or "pruned" from the search space.
func phase1Hueristic(eoCoord, coCood, ePermCoord int) int {
	return int(math.Max(
		math.Max(
			float64(eoMinMoves[eoCoord]),
			float64(coMinMoves[coCood]),
		),
		float64(eSliceP1MinMoves[ePermCoord]),
	))
}

// phase2Hueristic combines the CP, UD and ESliceP2 coorindate spaces
// to provide the lower bound in the same way we do in Phase 1.
func phase2Hueristic(cpCoord, eudCood, eeCoord int) int {
	return int(math.Max(
		math.Max(
			float64(cpMinMoves[cpCoord]),
			float64(udMinMoves[eudCood]),
		),
		float64(eSliceP2MinMoves[eeCoord]),
	))
}
