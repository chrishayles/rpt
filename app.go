package main

import (
	"fmt"

	"github.com/haylesnortal/rpt/rpt"
)

func main() {

	r, err := rpt.NewRptFromEnvironment()
	if err != nil {
		fmt.Printf("error: %s", err)
	}

	r.Init()
}
