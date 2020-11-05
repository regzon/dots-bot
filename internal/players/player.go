package players

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/regzon/dots-bot/internal/base"
)

type Player interface {
	ChooseCell(board *base.Board) base.CellPosition
}

type RealPlayer struct{}

func readPos(board *base.Board) (pos base.CellPosition) {
	b := make([]byte, 1)

	tried := false
	failed := false
	fmt.Print("Enter x coordinate: ")
	for !tried || failed {
		tried = true
		failed = false

		os.Stdin.Read(b)
		fmt.Println()

		var x int = int(b[0]) - 65

		if x < 0 || x >= board.Width {
			failed = true
			fmt.Print("Invalid character.\nTry again: ")
			continue
		}

		pos.X = x
	}

	tried = false
	failed = false
	fmt.Print("Enter y coordinate: ")
	for !tried || failed {
		tried = true
		failed = false

		os.Stdin.Read(b)
		fmt.Println()

		var y int = int(b[0]) - 65

		if y < 0 || y >= board.Height {
			failed = true
			fmt.Print("Invalid character.\nTry again: ")
			continue
		}

		pos.Y = y
	}

	return
}

func (player *RealPlayer) ChooseCell(board *base.Board) (pos base.CellPosition) {
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()

	tried := false
	failed := false

	for !tried || failed {
		tried = true
		failed = false

		pos = readPos(board)

		if board.GetCellType(pos) != base.Empty {
			fmt.Println("Cell is busy, try another")
			failed = true
			continue
		}
	}

	return
}

func NewRealPlayer() (player Player) {
	player = &RealPlayer{}
	return
}
