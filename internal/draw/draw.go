package draw

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/regzon/dots-bot/internal/base"
)

func DrawBoard(board *base.Board) {
	white := color.New(color.FgWhite)
	red := color.New(color.FgRed)
	cyan := color.New(color.FgCyan)

	for _, row := range board.Cells {
		for _, v := range row {

			switch v {
			case base.Empty:
				white.Print(".")
			case base.Occupied1:
				cyan.Print("*")
			case base.Occupied2:
				red.Print("*")
			case base.Captured1:
				cyan.Print("+")
			case base.Captured2:
				red.Print("+")
			}
		}
		fmt.Println()
	}
}
