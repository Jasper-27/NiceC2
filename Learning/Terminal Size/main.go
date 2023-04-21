package main

import (
	"fmt"

	"golang.org/x/term"
)

func main() {

	fmt.Println(find_terminal_size())
}

func find_terminal_size() int {
	if !term.IsTerminal(0) {
		return 0
	}

	width, _, err := term.GetSize(0)
	if err != nil {
		return 0
	}

	return width
}
