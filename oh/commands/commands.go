package commands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

type OpenHABCommand interface {
	Action(c *cli.Context)
}

func PrintStringMessageAndExit(err string) {
	fmt.Printf("Error: %s", err)
	os.Exit(1)
}

func PrintErrorMessageAndExit(err error) {
	PrintStringMessageAndExit(err.Error())
}
