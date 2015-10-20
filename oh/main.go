package main 

import (
  "net/http"
  "fmt"
  "os"
  "os/user"
  "bytes"
  "io"
  "io/ioutil"
  "log"
  "path"

  "github.com/olebedev/config"
)

var OpenHabCfg *config.Config

var OpenHABHost string

func main() {
  LoadConfigFromHome()
  args := os.Args
  if len(args) < 3 {
    fmt.Printf("Wrong usage\n")
    return
  }
  verb := args[1]
  item := args[2]

  switch verb {
  case "state":
    state, err := GetState(OpenHABHost,item)
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
    command := args[3]
    err := SendCommand(OpenHABHost, item, command)
    if err != nil {
      HandleError(err)
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

func NewOpenHabRequest(method, url, contentType, payload string) (*http.Request, error) {
  var payloadReader io.Reader = nil
  if payload != "" {
    payloadReader = bytes.NewReader([]byte(payload))
  }
  req, err := http.NewRequest(method, url, payloadReader)
  if err != nil {
    return nil,err
  }
  if contentType != "" {
    req.Header.Add("Content-Type",contentType)
  }
  if payload != "" {
    req.Header.Add("Accept", "application/json")
  }
  return req, nil
}

func ReaderToString(reader io.ReadCloser) string {
  defer reader.Close()
  body, _ := ioutil.ReadAll(reader)
  return string(body)
}

func GetState(host, item string) (string,error) {
  client := &http.Client{}
  url := fmt.Sprintf("%s/rest/items/%s/state",host, item)
  req, err := NewOpenHabRequest("GET", url, "","")
  if err != nil {
    return "",err
  }
  resp, err := client.Do(req)
  if err != nil {
    return "", nil
  }

  if resp.StatusCode >= 200 && resp.StatusCode < 300 {
    body := ReaderToString(resp.Body)
    return body, nil
  } else {
    errorBody := ReaderToString(resp.Body)
    return "", fmt.Errorf("Error response %d %s",resp.StatusCode, errorBody)
  }

}

func SendCommand(host, item, command string) error {
  client := &http.Client{}
  req, err := NewOpenHabRequest("POST", fmt.Sprintf("%s/rest/items/%s", host, item), "text/plain",command)
  if err != nil {
    return err
  }
  resp, err := client.Do(req)
  if err != nil {
    return err
  }
  if resp.StatusCode < 200 || resp.StatusCode >= 300 {
    log.Fatalf("Received error from server %d %s", resp.StatusCode, ReaderToString(resp.Body))
  }
  return nil
}