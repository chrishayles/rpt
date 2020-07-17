package rpt

import (
	"context"
	"fmt"
	"time"
)

// OPERATIONS

type DBOperation struct {
	errors    []error
	Name      string
	created   time.Time
	started   time.Time
	completed time.Time
	operation func(db DBClient, data DataSet) (s string, e error)
	client    DBClient
	result    string
	data      DataSet
}

func (dbo *DBOperation) Start() {
	dbo.started = time.Now()

	res, err := dbo.operation(dbo.client, dbo.data)

	dbo.result = res
	if err != nil {
		dbo.errors = append(dbo.errors, err)
	}

	d := dbo.Complete()
	fmt.Printf("\nCompleted in %s\n", d)
}

func (dbo *DBOperation) Started() time.Time {
	return dbo.started
}

func (dbo *DBOperation) Complete() time.Duration {
	dbo.completed = time.Now()
	return dbo.Duration()
}

func (dbo *DBOperation) Completed() time.Time {
	return dbo.completed
}

func (dbo *DBOperation) Duration() time.Duration {
	return dbo.completed.Sub(dbo.started)
}

func (dbo *DBOperation) GetData() DataSet {
	return dbo.data
}

func (dbo *DBOperation) GetResult() string {
	return dbo.result
}

type DBOperationSet struct {
	Operations []*DBOperation
	ctx        context.Context
}

func (dbos *DBOperationSet) Cancel() {}

// CLIENTS

type DBClient interface {
	Connect() error
	Disconnect() error
	Seed(d DataSet) error
}

// OPERATION FUNCTIONS

func newDBOperator(n string, c DBClient, d DataSet, o func(db DBClient, data DataSet) (string, error)) *DBOperation {

	e := []error{}

	return &DBOperation{
		result:    "",
		errors:    e,
		created:   time.Now(),
		started:   time.Time{},
		completed: time.Time{},
		Name:      n,
		client:    c,
		operation: o,
		data:      d,
	}
}

func SeedData(client DBClient, data DataSet) *DBOperation {

	dbo := newDBOperator("seed_data", client, data, func(db DBClient, data DataSet) (string, error) {

		err := db.Seed(data)
		if err != nil {
			return "", err
		}

		return "Created.", nil
	})

	return dbo
}
