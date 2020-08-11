package rpt

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// OPERATIONS

type DBOperation struct {
	errors    []error
	Name      string
	ID        string
	created   time.Time
	started   time.Time
	completed time.Time
	operation func(db DBClient, data DataSet) (i interface{}, e error)
	client    DBClient
	result    interface{}
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

	if dbo.completed.Sub(dbo.started) <= 0 {
		return 0
	}

	return dbo.completed.Sub(dbo.started)
}

func (dbo *DBOperation) GetData() DataSet {
	return dbo.data
}

func (dbo *DBOperation) GetResult() []byte {

	output := &map[string]interface{}{
		"Result": dbo.result,
		"Errors": dbo.errors,
	}

	outputJson, _ := json.MarshalIndent(output, "", "  ")

	return outputJson
}

func (dbo *DBOperation) GetOutputJSON() []byte {

	res := make(map[string]interface{})
	_ = json.Unmarshal(dbo.GetResult(), &res)

	newObject := &map[string]interface{}{
		"Completed": dbo.Completed(),
		"Started":   dbo.Started(),
		"Duration":  dbo.Duration().String(),
		"Output":    res,
	}
	return ToJSON(newObject)
}

type DBOperationSet struct {
	Operations      []*DBOperation
	ctx             context.Context
	ID              string
	lookupOperation map[string]*DBOperation
}

func (dbos *DBOperationSet) GetOutputJSON() []byte {

	ops := make(map[string]interface{})

	for id, o := range dbos.lookupOperation {

		newObject := &map[string]interface{}{
			"Completed": nil,
			"Started":   nil,
			"Duration":  nil,
			"Output":    nil,
		}

		_ = json.Unmarshal(o.GetOutputJSON(), newObject)

		ops[id] = newObject
	}

	output := &map[string]interface{}{
		"ID":         dbos.ID,
		"Operations": ops,
	}

	return ToJSON(output)
}

func (dbos *DBOperationSet) Cancel() {}

func (dbos *DBOperationSet) AddOperation(dbo *DBOperation) {
	dbos.lookupOperation[dbo.ID] = dbo
	dbos.Operations = append(dbos.Operations, dbo)
}

func (dbos *DBOperationSet) LookupOperation(guid string) *DBOperation {

	if val, ok := dbos.lookupOperation[guid]; ok {
		return val
	}

	return nil
}

// CLIENTS

type DBClient interface {
	Connect() error
	Disconnect() error
	Reconnect() error
	Seed(d DataSet) (interface{}, error)
	Query(s string) (interface{}, error)
	ListDB() (interface{}, error)
}

// OPERATION FUNCTIONS

func newDBOperation(n string, c DBClient, d DataSet, o func(db DBClient, data DataSet) (interface{}, error)) *DBOperation {

	e := []error{}

	dbo := &DBOperation{
		result:    "",
		errors:    e,
		created:   time.Now(),
		started:   time.Time{},
		completed: time.Time{},
		Name:      n,
		client:    c,
		operation: o,
		data:      d,
		ID:        NewGUID(),
	}

	return dbo
}

func newDBOperationSet(ctx context.Context) *DBOperationSet {

	lookup := &map[string]*DBOperation{}

	return &DBOperationSet{
		ID:              NewGUID(),
		ctx:             ctx,
		Operations:      []*DBOperation{},
		lookupOperation: *lookup,
	}
}

func SeedData(client DBClient, data DataSet) *DBOperation {

	dbo := newDBOperation("seed_data", client, data, func(db DBClient, data DataSet) (interface{}, error) {

		res, err := db.Seed(data)
		if err != nil {
			return "", err
		}

		return res, nil
	})

	return dbo
}

func ReadData(client DBClient, query string) *DBOperation {

	data := &DBQueryDataSet{
		Name:  "query",
		Query: query,
	}

	dbo := newDBOperation("read_data", client, data, func(db DBClient, data DataSet) (interface{}, error) {

		res, err := db.Seed(data)
		if err != nil {
			return "", err
		}

		return res, nil
	})

	return dbo
}

func WriteData(client DBClient, data DataSet) *DBOperation {

	dbo := newDBOperation("write_data", client, data, func(db DBClient, data DataSet) (interface{}, error) {

		res, err := db.Seed(data)
		if err != nil {
			return "", err
		}

		return res, nil
	})

	return dbo
}

func DeleteData(client DBClient, data DataSet) *DBOperation {

	dbo := newDBOperation("delete_data", client, data, func(db DBClient, data DataSet) (interface{}, error) {

		res, err := db.Seed(data)
		if err != nil {
			return "", err
		}

		return res, nil
	})

	return dbo
}

func Query(client DBClient, data DataSet) *DBOperation {

	dbo := newDBOperation("query", client, data, func(db DBClient, data DataSet) (interface{}, error) {

		log.Println(string(ToJSON(data)))

		q := DBQueryDataSet{}
		_ = json.Unmarshal(ToJSON(data), &q)

		log.Println(q)

		res, err := db.Query(q.Query)
		if err != nil {
			return "", err
		}

		return res, nil
	})

	return dbo
}
