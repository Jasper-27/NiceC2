package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/gosuri/uitable"
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

	// t.AddRow("Command", "Description")
	// t.AddRow("ls", "List all nodes")
	// t.AddRow("use [node]", "Set node you are working on")
	// t.AddRow("tasks [node]", "view the task queue for that node. Leaving blank will print all tasks")
	// t.AddRow("run [node]", "run a single command on a node")
	// t.AddRow("shutdown [node]", "ask a node to shutdown")
	// t.AddRow("reboot [node]", "ask a node to reboot")
	// t.AddRow("send-file [node] -f [filename] -d [destination file path]", "Send a file from the server to the node")
	// t.AddRow("get-file [node] -p [file path on node]", "Get a file from a node, and store it on the server")
	// t.AddRow("payloads", "List all the payloads available in the payloads folder")
	// t.AddRow("Exit", "Exit the NiceC2 command line")

	// fmt.Println(t)

	table := uitable.New()
	table.MaxColWidth = 50
	table.Wrap = true

	for _, help_message := range help_messages {
		table.AddRow(color.RedString("Command"), help_message.command)
		table.AddRow(color.BlueString("Description"), "  "+help_message.description)
		// table.AddRow("") // blank
	}
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
