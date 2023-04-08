package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func save_nodes_to_file() {
	out, _ := json.MarshalIndent(nodes, "", "  ")
	err := ioutil.WriteFile("nodes.json", out, 0644)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func read_nodes_from_file() []node {
	fmt.Println("These are the nodes allready in the file")
	fmt.Printf("\n\n")
	content, err := ioutil.ReadFile("nodes.json")
	if err != nil {
		fmt.Println(err.Error())
	}

	var nodes []node
	err3 := json.Unmarshal(content, &nodes)
	if err3 != nil {
		fmt.Println("error with Unmarshal")
		fmt.Println(err3.Error())
	}

	for _, x := range nodes {
		fmt.Println(x.ID)
	}

	fmt.Printf("\n\n")
	return nodes
}
