package cube

const (
	moveU = iota
	moveU2
	moveU3
	moveL
	moveL2
	moveL3
	moveF
	moveF2
	moveF3
	moveR
	moveR2
	moveR3
	moveB
	moveB2
	moveB3
	moveD
	moveD2
	moveD3

	moveCount
)

const (
	edgeUB = iota
	edgeUR
	edgeUF
	edgeUL
	edgeDF
	edgeDR
	edgeDB
	edgeDL
	edgeFL
	edgeFR
	edgeBR
	edgeBL

	edgeCount
)

const (
	cornerUBL = iota
	cornerURB
	cornerUFR
	cornerULF
	cornerDFL
	cornerDRF
	cornerDBR
	cornerDLB

	cornerCount
)
