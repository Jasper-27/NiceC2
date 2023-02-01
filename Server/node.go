package main

import "fmt"

////////////////////////////////////
/// Node slice  			    ////
////////////////////////////////////

func createNode(ID string, Hostname string, Platform string, Timestamp string) node {

	newNode := node{ID, Hostname, Platform, Timestamp, Timestamp}
	return newNode

}

func is_new_node(input_ID string) bool {

	for _, value := range nodes {
		if value.ID == input_ID {
			return false
		}
	}

	return true
}

// updates the node in the list
func update_node(ID string, Timestamp string) {

	node_position, error := find_node(ID)
	if error != "" {
		fmt.Println(error)
		return
	}

	nodes[node_position].Last_Check_In = Timestamp

}

// Finds a nodes position in the slice
func find_node(input_ID string) (int, string) {

	for i, value := range nodes {
		if value.ID == input_ID {
			return i, ""
		}
	}

	// returns 0 if it can't find anything.
	// pretty sure this small brain, but ehh
	return 0, "💀 Couldn't find node"
}

func display_all_nodes() {

	fmt.Println("")

	fmt.Println("[][][][][][][][][] Nodes [][][][][][][][][]")
	fmt.Println()
	for _, value := range nodes {
		fmt.Println("NodeID:        " + value.ID)
		fmt.Println("Hostname:      " + value.Hostname)
		fmt.Println("Platform:      " + value.Platform)
		fmt.Println("First Seen:    " + value.First_Check_In)
		fmt.Println("Last Seen:     " + value.Last_Check_In)
		fmt.Println("------------------ ===== ------------------")
	}
}
