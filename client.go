package tw

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	BaseURL    *url.URL
	UserAgent  string
	Token      string
	httpClient *http.Client
}

type ClientError struct {
	Err string
}

func NewClientError(err string) error {
	return &ClientError{Err: err}
}

func (err ClientError) Error() string {
	return err.Err
}

type ServerError struct {
	StatusCode int
	Err string `json:"err"`
}

func NewApiError(statusCode int, err string) error {
	return &ServerError{
		StatusCode: statusCode,
		Err: err,
	}
}

func NewApiErrorFromResponse(statusCode int, body []byte) error {
	apiErr := ServerError{StatusCode: statusCode}
	decodeErr := json.Unmarshal(body, &apiErr)
	if decodeErr != nil {
		body := body[:256]
		apiErr.Err = string(body)
	}
	return &apiErr
}

func (err ServerError) Error() string {
	return fmt.Sprintf("api error: %d: %s", err.StatusCode, err.Err)
}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, NewApiErrorFromResponse(resp.StatusCode, body)
	}
	if v == nil {
		// Don't bother unmarshalling response if none is expected
		return resp, nil
	}
	err = json.Unmarshal(body, v)
	return resp, err
}

func NewClient(apiUrl string) (*Client, error) {
	c := &Client{}
	u, err := url.Parse(apiUrl)
	if err != nil {
		return nil, err
	}
	c.BaseURL = u
	c.httpClient = &http.Client{}
	return c, nil
}

type TokenResponse struct {
	Token string `json:"token"`
}

func (c *Client) Login(username, password string) error {
	body := map[string]string{
		"username": username,
		"password": password,
	}
	req, err := c.newRequest("POST", "authenticate", body)
	if err != nil {
		return err
	}
	var token TokenResponse
	_, err = c.do(req, &token)
	if err != nil {
		return err
	}
	c.Token = token.Token
	return nil
}

func (c *Client) Logout() error {
	req, err := c.newRequest("POST", "logout", nil)
	if err != nil {
		return err
	}
	_, err = c.do(req, nil)
	return err
}

