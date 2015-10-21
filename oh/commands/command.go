package commands

import (
	"github.com/codegangsta/cli"
	"github.com/dereulenspiegel/openhab-cli/openhab"
)

type OpenHABSendCommand struct {
	Client *openhab.Client
}

func NewSendCommand(client *openhab.Client) cli.Command {
	sendCommand := &OpenHABSendCommand{Client: client}

	return cli.Command{
		Name:    "command",
		Aliases: []string{"c"},
		Usage:   "command <item name> <command>",
		Action:  sendCommand.Action,
	}
}

func (o *OpenHABSendCommand) Action(c *cli.Context) {
	if len(c.Args()) != 2 {
		PrintStringMessageAndExit("Item name or command missing")
	}
	item := c.Args().Get(0)
	command := c.Args().Get(1)
	err := o.Client.SendCommand(item, command)
	if err != nil {
		PrintErrorMessageAndExit(err)
	}
}
