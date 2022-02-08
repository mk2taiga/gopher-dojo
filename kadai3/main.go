package main

import (
	"fmt"
	game2 "gopher-dojo/kadai3/game"
	"os"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "Error:\n%s\n", err)
			os.Exit(1)
		}
	}()

	game := game2.Game{}
	os.Exit(game.Run())
}
