package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/regzon/dots-bot/internal/game"
)

func main() {
	fmt.Println("You've started playing dots")

	if len(os.Args) != 4 {
		fmt.Printf("Usage: %s WIDTH HEIGHT MAXDEPTH\n", os.Args[0])
		os.Exit(1)
	}

	width, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Failed to parse width parameter")
		os.Exit(1)
	}

	height, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Failed to parse height parameter")
		os.Exit(1)
	}

	maxDepth, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("Failed to parse maxDepth parameter")
		os.Exit(1)
	}

	game.MainLoop(width, height, maxDepth)
}
