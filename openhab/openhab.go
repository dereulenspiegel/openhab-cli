package openhab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type Client struct {
	httpClient *http.Client
	host       string
}

type Item struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	State string `json:"state"`
	Link  string `json:"link"`
}

type ItemsResponse struct {
	Items []Item `json:"item"`
}

func NewClient(host string) *Client {
	return &Client{
		httpClient: &http.Client{},
		host:       host,
	}
}

func (c *Client) GetState(item string) (string, error) {
	url := fmt.Sprintf("%s/rest/items/%s/state", c.host, item)
	req, err := newRequest("GET", url, "", "")
	if err != nil {
		return "", nil
	}
	return c.executeRequest(req)
}

func (c *Client) SendCommand(item, command string) error {
	url := fmt.Sprintf("%s/rest/items/%s", c.host, item)
	req, err := newRequest("POST", url, "text/plain", command)
	if err != nil {
		log.Printf("Error creating request")
		return err
	}
	_, err = c.executeRequest(req)
	return err
}

func (c *Client) ListItems() ([]Item, error) {
	url := fmt.Sprintf("%s/rest/items", c.host)
	req, err := newRequest("GET", url, "", "")
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	body, err := c.executeRequest(req)
	if err != nil {
		return nil, err
	}
	itemsResponse := &ItemsResponse{}
	err = json.Unmarshal([]byte(body), itemsResponse)
	if err != nil {
		return nil, err
	}
	return itemsResponse.Items, nil
}

func (c *Client) executeRequest(req *http.Request) (string, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	body := readerToString(resp.Body)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return body, nil
	} else {
		return "", fmt.Errorf("Got error from server %d", resp.StatusCode)
	}
}

func readerToString(reader io.ReadCloser) string {
	defer reader.Close()
	body, _ := ioutil.ReadAll(reader)
	return string(body)
}

func newRequest(method, url, contentType, payload string) (*http.Request, error) {
	var payloadReader io.Reader = nil
	if payload != "" {
		payloadReader = bytes.NewReader([]byte(payload))
	}
	req, err := http.NewRequest(method, url, payloadReader)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}
	if payload != "" {
		req.Header.Add("Accept", "application/json")
	}
	return req, nil
}
