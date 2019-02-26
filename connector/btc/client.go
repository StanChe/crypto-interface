package btc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Client - http client for node communication
type (
	Client struct {
		URL        string
		HTTPClient *http.Client
	}

	rpcResponse struct {
		Result json.RawMessage `json:"result"`
		Error  json.RawMessage `json:"error"`
	}
)

// NewClient creates new Client instance
func NewClient(rpcURL string, timeout int) *Client {
	client := Client{}
	httpClient := http.DefaultClient
	httpClient.Timeout = time.Duration(timeout) * time.Second
	client.HTTPClient = httpClient
	client.URL = rpcURL
	return &client
}

func (c *Client) send(data string) (string, error) {

	body := bytes.NewBuffer([]byte(data))
	resp, err := c.HTTPClient.Post(c.URL, "application/json", body)
	if resp != nil && resp.StatusCode >= http.StatusMultipleChoices {
		err = fmt.Errorf("http status: %s (%d)", resp.Status, resp.StatusCode)
	}
	if err != nil {
		return "", fmt.Errorf("coreclient.send.http: %s", err.Error())
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("coreclient.send.ReadAll: %s", err.Error())
	}
	var res rpcResponse
	// check the request error status
	err = json.Unmarshal(respBytes, &res)
	if err != nil {
		return "", err
	}
	if bytes.Compare(res.Error, []byte("null")) != 0 {
		return "", fmt.Errorf("%s", string(res.Error))
	}
	if string(res.Result) == "" {
		return "", fmt.Errorf("core response is empty")
	}
	return string(res.Result), nil
}
