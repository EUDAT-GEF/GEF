package pier

import (
	"bytes"
	"sync"
)

// Task exported
type Task struct {
	JobID	JobID
	Items   []TaskStatus
}

// TaskStatus exported
type TaskStatus struct {
	Name          string
	Error         error
	ExitCode      int
	ConsoleOutput *bytes.Buffer
}

// TaskList is a shared structure that stores info about all containers related to jobs
type TaskList struct {
	sync.Mutex
	cache map[JobID]Task
}

// NewTaskList exported
func NewTaskList() *TaskList {
	return &TaskList{
		cache: make(map[JobID]Task),
	}
}

func (taskList *TaskList) addTask(jobID JobID, taskName string, taskError error, taskExitCode int, taskConsoleOutput *bytes.Buffer) {
	taskList.Lock()
	defer taskList.Unlock()
	task := taskList.cache[jobID]

	var newTask TaskStatus
	newTask.Name = taskName
	newTask.Error = taskError
	newTask.ExitCode = taskExitCode
	newTask.ConsoleOutput = taskConsoleOutput
	task.Items = append(task.Items, newTask)

	taskList.cache[jobID] = task
}