package nico

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

// Client is a API client for niconico.
type Client struct {
	http.Client
	loginRawurl string
	userSession string
}

// NewClient return new niconico client.
func NewClient() *Client {
	return &Client{
		loginRawurl: "https://secure.nicovideo.jp/secure/login",
	}
}

// Login is login to niconico and get user session.
func (c *Client) Login(ctx context.Context, mail, password string) error {
	v := url.Values{}
	v.Set("mail", mail)
	v.Set("password", password)

	req, err := http.NewRequest(http.MethodPost, c.loginRawurl, strings.NewReader(v.Encode()))
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	cr := c.CheckRedirect
	defer func() { c.CheckRedirect = cr }()
	c.CheckRedirect = func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	cookies := resp.Cookies()
	for i := len(cookies) - 1; i >= 0; i-- {
		if cookies[i].Name == "user_session" {
			c.userSession = cookies[i].Value
			return nil
		}
	}

	return errors.New("user session not found")
}
