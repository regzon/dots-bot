package base

import (
	"log"
	"sort"

	"github.com/regzon/dots-bot/internal/algs"
)

type CellType int

const (
	Empty CellType = iota

	Occupied1
	Occupied2

	Captured1
	Captured2
)

var (
	Occupied = [2]CellType{Occupied1, Occupied2}
	Captured = [2]CellType{Captured1, Captured2}
)

type CellPosition struct {
	X, Y int
}

type Edge struct {
	Start, End CellPosition
}

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

func (b *Board) Occupy(pos CellPosition, player int) {
	if !b.CanOccupy(pos) {
		log.Fatal("Tried to occupy a non-empty cell")
	}
	b.Cells[pos.Y][pos.X] = Occupied[player]
	b.update(player)
}

func (b *Board) CanOccupy(pos CellPosition) bool {
	return b.Cells[pos.Y][pos.X] == Empty
}

type cellBaseNode struct {
	pos           CellPosition
	adjacentCells []int
}

func (c *cellBaseNode) GetAdjacent() []int {
	return c.adjacentCells
}

func (b *Board) update(player int) {
	baseNodes := []*cellBaseNode{}

	for y, row := range b.Cells {
		for x, cellType := range row {
			if cellType != Occupied[player] {
				continue
			}

			pos := CellPosition{x, y}
			baseNodes = append(baseNodes, &cellBaseNode{pos: pos})
		}
	}

	// Find adjacent nodes
	for _, node := range baseNodes {
		adjacent := []int{}

		// TODO: rewrite using more efficient algorithm (not O(n*n))
		for _, pos := range b.getCellNeighbors(node.pos, Occupied[player]) {
			for i, adjNode := range baseNodes {
				if adjNode.pos == pos {
					adjacent = append(adjacent, i)
				}
			}
		}

		node.adjacentCells = adjacent
	}

	nodes := make([]algs.BaseNode, len(baseNodes))
	for i, n := range baseNodes {
		nodes[i] = n
	}

	cycles := algs.DetectCycles(nodes)
	cyclesEdges := [][]Edge{}

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

		cyclesEdges = append(cyclesEdges, edges)
	}

	// Sort is required to handle nested cycles correctly
	sort.Slice(cyclesEdges, func(i, j int) bool {
		return len(cyclesEdges[i]) < len(cyclesEdges[j])
	})

	for _, edges := range cyclesEdges {
		for y, row := range b.Cells {
			for x, _ := range row {
				pos := CellPosition{X: x, Y: y}
				contains := false
				cross := 0

				for _, edge := range edges {
					if edge.Start == pos {
						contains = true
						break
					}

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

				if contains {
					continue
				}

				if cross%2 == 1 {
					// TODO: fix the problem of extra captured cells
					// when cycles are not induced
					b.Cells[y][x] = Captured[player]
				}
			}
		}
	}
}

type posDiff struct{ x, y int }

var neighborsDiff = [...]posDiff{
	{0, 1}, {1, 1}, {-1, 1},
	{0, -1}, {1, -1}, {-1, -1},
	{1, 0}, {-1, 0},
}

func (b *Board) getCellNeighbors(pos CellPosition, t CellType) (res []CellPosition) {
	for _, diff := range neighborsDiff {
		newPos := CellPosition{pos.X + diff.x, pos.Y + diff.y}

		isValidX := newPos.X >= 0 && newPos.X < b.Width
		isValidY := newPos.Y >= 0 && newPos.Y < b.Height

		if isValidX && isValidY {
			if b.getCellType(newPos) == t {
				res = append(res, newPos)
			}
		}
	}

	return
}

func (b *Board) getCellType(pos CellPosition) CellType {
	return b.Cells[pos.Y][pos.X]
}
