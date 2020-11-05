package game

import (
	"fmt"

	"github.com/regzon/dots-bot/internal/base"
	"github.com/regzon/dots-bot/internal/draw"
	"github.com/regzon/dots-bot/internal/players"
)

func MainLoop(width, height, maxDepth int) {
	board := base.NewBoard(width, height)
	board.AddInitCells()

	var playerArr [2]players.Player
	playerArr[0] = players.NewMinimaxBot(maxDepth, 0)
	playerArr[1] = players.NewRealPlayer()

	draw.DrawBoard(board)

	for i := 0; board.HasEmpty(); i++ {
		turn := i + 1
		playerInd := (turn - 1) % 2

		fmt.Printf("Turn %d\n", turn)
		fmt.Printf("Player index %d\n", playerInd)

		pos := playerArr[playerInd].ChooseCell(board)
		board.Occupy(pos, playerInd)

		draw.DrawBoard(board)

		fmt.Printf("Player1 score %d\n", board.GetScore(0))
		fmt.Printf("Player2 score %d\n", board.GetScore(1))
		fmt.Println()
	}

	fmt.Println("Finished")
}
