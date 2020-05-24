package puzzle

type Combo []Chain

type Chain struct {
	Coordinates []ComboCoordinate
	Phase       int
}

type ComboCoordinate struct {
	Coordinate
	Initial Coordinate
}

type Coordinate struct {
	X int
	Y int
}
