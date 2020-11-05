package base

import (
	"errors"
	"sort"

	"github.com/regzon/dots-bot/internal/algs"
)

type Board struct {
	Width, Height int
	Cells         [][]CellType
}

func NewBoard(width, height int) *Board {
	// Create a 2-d array filled with empty board cells
	cells := make([][]CellType, height)
	buf := make([]CellType, height*width)
	for i := range cells {
		cells[i] = buf[i*width : (i+1)*width]
	}

	return &Board{width, height, cells}
}

func (b *Board) AddInitCells() {
	x := b.Width/2 - 1
	y := b.Height/2 - 1

	b.Cells[y][x] = Occupied[0]
	b.Cells[y+1][x+1] = Occupied[0]

	b.Cells[y][x+1] = Occupied[1]
	b.Cells[y+1][x] = Occupied[1]
}

func (b *Board) Copy() (bCopy *Board) {
	bCopy = &Board{}

	bCopy.Width = b.Width
	bCopy.Height = b.Height

	bCopy.Cells = make([][]CellType, len(b.Cells))
	buf := make([]CellType, len(b.Cells)*len(b.Cells[0]))
	for i, row := range b.Cells {
		bCopy.Cells[i] = buf[i*len(row) : (i+1)*len(row)]
		for j, cellType := range row {
			bCopy.Cells[i][j] = cellType
		}
	}

	return
}

func (b *Board) HasEmpty() bool {
	for _, row := range b.Cells {
		for _, cellType := range row {
			if cellType == Empty {
				return true
			}
		}
	}
	return false
}

func (b *Board) Occupy(pos CellPosition, playerInd int) (err error) {
	if !b.CanOccupy(pos) {
		err = errors.New("Tried to occupy a non-empty cell")
		return
	}
	b.Cells[pos.Y][pos.X] = Occupied[playerInd]
	b.update(playerInd)

	return
}

func (b *Board) CanOccupy(pos CellPosition) bool {
	return b.Cells[pos.Y][pos.X] == Empty
}

func (b *Board) GetScore(playerInd int) (score int) {
	for _, row := range b.Cells {
		for _, cellType := range row {
			if cellType == Captured[playerInd] {
				score++
			}
		}
	}
	return
}

type cellBaseNode struct {
	pos           CellPosition
	adjacentCells []int
}

func (c *cellBaseNode) GetAdjacent() []int {
	return c.adjacentCells
}

func (b *Board) findAdjacentNodes(baseNodes []*cellBaseNode, playerInd int) {
	buf := make([]int, 8*len(baseNodes))
	for i, node := range baseNodes {
		adjacent := buf[i*8 : (i+1)*8]
		aLen := 0

		// TODO: rewrite using more efficient algorithm (not O(n*n))
		for _, pos := range b.GetCellNeighbors(node.pos, Occupied[playerInd]) {
			for i, adjNode := range baseNodes {
				if adjNode.pos == pos {
					adjacent[aLen] = i
					aLen++
				}
			}
		}

		node.adjacentCells = adjacent[:aLen]
	}
}

func (b *Board) createBaseNodes(playerInd int) []*cellBaseNode {
	baseNodes := make([]*cellBaseNode, b.Width*b.Height)
	bnLen := 0
	for y, row := range b.Cells {
		for x, cellType := range row {
			if cellType != Occupied[playerInd] {
				continue
			}

			pos := CellPosition{x, y}
			baseNodes[bnLen] = &cellBaseNode{pos: pos}
			bnLen++
		}
	}
	return baseNodes[:bnLen]
}

func (b *Board) getCyclesEdges(baseNodes []*cellBaseNode) [][]Edge {
	nodes := make([]algs.BaseNode, len(baseNodes))
	for i, n := range baseNodes {
		nodes[i] = n
	}

	cycles := algs.DetectCycles(nodes)
	cyclesEdges := make([][]Edge, len(cycles))
	ceLen := 0

	for _, cycle := range cycles {
		if len(cycle) < 4 {
			// Ignore cycles with length 3 and less
			// (they can't have any cell inside)
			continue
		}

		cycleLen := len(cycle)
		edges := make([]Edge, cycleLen)
		for i := range cycle {
			iS := i
			iE := (i + 1) % cycleLen
			edges[i].Start = cycle[iS].(*cellBaseNode).pos
			edges[i].End = cycle[iE].(*cellBaseNode).pos
		}

		cyclesEdges[ceLen] = edges
		ceLen++
	}
	return cyclesEdges[:ceLen]
}

func (b *Board) detectCaptured(cyclesEdges [][]Edge, iT, rT CellType, playerInd int) (has bool) {
	has = false

	for _, edges := range cyclesEdges {
		for y, row := range b.Cells {
			for x, cellType := range row {
				if cellType != iT {
					continue
				}

				pos := CellPosition{X: x, Y: y}
				cross := 0

				for _, edge := range edges {
					var high, low CellPosition

					if edge.Start.Y > edge.End.Y {
						high = edge.End
						low = edge.Start
					} else {
						high = edge.Start
						low = edge.End
					}

					isValidY := low.Y == pos.Y && high.Y == (pos.Y-1)
					isValidX := low.X > pos.X

					if isValidX && isValidY {
						cross++
					}
				}

				if cross%2 == 1 {
					b.Cells[y][x] = rT
					has = true
				}
			}
		}
	}

	return
}

func (b *Board) update(playerInd int) {
	baseNodes := b.createBaseNodes(playerInd)
	b.findAdjacentNodes(baseNodes, playerInd)
	cyclesEdges := b.getCyclesEdges(baseNodes)

	// Sort is required to handle nested cycles correctly
	sort.Slice(cyclesEdges, func(i, j int) bool {
		return len(cyclesEdges[i]) < len(cyclesEdges[j])
	})

	has := b.detectCaptured(
		cyclesEdges,
		Occupied[(playerInd+1)%2],
		Captured[playerInd],
		playerInd,
	)
	if has {
		b.detectCaptured(
			cyclesEdges,
			Empty,
			SemiCaptured[playerInd],
			playerInd,
		)
	}
}

type posDiff struct{ x, y int }

var neighborsDiff = [...]posDiff{
	{0, 1}, {1, 1}, {-1, 1},
	{0, -1}, {1, -1}, {-1, -1},
	{1, 0}, {-1, 0},
}

func (b *Board) GetCellNeighbors(pos CellPosition, t CellType) []CellPosition {
	neigh := make([]CellPosition, 8)
	nLen := 0

	for _, diff := range neighborsDiff {
		newPos := CellPosition{pos.X + diff.x, pos.Y + diff.y}

		isValidX := newPos.X >= 0 && newPos.X < b.Width
		isValidY := newPos.Y >= 0 && newPos.Y < b.Height

		if isValidX && isValidY {
			if b.GetCellType(newPos) == t {
				neigh[nLen] = newPos
				nLen++
			}
		}
	}

	return neigh[:nLen]
}

func (b *Board) areNeighbors(f, s CellPosition) bool {
	diffY := f.Y - s.Y
	diffX := f.X - s.X
	return diffY >= -1 && diffY <= 1 && diffX >= -1 && diffX <= 1
}

func (b *Board) GetCellType(pos CellPosition) CellType {
	return b.Cells[pos.Y][pos.X]
}
