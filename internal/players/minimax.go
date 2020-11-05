package players

import (
	"fmt"

	"github.com/regzon/dots-bot/internal/base"
)

type MinimaxBot struct {
	maxDepth  int
	playerInd int
}

type cellResult struct {
	score int
	pos   base.CellPosition
}

func getScore(board *base.Board, playerInd int, isMax bool) int {
	if isMax {
		return board.GetScore(playerInd) - board.GetScore((playerInd+1)%2)
	}
	return board.GetScore((playerInd+1)%2) - board.GetScore(playerInd)
}

func cellWorker(
	initBoard *base.Board,
	playerInd int,
	isMax bool,
	curDepth int,
	maxDepth int,
	pos base.CellPosition,
	resultScore chan<- cellResult,
) {
	var score int

	board := initBoard.Copy()

	err := board.Occupy(pos, playerInd)

	if curDepth == maxDepth || err != nil {
		score = getScore(board, playerInd, isMax)
	} else {
		score, _ = minimaxIter(
			board,
			(playerInd+1)%2,
			!isMax,
			curDepth+1,
			maxDepth,
		)
	}

	resultScore <- cellResult{score, pos}
}

func minimaxIter(
	initBoard *base.Board,
	playerInd int,
	isMax bool,
	curDepth int,
	maxDepth int,
) (bestScore int, choice base.CellPosition) {

	hasBest := false
	choice = base.CellPosition{X: -1, Y: -1}

	cnt := 0
	ch := make(chan cellResult, initBoard.Width*initBoard.Height)

	hasEmp := false
	hasNei := false
	for y, row := range initBoard.Cells {
		for x, cellType := range row {
			if cellType != base.Empty {
				continue
			}

			hasEmp = true

			pos := base.CellPosition{X: x, Y: y}
			if len(initBoard.GetCellNeighbors(pos, base.Occupied[playerInd])) == 0 {
				continue
			}

			hasNei = true
			break
		}

		if hasNei && hasEmp {
			break
		}
	}

	if !hasEmp {
		bestScore = getScore(initBoard, playerInd, isMax)
		return
	}

	for y, row := range initBoard.Cells {
		for x, cellType := range row {
			if cellType != base.Empty {
				continue
			}

			pos := base.CellPosition{X: x, Y: y}
			if hasNei {
				if len(initBoard.GetCellNeighbors(pos, base.Occupied[playerInd])) == 0 {
					// Skip cells that don't have any occupied cell in neighbors
					continue
				}
			}

			cnt += 1

			if curDepth <= 2 {
				go cellWorker(initBoard, playerInd, isMax, curDepth, maxDepth, pos, ch)
			} else {
				cellWorker(initBoard, playerInd, isMax, curDepth, maxDepth, pos, ch)
			}
		}
	}

	for i := 0; i < cnt; i++ {
		res := <-ch

		isBest := (isMax && res.score > bestScore) || (!isMax && res.score < bestScore)
		if !hasBest || isBest {
			hasBest = true
			bestScore = res.score
			choice = res.pos
		}
	}

	return
}

func (bot *MinimaxBot) ChooseCell(board *base.Board) base.CellPosition {
	score, choice := minimaxIter(board, bot.playerInd, true, 1, bot.maxDepth)
	fmt.Printf("Best score is: %d\n", score)
	return choice
}

func NewMinimaxBot(depth int, playerInd int) (player Player) {
	player = &MinimaxBot{depth, playerInd}
	return
}
