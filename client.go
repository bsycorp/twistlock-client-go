package tw

import (
	"bytes"
	"encoding/json"
	"errors"
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

func NewServerError(statusCode int, err string) error {
	return &ServerError{
		StatusCode: statusCode,
		Err: err,
	}
}

func NewServerErrorFromResponse(statusCode int, body []byte) error {
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
		req.Header.Set("Authorization", "Bearer " + c.Token)
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
		return nil, NewServerErrorFromResponse(resp.StatusCode, body)
	}
	if v == nil {
		// Don't bother unmarshalling response if none is expected
		return resp, nil
	} else {
		err = json.Unmarshal(body, v)
		return resp, err
	}
}

func NewClient(apiUrl string) (*Client, error) {
	c := &Client{}
	u, err := url.Parse(apiUrl)
	if err != nil {
		return nil, err
	}
	c.BaseURL = u
	c.httpClient = http.DefaultClient
	return c, nil
}

type TokenResponse struct {
	Token string `json:"token"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// If not initialized, must create an initial admin account
func (c *Client) Signup(username, password string) error {
	req, err := c.newRequest("POST", "signup", &Credentials{
		Username: username,
		Password: password,
	})
	if err != nil {
		return err
	}
	_, err = c.do(req, nil)
	if err != nil {
		return err
	}
	return nil
}

// Exchange credentials for JWT
func (c *Client) Login(username, password string) error {
	req, err := c.newRequest("POST", "authenticate", &Credentials{
		Username: username,
		Password: password,
	})
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

// Invalidate JWT
func (c *Client) Logout() error {
	req, err := c.newRequest("POST", "logout", nil)
	if err != nil {
		return err
	}
	_, err = c.do(req, nil)
	return err
}

// Health check the API
func (c *Client) Ping() error {
	req, err := c.newRequest("GET", "_ping", nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return NewServerErrorFromResponse(resp.StatusCode, body)
	}
	if string(body) != "OK" {
		return errors.New("health check did not return OK")
	}
	return nil
}
