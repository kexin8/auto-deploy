package test

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"testing"
)

func TestCli(t *testing.T) {
	app := &cli.App{
		Name:  "boom",
		Usage: "make an explosive entrance",
		Action: func(*cli.Context) error {
			fmt.Println("boom! I say!")
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
