package zabbix

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var ErrNotFound = errors.New("No results were found matching the given search parameters")

type Session struct {
	URL string `json:"url"`

	Token string `json:"token"`

	APIVersion string `json:"apiVersion"`

	client *http.Client
}

func NewSession(url string, username string, password string) (session *Session, err error) {
	session = &Session{URL: url}
	err = session.login(username, password)
	return
}

func (c *Session) login(username, password string) error {
	_, err := c.GetVersion()
	if err != nil {
		return fmt.Errorf("Failed to retrieve Zabbix API version: %v", err)
	}

	params := map[string]string{
		"user":     username,
		"password": password,
	}

	res, err := c.Do(NewRequest("user.login", params))
	if err != nil {
		return fmt.Errorf("Error logging in to Zabbix API: %v", err)
	}

	err = res.Bind(&c.Token)
	if err != nil {
		return fmt.Errorf("Error failed to decode Zabbix login response: %v", err)
	}

	return nil
}

func (c *Session) GetVersion() (string, error) {
	if c.APIVersion == "" {
		res, err := c.Do(NewRequest("apiinfo.version", nil))
		if err != nil {
			return "", err
		}

		err = res.Bind(&c.APIVersion)
		if err != nil {
			return "", err
		}
	}
	return c.APIVersion, nil
}

func (c *Session) AuthToken() string {
	return c.Token
}

func (c *Session) Do(req *Request) (resp *Response, err error) {
	req.AuthToken = c.Token

	b, err := json.Marshal(req)
	if err != nil {
		return
	}

	fmt.Printf("Call     [%s:%d]: %s\n", req.Method, req.RequestID, b)

	r, err := http.NewRequest("POST", c.URL, bytes.NewReader(b))
	if err != nil {
		return
	}
	r.ContentLength = int64(len(b))
	r.Header.Add("Content-Type", "application/json-rpc")

	client := c.client
	if client == nil {
		client = http.DefaultClient
	}
	res, err := client.Do(r)
	if err != nil {
		return
	}

	defer res.Body.Close()

	b, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response: %v", err)
	}

	fmt.Printf("Response [%s:%d]: %s\n", req.Method, req.RequestID, b)

	resp = &Response{
		StatusCode: res.StatusCode,
	}

	err = json.Unmarshal(b, &resp)
	if err != nil {
		return nil, fmt.Errorf("Error decoding JSON response body: %v", err)
	}

	if err = resp.Err(); err != nil {
		return
	}

	return
}

func (c *Session) Get(method string, params interface{}, v interface{}) error {
	req := NewRequest(method, params)
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	err = resp.Bind(v)
	if err != nil {
		return err
	}

	return nil
}
