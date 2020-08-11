package rpt

import (
	"fmt"
	"time"
)

// STRUCTS
type Workflow struct {
	operations DBOperationSet
	Name       string
	started    time.Time
	completed  time.Time
}

// FUNCTIONS
func (w *Workflow) Start() {
	w.started = time.Now()

	// res, err := w.operation(w.client, w.data)

	// w.result = res
	// if err != nil {
	// 	w.errors = append(w.errors, err)
	// }

	d := w.Complete()
	fmt.Printf("\nCompleted in %s\n", d)
}

func (w *Workflow) Started() time.Time {
	return w.started
}

func (w *Workflow) Complete() time.Duration {
	w.completed = time.Now()
	return w.Duration()
}

func (w *Workflow) Completed() time.Time {
	return w.completed
}

func (w *Workflow) Duration() time.Duration {
	return w.completed.Sub(w.started)
}

// IMPLEMENTATIONS

func NewWorkflow(w string, r RptClient) *Workflow {

	//Do stuff

	return &Workflow{}
}

func startupWorkflow(p DBClient, s DBClient) *Workflow {
	return &Workflow{}
}

func replicationWorkflow(p DBClient, s DBClient) *Workflow {
	return &Workflow{}
}

func reconfigureClientWorkflow(p DBClient, s DBClient) *Workflow {
	return &Workflow{}
}
