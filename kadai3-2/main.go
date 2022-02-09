package main

import (
	"fmt"
	"gopher-dojo/kadai3-2/rangedownloader"
	"os"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "Error:\n%s\n", err)
			os.Exit(1)
		}
	}()
	cli := rangedownloader.New()
	os.Exit(cli.Run())
}
