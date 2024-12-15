package cube

// Edge Orientation coordinate space

// toEOCordinate converts a cube structure into a number
// from 0-2047. There are 2^11 possible edge orientations.
// the orientation of the last edge is determined by the
// rest
func toEOCoordinate(c cube) int {
	result := 0
	for _, e := range c.Edges[0 : edgeCount-1] {
		result = (result << 1) | (e.Orientation & 1)
	}
	return result
}

// fromEOCoordinate turns a coordinate into a Cube. Importantly,
// it does not produce an actually valid Cube - just valid from
// the perspective of EO. This means that this Cube is not
// suitable for anything other than generating pruning tables.
func fromEOCoordinate(c int) cube {
	var result cube
	totalOrientation := 0
	for i := edgeCount - 2; i >= 0; i-- {
		result.Edges[i].Orientation = c & 1
		totalOrientation += result.Edges[i].Orientation
		c >>= 1
	}
	if totalOrientation%2 != 0 {
		result.Edges[edgeCount-1].Orientation = 1
	}
	return result
}

// toCOCoordinate converts a Cube into the corner orientation
// coordinate space. There are 3^7 possible corner orientations.
// The orientation of the last corner is determined by the rest.
// The total of all corner orientations must be divisible by 3.
func toCOCoordinate(c cube) int {
	result := 0
	for _, c := range c.Corners[0 : cornerCount-1] {
		result = result*3 + c.Orientation
	}
	return result
}

// fromCOCoordinate converts a coordinate in the CO space into
// a Cube. This Cube is only valid from the perspective of CO and
// is not suitable for anything other than generating pruning tables.
func fromCOCoordinate(c int) cube {
	var result cube
	totalOrientation := 0
	for i := cornerCount - 2; i >= 0; i-- {
		result.Corners[i].Orientation = c % 3
		totalOrientation += result.Corners[i].Orientation
		c /= 3
	}
	if totalOrientation%3 != 0 {
		result.Corners[cornerCount-1].Orientation = 3 - totalOrientation%3
	}
	return result
}

// toESliceP1Coordinate converts a cube into the E slice coordinate
// space for phase 1. This space represents 12 choose 4, or 495 possible
// configurations of the 4 edge slice edges, order independant. The solved
// state of this problem space is when all 4 edges of the E slice are
// in the E slice, but not necessarily solved yet.
func toESliceP1Coordinate(c cube) int {
	result := 0

	k := 1

	for i := edgeCount - 1; i >= 0; i-- {
		if c.Edges[i].Index >= edgeFL && c.Edges[i].Index <= edgeBL {
			result += cNK(edgeCount-1-i, k)
			k += 1
		}
	}
	return result
}

// fromESliceP1Coordinate creates a cube from an E slice coordinate
// this cube will have the 4 E slice edges taking up the 4 specific
// in any order. When c = 0, the will generate a cube with all 4 E slice
// edges in the E slice, but potentially out of order.
func fromESliceP1Coordinate(c int) cube {
	var result cube

	// i := 0

	for k := 4; k >= 1; k-- {
		for i := 0; i < edgeCount; i++ {
			n := edgeCount - 1 - i
			coefficient := cNK(n, k)
			if coefficient <= c {
				result.Edges[edgeCount-1-n].Index = k - 1 + edgeFL
				c -= coefficient
				break
			}
		}
	}
	return result
}

// cNK computes binomial coefficent N choose K
func cNK(n, k int) int {
	if k > n {
		return 0
	}
	result := 1
	for i := n; i > n-k; i-- {
		result *= i
	}
	for i := k; i >= 2; i-- {
		result /= i
	}
	return result
}

// toCPCoordinate converts a Cube into the corner permutatation coordinate
// space. There are 8! possible corner permutations, where 0 is the solved
// state of all corner permutations
func toCPCoordinate(c cube) int {
	result := 0
	for i := 0; i < cornerCount; i++ {
		count := 0
		for j := i + 1; j < cornerCount; j++ {
			if c.Corners[j].Index < c.Corners[i].Index {
				count++
			}
		}
		result += count * factorial(cornerCount-i-1)
	}
	return result
}

// fromCPCoordinate converts a coordinate in the CP space back into a Cube.
// This cube is only valid for generating lookup tables for CP.
func fromCPCoordinate(c int) cube {
	var result cube

	available := make([]int, cornerCount)
	for i := 0; i < cornerCount; i++ {
		available[i] = i
	}

	for i := 0; i < cornerCount; i++ {
		fact := factorial(cornerCount - i - 1)

		pos := c / fact
		c = c % fact

		result.Corners[i].Index = available[pos]
		available = append(available[:pos], available[pos+1:]...)
	}
	return result
}

// toUDCoordinate converts a Cube into the coordinate space of the perumation
// of the top and bottom layer edges. Importantly ignores the permutation of
// E slice edges and therefore has a problem size of 8!, just like CP.
// This problem space is undefined in Phase 1 because the 8 UD edges might not
// all be in the UD layers yet.
func toUDCoordinate(c cube) int {
	result := 0
	for i := 0; i < edgeFL; i++ {
		count := 0
		for j := i + 1; j < edgeFL; j++ {
			if c.Edges[j].Index < c.Edges[i].Index {
				count++
			}
		}
		result += count * factorial(edgeFL-i-1)
	}
	return result
}

// fromUDCoordinate converts a coordinate in the UD space into a cube.
// This cube is only valid for generating pruning tables for UD
func fromUDCoordinate(c int) cube {
	var result cube

	available := make([]int, edgeFL)
	for i := 0; i < edgeFL; i++ {
		available[i] = i
	}

	for i := 0; i < edgeFL; i++ {
		fact := factorial(edgeFL - i - 1)

		pos := c / fact
		c = c % fact

		result.Edges[i].Index = available[pos]
		available = append(available[:pos], available[pos+1:]...)
	}
	return result
}

// toESliceP2Coordinate converts a Cube into a coordinate in the
// E slice permutation space. There are only 4 E slice edges, so
// there are only 4! permutations of this space! This problem space
// is undefined in Phase 1 because all 4 edges may not yet be in the
// E slice.
func toESliceP2Coordinate(c cube) int {
	result := 0
	for i := 0; i < 4; i++ {
		count := 0
		for j := i + 1; j < 4; j++ {
			if c.Edges[j+edgeFL].Index < c.Edges[i+edgeFL].Index {
				count++
			}
		}
		result += count * factorial(4-i-1)
	}
	return result
}

// fromESliceP2Coordinate converts a coordinate into a Cube. This
// coordinate is only valid for generating pruning tables for the
// phase 2 E slice problem space.
func fromESliceP2Coordinate(c int) cube {
	var result cube

	available := make([]int, 4)
	for i := 0; i < 4; i++ {
		available[i] = i
	}

	for i := 0; i < 4; i++ {
		fact := factorial(4 - i - 1)

		pos := c / fact
		c = c % fact

		result.Edges[i+edgeFL].Index = available[pos] + edgeFL
		available = append(available[:pos], available[pos+1:]...)
	}
	return result
}

func factorial(n int) int {
	if n == 0 {
		return 1
	}
	return n * factorial(n-1)
}
