package main

import "strconv"

func create_task(node string, task string, arg string) Task {

	// Wouldn't want a small task ID
	taskID = taskID + 1

	// Not sure why i need to convert here, but 🤷‍♀️
	ID_string := strconv.Itoa(taskID)

	newCommand := Task{ID_string, node, task, arg, "waiting", ""}

	return newCommand

}

// Finds a nodes position in the slice
func find_task(Node_input_ID string) (int, string) {

	for i, value := range task_queue {

		if string(value.NodeID) == Node_input_ID {
			return i, ""
		}
	}

	// returns 0 if it can't find anything.
	// pretty sure this small brain, but ehh
	return 0, "💀 Couldn't find node"
}

// this is for finding a task that we want to get executed
func find_task_unsent(Node_input_ID string) (int, string) {

	for i, value := range task_queue {

		if string(value.NodeID) == Node_input_ID {

			if value.Progress == "waiting" {
				return i, ""
			}

		}
	}

	// returns 0 if it can't find anything.
	// pretty sure this small brain, but ehh
	return 0, "💀 Couldn't find node"
}

// this is for finding a task that we want to get executed
func find_task_by_id(input_task_id string) (int, string) {

	for i, value := range task_queue {

		if string(value.TaskID) == input_task_id {
			return i, ""
		}

	}

	// returns 0 if it can't find anything.
	// pretty sure this small brain, but ehh
	return 0, "💀 Couldn't find node"
}
