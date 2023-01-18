package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/ebauman/ldifgen/cmd"
	"github.com/urfave/cli/v2"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	app := &cli.App{
		Name:    "ldifgen",
		Authors: []*cli.Author{{Name: "Eamon Bauman", Email: "eamon@eamonbauman.com"}},
		Usage:   "Generate LDIF files with complex structures",
		Commands: []*cli.Command{
			cmd.GenerateCommand(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
