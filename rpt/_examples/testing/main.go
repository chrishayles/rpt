package main

import (
	"fmt"
	"os"

	"github.com/haylesnortal/rpt"
)

func main() {

	// file := "/Users/chrishayles/go/src/github.com/haylesnortal/rpt/sample_data_01.json"

	// ds, errs := rpt.ImportDBDataSet(file)
	// if len(errs) > 0 {
	// 	fmt.Print(errs)
	// }

	// psql := rpt.NewPostgresClient()
	// psql.Host = "localhost"
	// psql.Port = 5432
	// psql.User = "postgres"
	// psql.Password = "mysecretpassword"
	// psql.SSLMode = "disable"
	// //psql.DBName = "postgres"

	// err := psql.Connect()
	// if err != nil {
	// 	fmt.Print(err)
	// }

	// operation := rpt.SeedData(psql, ds)
	// //fmt.Println(operation)
	// operation.Start()
	// res := operation.GetResult()

	// fmt.Println(res)

	os.Setenv("RPT_PRIMARY_HOST", "postgres:localhost") //required
	os.Setenv("RPT_PRIMARY_PORT", "5432")               //defaults to 5432
	os.Setenv("RPT_PRIMARY_USER", "postgres")           //required
	os.Setenv("RPT_PRIMARY_PASS", "mysecretpassword")   //required
	os.Setenv("RPT_PRIMARY_SSLMODE", "disable")         //defaults to disable

	os.Setenv("RPT_SECONDARY_HOST", "postgres:localhost") //required
	os.Setenv("RPT_SECONDARY_PORT", "5432")               //defaults to 5432
	os.Setenv("RPT_SECONDARY_USER", "postgres")           //required
	os.Setenv("RPT_SECONDARY_PASS", "mysecretpassword")   //required
	os.Setenv("RPT_SECONDARY_SSLMODE", "disable")         //defaults to disable

	os.Setenv("RPT_SEED_FILE", "/Users/chrishayles/go/src/github.com/haylesnortal/rpt/sample_data_01.json") //optional

	os.Setenv("RPT_API", "TRUE") //defaults to false
	//os.Setenv("RPT_API_BASEPATH", "/api")     //defaults to /api
	//os.Setenv("RPT_API_LISTEN_ADDR", ":5000") //defaults to [localhost]:5000

	r, err := rpt.NewRptFromEnvironment()
	if err != nil {
		fmt.Printf("error: %s", err)
	}

	r.Init()
}
