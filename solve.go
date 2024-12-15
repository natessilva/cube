package cube

type solver struct {
	scrambledCube cube
	path          []int
	pathPhase2    []int
}

func newSolver(scramble string) (*solver, error) {
	c, err := parseScrambe(scramble)
	if err != nil {
		return nil, err
	}
	return &solver{
		path:          make([]int, 0),
		pathPhase2:    make([]int, 0),
		scrambledCube: c,
	}, nil
}

func Solve(scrambe string) (string, error) {
	solver, err := newSolver(scrambe)
	if err != nil {
		return "", err
	}
	return solver.solve(), nil
}

func (s *solver) solve() string {
	// convert the cube into the 3 phase 1 coordinates:
	eoCoord := toEOCoordinate(s.scrambledCube)
	coCoord := toCOCoordinate(s.scrambledCube)
	ePermCoord := toESliceP1Coordinate(s.scrambledCube)

	cost := phase1Hueristic(eoCoord, coCoord, ePermCoord)

	// because our hueristic find the lower bound for possible
	// solves, it might actually under estimate. So if we don't
	// actually find any solves in our cost, we will just
	// keep bumping the estimated cost until we do find a solution.
	// we also do not stop with the first phase 1 solution we find.
	// we continue to find sub-optimal solutions (up to 20 moves) for
	// phase1 in hopes that we find one that "sets up" a good phase 2
	// and provides a generally efficient solve.
	for i := cost; i < 20; i++ {
		done := s.search(eoCoord, coCoord, ePermCoord, i)
		if done {
			break
		}
	}

	return toString(append(s.path, s.pathPhase2...))
}

func (s *solver) search(eoCoord, coCoord, ePermCoord, phase1Cost int) bool {
	// if the cost is zero than we have in fact found a solution to phase1 🎉
	// begin searching for phase2 solutions
	if phase1Cost == 0 {
		// Don't ever end a phase 1 solution with phase 2 moves (half turns of LRFB or any UD turns)
		// because these turns do not break phase1 (meaning we already solved it before we got here)
		if len(s.path) > 0 {
			m := s.path[len(s.path)-1]
			if m < moveL || m > moveB3 || m%3 == 1 {
				return false
			}
		}
		// since we haven't actually been performing moves on a real Cube, we don't know the current
		// state. Compute it by applying all phase 1 moves to the original scramble.
		newCube := s.scrambledCube
		for i := 0; i < len(s.path); i++ {
			newCube = transform(newCube, moves[s.path[i]])
		}
		// convert to the phase 2 coordinate system
		cpCoord := toCPCoordinate(newCube)
		eudCoord := toUDCoordinate(newCube)
		eePermCoord := toESliceP2Coordinate(newCube)

		// limit phase 2 solutions to size 10
		// The reason to do this is that the search space for Phase 2
		// is very large and searching beyong 10 deep into the tree becomes
		// quite time consuming. Also, if we have a more than 10 move Phase 2
		// there is almost certainly a less optimal phase 1 that leads to
		// a much shorter phase 2
		phase2Limit := 11

		newCost := phase2Hueristic(cpCoord, eudCoord, eePermCoord)
		// same as in phase 1 - our huerist might under estimate and so
		// try again with bigger numbers if we don't find anything
		for i := newCost; i < phase2Limit; i++ {
			if s.searchPhase2(cpCoord, eudCoord, eePermCoord, i) {
				// finally, if we have found a phase 2 solutions we are done
				// and we return true back through the call stack
				return true
			}
		}
		return false
	} else {
		// otherwise, apply every possible move and see if it is worth searching
		for m := 0; m < moveCount; m++ {
			// don't perform sequential moves of the same face
			if len(s.path) > 0 {
				lastMoveFace := s.path[len(s.path)-1] / 3
				mFace := m / 3
				if mFace == lastMoveFace {
					continue
				}
			}

			// use our lookup tables to get new coordinates and cost quickly
			eoNew := eoLookup[eoCoord][m]
			coNew := coLookup[coCoord][m]
			ePermNew := eSliceP1Lookup[ePermCoord][m]
			costNew := phase1Hueristic(eoNew, coNew, ePermNew)

			// don't explore paths that are too expensive
			if costNew >= phase1Cost {
				continue
			}

			s.path = append(s.path, m)
			done := s.search(eoNew, coNew, ePermNew, phase1Cost-1)
			if done {
				return true
			}
			s.path = s.path[0 : len(s.path)-1]
		}
	}
	return false
}

func (s *solver) searchPhase2(cpCoord, eudCoord, eePermCoord, phase2Cost int) bool {
	if phase2Cost == 0 {
		return true
	} else {
		for m := 0; m < moveCount; m++ {
			// don't do quarter turns on the side faces as these are not valid moves in phase 2
			if m%3 != 1 && m > moveU3 && m < moveD {
				continue
			}

			// don't perform sequential moves of the same face in phase 2
			if len(s.pathPhase2) > 0 {
				lastMoveFace := s.pathPhase2[len(s.pathPhase2)-1] / 3
				mFace := m / 3
				if mFace == lastMoveFace {
					continue
				}
				// also don't start a phase 2 solution with the same face as our phase 1 solution ended with
			} else if len(s.path) > 0 {
				lastMoveFace := s.path[len(s.path)-1] / 3
				mFace := m / 3
				if mFace == lastMoveFace {
					continue
				}
			}

			// use our lookup tables to get the new coordinates and cost
			cpNew := cpLookup[cpCoord][m]
			eudNew := udLookup[eudCoord][m]
			eePermNew := eSliceP2Lookup[eePermCoord][m]
			costNew := phase2Hueristic(cpNew, eudNew, eePermNew)

			// don't explore branches that are too expensive
			if costNew >= phase2Cost {
				continue
			}
			s.pathPhase2 = append(s.pathPhase2, m)
			done := s.searchPhase2(cpNew, eudNew, eePermNew, phase2Cost-1)
			if done {
				return true
			}
			s.pathPhase2 = s.pathPhase2[0 : len(s.pathPhase2)-1]
		}
	}
	return false
}
