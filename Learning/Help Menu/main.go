package main

import (
	"fmt"

	"github.com/fatih/color"

	"github.com/cheynewallace/tabby"
)

func main() {
	fmt.Println("Hello")

	color.Blue("Hello there")

	color.Red("We have red")
	color.Magenta("And many others ..")

	// Create a new color object
	c := color.New(color.FgCyan).Add(color.Underline)
	c.Println("Prints cyan text with an underline.")

	// Or just add them to New()
	d := color.New(color.FgCyan, color.Bold)
	d.Printf("This prints bold cyan %s\n", "too!.")

	// Mix up foreground and background colors, create new mixes!
	red := color.New(color.FgRed)

	boldRed := red.Add(color.Bold)
	boldRed.Println("This will print text in bold red.")

	fmt.Println()

	t := tabby.New()
	t.AddHeader("NAME", "TITLE", "DEPARTMENT")
	t.AddLine("John Smith", "Developer", "Engineering")
	t.Print()

}
