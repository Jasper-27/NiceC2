package main

type Task struct {
	TaskID   string
	NodeID   string
	Action   string
	Content  string
	Progress string
}

type Task_Response struct {
	TaskID   string
	Progress string // Completed / Failed
	Result   string // Data from the task
}

type check_in struct {
	ID       string
	Hostname string
	Platform string
}

type node struct {
	ID             string
	Hostname       string
	Platform       string
	First_Check_In string
	Last_Check_In  string
}
