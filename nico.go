package nico

import (
	"bufio"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Client is a API client for niconico.
type Client struct {
	http.Client
	loginRawurl         string
	liveBaseRawurl      string
	communityBaseRawurl string
	UserSession         string
}

// NewClient return new niconico client.
func NewClient() *Client {
	return &Client{
		loginRawurl:         "https://secure.nicovideo.jp/secure/login",
		liveBaseRawurl:      "http://live.nicovideo.jp",
		communityBaseRawurl: "http://com.nicovideo.jp",
	}
}

// Login is login to niconico and get user session.
func (c *Client) Login(ctx context.Context, mail, password string) (string, error) {
	v := url.Values{}
	v.Set("mail", mail)
	v.Set("password", password)

	req, err := http.NewRequest(http.MethodPost, c.loginRawurl, strings.NewReader(v.Encode()))
	if err != nil {
		return "", err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	cr := c.CheckRedirect
	defer func() { c.CheckRedirect = cr }()
	c.CheckRedirect = func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }
	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
		return "", errors.New(resp.Status)
	}

	cookies := resp.Cookies()
	for i := len(cookies) - 1; i >= 0; i-- {
		if cookies[i].Name == "user_session" {
			c.UserSession = cookies[i].Value
			return c.UserSession, nil
		}
	}

	return "", errors.New("login failed")
}

func (c *Client) GetPlayerStatus(ctx context.Context, liveID string) (*PlayerStatus, error) {
	u, err := url.Parse(c.liveBaseRawurl)
	if err != nil {
		return nil, err
	}
	u.Path = "api/getplayerstatus"

	v := url.Values{}
	v.Set("v", liveID)
	u.RawQuery = v.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.AddCookie(&http.Cookie{Name: "user_session", Value: c.UserSession})

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	ps := PlayerStatus{}
	if err := xml.NewDecoder(resp.Body).Decode(&ps); err != nil {
		return nil, err
	}
	if ps.Status != "ok" {
		return nil, PlayerStatusError{Status: ps.Status, Code: ps.Error.Code}
	}

	return &ps, nil
}

func (c *Client) GetPostkey(ctx context.Context, thread int64) (string, error) {
	u, err := url.Parse(c.liveBaseRawurl)
	if err != nil {
		return "", err
	}
	u.Path = "api/getpostkey"

	v := url.Values{}
	v.Set("thread", fmt.Sprint(thread))
	u.RawQuery = v.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return "", err
	}
	req = req.WithContext(ctx)
	req.AddCookie(&http.Cookie{Name: "user_session", Value: c.UserSession})

	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	rv, err := url.ParseQuery(string(b))
	if err != nil {
		return "", err
	}
	postkey := rv.Get("postkey")
	if postkey == "" {
		return "", errors.New("postkey is empty")
	}

	return postkey, nil
}

func (c *Client) FollowCommunity(ctx context.Context, communityID string) error {
	u, err := url.Parse(c.communityBaseRawurl)
	if err != nil {
		return err
	}
	u.Path = fmt.Sprintf("motion/%s", communityID)

	v := url.Values{}
	v.Set("mode", "commit")

	req, err := http.NewRequest(http.MethodPost, u.String(), strings.NewReader(v.Encode()))
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Referer", u.String())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "user_session", Value: c.UserSession})

	cr := c.CheckRedirect
	defer func() { c.CheckRedirect = cr }()
	c.CheckRedirect = func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
		return errors.New(resp.Status)
	} else if !strings.Contains(resp.Header.Get("Location"), "done") {
		return errors.New("community follow failed")
	}

	return nil
}

func (c *Client) getLeaveCommunityFormData(ctx context.Context, communityID string) (url.Values, error) {
	u, err := url.Parse(c.communityBaseRawurl)
	if err != nil {
		return nil, err
	}
	u.Path = fmt.Sprintf("leave/%s", communityID)

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.AddCookie(&http.Cookie{Name: "user_session", Value: c.UserSession})

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}

	v := url.Values{}
	doc.Find(".leave_form input").Each(func(i int, s *goquery.Selection) {
		name, ok := s.Attr("name")
		if !ok {
			return
		}
		value, ok := s.Attr("value")
		if !ok {
			return
		}
		v.Set(name, value)
	})

	return v, nil
}

func (c *Client) LeaveCommunity(ctx context.Context, communityID string) error {
	u, err := url.Parse(c.communityBaseRawurl)
	if err != nil {
		return err
	}
	u.Path = fmt.Sprintf("leave/%s", communityID)

	v, err := c.getLeaveCommunityFormData(ctx, communityID)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), strings.NewReader(v.Encode()))
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Referer", u.String())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "user_session", Value: c.UserSession})

	cr := c.CheckRedirect
	defer func() { c.CheckRedirect = cr }()
	c.CheckRedirect = func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
		return errors.New(resp.Status)
	} else if !strings.Contains(resp.Header.Get("Location"), "done") {
		return errors.New("community leave failed")
	}

	return nil
}

func (c *Client) MakeLiveClient(ctx context.Context, liveID string) (*LiveClient, error) {
	ps, err := c.GetPlayerStatus(ctx, liveID)
	if err != nil {
		return nil, err
	}

	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", ps.Ms.Addr, ps.Ms.Port))
	if err != nil {
		return nil, err
	}
	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	return &LiveClient{Client: c, ps: ps, conn: conn}, nil
}

type LiveClient struct {
	*Client
	ps   *PlayerStatus
	conn net.Conn
}

func (c *LiveClient) StreamingComment(ctx context.Context, resFrom int64) (chan Comment, error) {
	b, err := xml.Marshal(SendThread{Thread: c.ps.Ms.Thread, Version: 20061206, ResFrom: resFrom})
	if err != nil {
		return nil, err
	}
	b = append(b, 0)
	_, err = c.conn.Write(b)
	if err != nil {
		return nil, err
	}

	r := bufio.NewReader(c.conn)
	ch := make(chan Comment)
	go func() {
		defer close(ch)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			rb, err := r.ReadSlice(0)
			if err != nil {
				ch <- &CommentError{err}
			}
			var thread Thread
			if err := xml.Unmarshal(rb, &thread); err == nil {
				ch <- &thread
			}
			var chat Chat
			if err := xml.Unmarshal(rb, &chat); err == nil {
				ch <- &chat
			}
			var chatResult ChatResult
			if err := xml.Unmarshal(rb, &chatResult); err == nil {
				ch <- &chatResult
			}
		}
	}()
	return ch, nil
}

func (c *LiveClient) PostComment(ctx context.Context, comment string, mail Mail) error {
	postkey, err := c.GetPostkey(ctx, c.ps.Ms.Thread)
	if err != nil {
		return err
	}
	chat := SendChat{
		Vpos:    (time.Now().UnixNano() - c.ps.Stream.BaseTime*int64(time.Second)) / (int64(time.Millisecond) * 10),
		Mail:    mail.String(),
		UserID:  fmt.Sprint(c.ps.User.UserID),
		Postkey: postkey,
		Comment: comment,
	}
	b, err := xml.Marshal(chat)
	if err != nil {
		return err
	}
	b = append(b, 0)
	if _, err := c.conn.Write(b); err != nil {
		return err
	}
	return nil
}

type SendThread struct {
	XMLName xml.Name `xml:"thread"`
	Thread  int64    `xml:"thread,attr"`
	Version int64    `xml:"version,attr"`
	ResFrom int64    `xml:"res_from,attr"`
}

type Comment interface {
	comment()
}

type Thread struct {
	XMLName    xml.Name `xml:"thread"`
	Resultcode int64    `xml:"resultcode,attr"`
	Thread     int64    `xml:"thread,attr"`
	LastRes    int64    `xml:"last_res,attr"`
	Ticket     string   `xml:"ticket,attr"`
	Revision   int64    `xml:"revision,attr"`
	ServerTime int64    `xml:"server_time,attr"`
}

func (t *Thread) comment() {}

type Chat struct {
	XMLName   xml.Name `xml:"chat"`
	Thread    int64    `xml:"thread,attr"`
	No        int64    `xml:"no,attr"`
	Vpos      int64    `xml:"vpos,attr"`
	Date      int64    `xml:"date,attr"`
	DateUsec  int64    `xml:"date_usec,attr"`
	Mail      string   `xml:"mail,attr"`
	Yourpost  int64    `xml:"yourpost,attr"`
	UserID    string   `xml:"user_id,attr"`
	Premium   int64    `xml:"premium,attr"`
	Anonymity int64    `xml:"anonymity,attr"`
	Locale    string   `xml:"locale,attr"`
	Score     int64    `xml:"score,attr"`
	Comment   string   `xml:",chardata"`
}

func (c *Chat) comment() {}

type ChatResult struct {
	XMLName xml.Name `xml:"chat_result"`
	Thread  int64    `xml:"thread,attr"`
	Status  int64    `xml:"status,attr"`
	No      int64    `xml:"no,attr"`
}

func (r *ChatResult) comment() {}

type CommentError struct{ error }

func (e *CommentError) comment() {}

type SendChat struct {
	XMLName xml.Name `xml:"chat"`
	Vpos    int64    `xml:"vpos,attr"`
	Mail    string   `xml:"mail,attr,omitempty"`
	UserID  string   `xml:"user_id,attr"`
	Postkey string   `xml:"postkey,attr"`
	Comment string   `xml:",chardata"`
}

const PlayerStatusErrorCodeFull = "full"

type PlayerStatusError struct {
	Status string
	Code   string
}

func (e PlayerStatusError) Error() string {
	return fmt.Sprintf("%s: %s", e.Status, e.Code)
}
