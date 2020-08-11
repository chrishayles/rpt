package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

/*

	--config-file
	--continuous
	--api

*/

func main() {
	app := cli.NewApp()
	app.Name = "My App Name"
	app.Usage = "This app does something"

	myFlags := []cli.Flag{
		&cli.StringFlag{
			Name:  "primaryhost",
			Value: "localhost",
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:  "test",
			Usage: "runs a test",
			Flags: myFlags,
			Action: func(c *cli.Context) error {
				//do stuff
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
