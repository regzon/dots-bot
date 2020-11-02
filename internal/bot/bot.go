package bot

import "github.com/regzon/dots-bot/internal/base"

type Bot interface {
	ChooseCell(board *base.Board) base.CellPosition
}

func NewBot() Bot {
	return &RandomBot{}
}
