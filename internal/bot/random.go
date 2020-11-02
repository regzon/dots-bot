package bot

import (
	"math/rand"

	"github.com/regzon/dots-bot/internal/base"
)

type RandomBot struct{}

func (bot *RandomBot) ChooseCell(board *base.Board) base.CellPosition {
	var cells []base.CellPosition
	for y, row := range board.Cells {
		for x, cellType := range row {
			if cellType == base.Empty {
				cells = append(cells, base.CellPosition{X: x, Y: y})
			}
		}
	}

	randInd := rand.Intn(len(cells))
	return cells[randInd]
}
