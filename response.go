package zabbix

import (
	"encoding/json"
	"fmt"
)

type Response struct {
	StatusCode int `json:"-"`

	JSONRPCVersion string `json:"jsonrpc"`

	Body json.RawMessage `json:"result"`

	RequestID int `json:"id"`

	Error APIError `json:"error"`
}

func (c *Response) Err() error {
	if c.Error.Code != 0 {
		return fmt.Errorf("HTTP %d %s (%d)\n%s", c.StatusCode, c.Error.Message, c.Error.Code, c.Error.Data)
	}

	return nil
}

func (c *Response) Bind(v interface{}) error {
	err := json.Unmarshal(c.Body, v)
	if err != nil {
		return fmt.Errorf("Error decoding JSON response body: %v", err)
	}

	return nil
}
