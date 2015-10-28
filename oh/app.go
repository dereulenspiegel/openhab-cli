package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"strings"

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

func AskForHostAndRenderConfig() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal("Can't obtain current user")
	}
	configPath := path.Join(usr.HomeDir, ".openhab-cli")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Please enter openHAB host URL:\n")
		host, _ := reader.ReadString('\n')
		cfg := map[string]interface{}{
			"host": host,
		}
		cfgContent, err := config.RenderYaml(cfg)
		if err != nil {
			log.Fatal("Can't render config")
		}
		cfgContent = strings.Replace(cfgContent, "|", "", -1)
		cfgContent = strings.Replace(cfgContent, "\n", "", -1)
		cfgContent = strings.TrimSpace(cfgContent)
		err = ioutil.WriteFile(configPath, []byte(cfgContent), 0644)
		if err != nil {
			log.Fatal("Can't write config")
		}
	}
}

func CollectCommands() []cli.Command {
	commandList := make([]cli.Command, 0, 10)
	commandList = append(commandList, commands.NewListCommand(OpenHABClient))
	commandList = append(commandList, commands.NewSendCommand(OpenHABClient))
	commandList = append(commandList, commands.NewStateCommand(OpenHABClient))

	return commandList
}

func CreateCliApp() *cli.App {
	AskForHostAndRenderConfig()
	LoadConfigFromHome()
	OpenHABClient = openhab.NewClient(OpenHABHost)
	App = cli.NewApp()
	App.Name = "openHAB CLI Client"
	App.Usage = "oh [cli command] [item] [item command]"
	App.Commands = CollectCommands()

	return App
}
