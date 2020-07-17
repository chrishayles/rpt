package main

import (
	"fmt"

	"github.com/haylesnortal/rpt"
)

func main() {

	file := "/path/to/sample_data_01.json"

	ds, errs := rpt.ImportDBDataSet(file)
	if errs != nil {
		fmt.Print(errs)
	}

	operation := rpt.SeedData(rpt.NewPostgresClient(), ds)

	fmt.Println(operation)

	operation.Start()

	fmt.Println(operation.GetResult())

}
