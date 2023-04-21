package main

import (
	"fmt"

	"github.com/gosuri/uitable"
	"golang.org/x/term"
)

type hacker struct {
	Name, Birthday, Bio string
}

var hackers = []hacker{
	{"Ada Lovelace", "December 10, 1815", "Ada was a British mathematician and writer, chiefly known for her work on Charles Babbage's early mechanical general-purpose computer, the Analytical Engine"},
	{"Alan Turing", "June 23, 1912", "Alan was a British pioneering computer scientist, mathematician, logician, cryptanalyst and theoretical biologist"},
}

type help_message struct {
	command, description string
}

var help_messages = []help_message{
	{"Command", "Description"},
	{"ls", "List all nodes"},
	{"use [node]", "Set node you are working on"},
	{"tasks [node]", "view the task queue for that node. Leaving blank will print all tasks"},
	{"run [node]", "run a single command on a node"},
	{"shutdown [node]", "ask a node to shutdown"},
	{"reboot [node]", "ask a node to reboot"},
	{"send-file [node] -f [filename] -d [destination file path]", "Send a file from the server to the node"},
	{"get-file [node] -p [file path on node]", "Get a file from a node, and store it on the server"},
	{"payloads", "List all the payloads available in the payloads folder"},
	{"Exit", "Exit the NiceC2 command line"},
}

func main() {

	// t := uitable.New()
	// t.MaxColWidth = 50
	// t.Wrap = true

	// fmt.Println(t)

	table := uitable.New()
	table.MaxColWidth = 20

	x := find_terminal_size()

	if x > 10 {
		table.MaxColWidth = uint(x)/2 - 5 // Devides x by 2, and rounds down if it's odd
	}
	table.Wrap = true

	table.AddRow("Command", "Description")
	table.AddRow("ls", "List all nodes")
	table.AddRow("use [node]", "Set node you are working on")
	table.AddRow("tasks [node]", "view the task queue for that node. Leaving blank will print all tasks")
	table.AddRow("run [node]", "run a single command on a node")
	table.AddRow("shutdown [node]", "ask a node to shutdown")
	table.AddRow("reboot [node]", "ask a node to reboot")
	table.AddRow("send-file [node] -f [filename] -d [destination file path]", "Send a file from the server to the node")
	table.AddRow("get-file [node] -p [file path on node]", "Get a file from a node, and store it on the server")
	table.AddRow("payloads", "List all the payloads available in the payloads folder")
	table.AddRow("Exit", "Exit the NiceC2 command line")

	// for _, help_message := range help_messages {
	// 	table.AddRow(color.RedString("Command"), help_message.command)
	// 	table.AddRow(color.BlueString("Description"), "  "+help_message.description)
	// 	// table.AddRow("") // blank
	// }
	fmt.Println(table)

	// fmt.Println("==> List")
	// table.AddRow(color.RedString("Name:"), "BIRTHDAY", "BIO")
	// for _, hacker := range hackers {
	// 	table.AddRow(hacker.Name, hacker.Birthday, hacker.Bio)
	// }
	// fmt.Println(table)

	// fmt.Print("\n==> Details\n")
	// table = uitable.New()
	// table.MaxColWidth = 80
	// table.Wrap = true
	// for _, hacker := range hackers {
	// 	table.AddRow("Name:", hacker.Name)
	// 	table.AddRow("Birthday:", hacker.Birthday)
	// 	table.AddRow("Bio:", hacker.Bio)
	// 	table.AddRow("") // blank
	// }
	// fmt.Println(table)

	// fmt.Print("\n==> Multicolor Support\n")
	// table = uitable.New()
	// table.MaxColWidth = 80
	// table.Wrap = true
	// for _, hacker := range hackers {
	// 	table.AddRow(color.RedString("Name:"), color.WhiteString(hacker.Name))
	// 	table.AddRow(color.BlueString("Birthday:"), hacker.Birthday)
	// 	table.AddRow(color.GreenString("Bio:"), hacker.Bio)
	// 	table.AddRow("") // blank
	// }
	// fmt.Println(table)
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
