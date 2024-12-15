package cube

type cube = struct {
	Edges   [edgeCount]piece
	Corners [cornerCount]piece
}

type piece = struct {
	Index       int
	Orientation int
}

// transform applies a permutation and orientation to a cube.
// Defining the transform we'd like to apply is actually the same
// shape as defining the current state of a cube. You can think
// of the current state of the cube a as a tranform of the solved cube
func transform(a, b cube) cube {
	var result cube
	for i, tEdge := range b.Edges {
		cEdge := a.Edges[tEdge.Index]
		result.Edges[i] = piece{
			Index:       cEdge.Index,
			Orientation: (cEdge.Orientation + tEdge.Orientation) % 2,
		}
	}
	for i, tCorner := range b.Corners {
		cCorner := a.Corners[tCorner.Index]
		result.Corners[i] = piece{
			Index:       cCorner.Index,
			Orientation: (cCorner.Orientation + tCorner.Orientation) % 3,
		}
	}
	return result
}
