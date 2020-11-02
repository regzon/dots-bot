package game

import (
	"fmt"

	"github.com/regzon/dots-bot/internal/base"
	"github.com/regzon/dots-bot/internal/bot"
	"github.com/regzon/dots-bot/internal/draw"
)

func MainLoop(width, height int) {
	board := base.NewBoard(width, height)
	bot := bot.NewBot()

	draw.DrawBoard(board)

	for i := 1; i < 500; i++ {
		// time.Sleep(time.Second)
		fmt.Printf("Turn %d\n", i)

		pos := bot.ChooseCell(board)
		board.Occupy(pos, i%2)

		draw.DrawBoard(board)
		fmt.Println()
	}
}
