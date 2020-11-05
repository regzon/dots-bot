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

	fmt.Print(" ")
	for i := range board.Cells[0] {
		fmt.Printf("%c", 'A'+i)
	}
	fmt.Println()

	for i, row := range board.Cells {
		fmt.Printf("%c", 'A'+i)

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
			case base.SemiCaptured1:
				cyan.Print("-")
			case base.SemiCaptured2:
				red.Print("-")
			}
		}
		fmt.Println()
	}
}
