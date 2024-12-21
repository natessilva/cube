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
var coLookup [2187][moveCount]int
var eSliceP1Lookup [495][moveCount]int
var cpLookup [40320][moveCount]int
var udLookup [40320][moveCount]int
var eSliceP2Lookup [24][moveCount]int

// in Phase 1 we are trying to orient all pieces and put the 4 e slice
// edges in their slice (we don't care about where in the slice)
// We will calculate the min moves to solve EO + ESlice and
// CO + ESlice and we will use the max of the two as our heuristic
var phase1EoAndSlice [2048 * 495]byte
var phase1CoAndSlice [2187 * 495]byte

// In phase2 we are trying to permute everything. We will calculate
// the min moves to all all corners + all e slice edges and all edges
// and use the max of the two as our heuristic
var phase2CornerESliceMinMoves [40320 * 24]byte
var phase2AllEdgesMinMoves [40320 * 24]byte

func init() {
	initEO()
	initCO()
	initESliceP1()
	initCP()
	initUD()
	initESliceP2()
	initPhase2CornersESlice()
	initPhase2Edges()
	initPhase1COAndSlice()
	initPhase1EOAndSlice()
}

func initEO() {
	for i := 0; i < 2048; i++ {
		c := fromEOCoordinate(i)
		for j := 0; j < moveCount; j++ {
			movedCube := transform(c, moves[j])
			eoLookup[i][j] = toEOCoordinate(movedCube)
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
}

func initESliceP1() {
	for i := 0; i < 495; i++ {
		c := fromESliceP1Coordinate(i)
		for j := 0; j < moveCount; j++ {
			movedCube := transform(c, moves[j])
			eSliceP1Lookup[i][j] = toESliceP1Coordinate(movedCube)
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
}

func initUD() {
	for i := 0; i < 40320; i++ {
		c := fromUDCoordinate(i)
		for j := 0; j < moveCount; j++ {
			movedCube := transform(c, moves[j])
			udLookup[i][j] = toUDCoordinate(movedCube)
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
}

func initPhase1COAndSlice() {
	queue := []int{0}
	max := byte(0)
	for len(queue) > 0 {
		current := queue[0]
		depth := phase1CoAndSlice[current]
		queue = queue[1:]

		eSliceP1Coord := current % 495
		coCoord := current / 495

		for i := 0; i < moveCount; i++ {
			nextCo := coLookup[coCoord][i]
			nextESlice := eSliceP1Lookup[eSliceP1Coord][i]

			next := nextCo*495 + nextESlice
			if phase1CoAndSlice[next] == 0 && next != 0 {
				phase1CoAndSlice[next] = 1 + depth
				if phase1CoAndSlice[next] > max {
					max = phase1CoAndSlice[next]
				}
				queue = append(queue, next)
			}
		}
	}
}

func initPhase1EOAndSlice() {
	queue := []int{0}
	max := byte(0)
	for len(queue) > 0 {
		current := queue[0]
		depth := phase1EoAndSlice[current]
		queue = queue[1:]

		eSliceP1Coord := current % 495
		eoCoord := current / 495

		for i := 0; i < moveCount; i++ {
			nextCo := eoLookup[eoCoord][i]
			nextESlice := eSliceP1Lookup[eSliceP1Coord][i]

			next := nextCo*495 + nextESlice
			if phase1EoAndSlice[next] == 0 && next != 0 {
				phase1EoAndSlice[next] = 1 + depth
				if phase1EoAndSlice[next] > max {
					max = phase1EoAndSlice[next]
				}
				queue = append(queue, next)
			}
		}
	}
}

func initPhase2CornersESlice() {
	queue := []int{0}
	max := byte(0)
	for len(queue) > 0 {
		current := queue[0]
		depth := phase2CornerESliceMinMoves[current]
		queue = queue[1:]

		eSliceP2Coord := current % 24
		cpCoord := current / 24

		for i := 0; i < moveCount; i++ {
			// don't do quarter turns on the side faces as these are not valid moves in phase 2
			if i%3 != 1 && i > moveU3 && i < moveD {
				continue
			}
			nextCp := cpLookup[cpCoord][i]
			nextESlice := eSliceP2Lookup[eSliceP2Coord][i]

			next := nextCp*24 + nextESlice
			if phase2CornerESliceMinMoves[next] == 0 && next != 0 {
				phase2CornerESliceMinMoves[next] = 1 + depth
				if phase2CornerESliceMinMoves[next] > max {
					max = phase2CornerESliceMinMoves[next]
				}
				queue = append(queue, next)
			}
		}
	}
}

func initPhase2Edges() {
	queue := []int{0}
	max := byte(0)
	for len(queue) > 0 {
		current := queue[0]
		depth := phase2AllEdgesMinMoves[current]
		queue = queue[1:]

		eSliceP2Coord := current % 24
		udCoord := current / 24

		for i := 0; i < moveCount; i++ {
			// don't do quarter turns on the side faces as these are not valid moves in phase 2
			if i%3 != 1 && i > moveU3 && i < moveD {
				continue
			}
			nextCp := udLookup[udCoord][i]
			nextESlice := eSliceP2Lookup[eSliceP2Coord][i]

			next := nextCp*24 + nextESlice
			if phase2AllEdgesMinMoves[next] == 0 && next != 0 {
				phase2AllEdgesMinMoves[next] = 1 + depth
				if phase2AllEdgesMinMoves[next] > max {
					max = phase2AllEdgesMinMoves[next]
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
	coAndSlice := coCood*495 + ePermCoord
	eoAnSlice := eoCoord*495 + ePermCoord
	return int(math.Max(
		float64(phase1EoAndSlice[eoAnSlice]),
		float64(phase1CoAndSlice[coAndSlice]),
	))
}

// phase2Hueristic combines the CP, UD and ESliceP2 coorindate spaces
// to provide the lower bound in the same way we do in Phase 1.
func phase2Hueristic(cpCoord, eudCood, eeCoord int) int {
	cornerSliceCoord := cpCoord*24 + eeCoord
	edgesCoord := eudCood*24 + eeCoord
	return int(math.Max(
		float64(phase2CornerESliceMinMoves[cornerSliceCoord]),
		float64(phase2AllEdgesMinMoves[edgesCoord]),
	))
}
