package base

type CellType int

const (
	Empty CellType = iota

	Occupied1
	Occupied2

	Captured1
	Captured2

	SemiCaptured1
	SemiCaptured2
)

var (
	Occupied     = [2]CellType{Occupied1, Occupied2}
	Captured     = [2]CellType{Captured1, Captured2}
	SemiCaptured = [2]CellType{SemiCaptured1, SemiCaptured2}
)

type CellPosition struct {
	X, Y int
}

type Edge struct {
	Start, End CellPosition
}
