package main

// task sent to node to complete. and task stored in server
type Task struct {
	TaskID   string
	NodeID   string
	Action   string
	Content  string
	Progress string
	Result   string
}

// Send to server when node completes response
type Task_Response struct {
	TaskID   string
	Progress string // Completed / Failed
	Result   string // Data from the task
}

// Data sent to server each time node checks in
type check_in struct {
	ID       string
	Hostname string
	Platform string
}

// Data stored about each node
type node struct {
	ID             string
	Hostname       string
	Platform       string
	First_Check_In string
	Last_Check_In  string
}
