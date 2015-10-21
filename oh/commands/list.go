package commands

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/dereulenspiegel/openhab-cli/openhab"
)

type ListCommand struct {
	Client *openhab.Client
}

func NewListCommand(client *openhab.Client) cli.Command {
	listCommand := &ListCommand{
		Client: client,
	}
	command := cli.Command{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "list availavle items",
		Action:  listCommand.Action,
	}

	return command
}

func (l *ListCommand) Action(c *cli.Context) {
	items, err := l.Client.ListItems()
	if err != nil {
		PrintErrorMessageAndExit(err)
	}
	for _, item := range items {
		fmt.Printf("%s\n", item.Name)
	}
}
