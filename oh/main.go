package main 

import (
  "fmt"
  "os"
  "os/user"
  "log"
  "path"

  "github.com/olebedev/config"
  "github.com/dereulenspiegel/openhab-cli/openhab"
)

var OpenHabCfg *config.Config

var OpenHABHost string

func main() {
  LoadConfigFromHome()
  client := openhab.NewClient(OpenHABHost)
  args := os.Args
  if len(args) < 2 {
    fmt.Printf("Wrong usage\n")
    return
  }
  verb := args[1]

  switch verb {
  case "state":
    if len(args) < 3 {
      log.Fatalf("Print usage")
    }
    item := args[2]
    state, err := client.GetState(item)
    if err != nil {
      HandleError(err)
    } else {
      fmt.Print(state)
    }
  case "command":
    if len(args) < 4 {
      // TODO print usage
      log.Fatalf("Not enough arguments")
    }
    item := args[2]
    command := args[3]
    err := client.SendCommand(item, command)
    if err != nil {
      HandleError(err)
    }
  case "list":
    items, err := client.ListItems()
    if err != nil {
      HandleError(err)
    }
    for _, item := range items {
      fmt.Printf("%s\t\t%s\n", item.Name, item.State)
    }
  default:
    log.Fatalf("Unknown command %s",verb)
  }
}

func HandleError(err error) {
  log.Fatalf("Error: %v", err)
}

func LoadConfigFromHome() {
  usr, err := user.Current()
  if err != nil {
    log.Fatal("Can't obtain current user")
  }
  OpenHabCfg, err = config.ParseYamlFile(path.Join(usr.HomeDir,".openhab-cli"))
  if err != nil {
    log.Fatal("Can't parse config")
  }
  OpenHABHost, err = OpenHabCfg.String("host")
  if err != nil {
    log.Fatal("Missing config key host")
  }
}

