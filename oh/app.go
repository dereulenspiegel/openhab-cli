package main

import (
	"log"
	"os/user"
	"path"

	"github.com/codegangsta/cli"
	"github.com/dereulenspiegel/openhab-cli/oh/commands"
	"github.com/dereulenspiegel/openhab-cli/openhab"
	"github.com/olebedev/config"
)

var App *cli.App

var OpenHabCfg *config.Config

var OpenHABHost string

var OpenHABClient *openhab.Client

func LoadConfigFromHome() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal("Can't obtain current user")
	}
	OpenHabCfg, err = config.ParseYamlFile(path.Join(usr.HomeDir, ".openhab-cli"))
	if err != nil {
		log.Fatal("Can't parse config")
	}
	OpenHABHost, err = OpenHabCfg.String("host")
	if err != nil {
		log.Fatal("Missing config key host")
	}
}

func CollectCommands() []cli.Command {
	commandList := make([]cli.Command, 0, 10)
	commandList = append(commandList, commands.NewListCommand(OpenHABClient))
	commandList = append(commandList, commands.NewSendCommand(OpenHABClient))

	return commandList
}

func CreateCliApp() *cli.App {
	LoadConfigFromHome()
	OpenHABClient = openhab.NewClient(OpenHABHost)
	App = cli.NewApp()
	App.Name = "openHAB CLI Client"
	App.Usage = "oh [cli command] [item] [item command]"
	App.Commands = CollectCommands()

	return App
}
