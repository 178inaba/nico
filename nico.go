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
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var findCommunityID = func() func(string) string {
	re := regexp.MustCompile(`co\d+`)
	return func(s string) string {
		return re.Copy().FindString(s)
	}
}()

// Client is a API client for niconico.
type Client struct {
	http.Client
	loginRawurl         string
	liveBaseRawurl      string
	communityBaseRawurl string
	ceBaseRawurl        string
	UserSession         string
}

// NewClient return new niconico client.
func NewClient() *Client {
	return &Client{
		loginRawurl:         "https://secure.nicovideo.jp/secure/login",
		liveBaseRawurl:      "http://live.nicovideo.jp",
		communityBaseRawurl: "http://com.nicovideo.jp",
		ceBaseRawurl:        "http://api.ce.nicovideo.jp",
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

// GetPlayerStatus gets the player status.
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

// GetPostkey gets the key to be specified when posting a comment.
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

// FollowCommunity follows community of communityID.
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

// LeaveCommunity leaves the community of communityID.
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

// GetCommunityIDFromLiveID gets the ID of the community that is broadcasting in liveID.
func (c *Client) GetCommunityIDFromLiveID(ctx context.Context, liveID string) (string, error) {
	u, err := url.Parse(c.liveBaseRawurl)
	if err != nil {
		return "", err
	}
	u.Path = path.Join("watch", liveID)

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return "", err
	}
	req = req.WithContext(ctx)

	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return "", err
	}

	href, ok := doc.Find(".shosai > a").Attr("href")
	if !ok {
		return "", errors.New("community not found")
	}

	return findCommunityID(href), nil
}

// GetUserInfo is get user's information.
func (c *Client) GetUserInfo(ctx context.Context, userID int64) (*UserInfo, error) {
	u, err := url.Parse(c.ceBaseRawurl)
	if err != nil {
		return nil, err
	}
	u.Path = "api/v1/user.info"

	v := url.Values{}
	v.Set("user_id", strconv.FormatInt(userID, 10))
	u.RawQuery = v.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	nur := nicovideoUserResponse{}
	if err := xml.NewDecoder(resp.Body).Decode(&nur); err != nil {
		return nil, err
	}
	if nur.Status != "ok" {
		nur.Error.Status = nur.Status
		return nil, nur.Error
	}

	return &nur.UserInfo, nil
}

// MakeLiveClient creates a client with broadcast information from liveID.
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

// LiveClient is a client with broadcast information.
type LiveClient struct {
	*Client
	ps   *PlayerStatus
	conn net.Conn
}

// StreamingComment return the channel that receives comment.
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

// PostComment post the comment.
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

// SendThread is an xml struct of the thread to send.
type SendThread struct {
	XMLName xml.Name `xml:"thread"`
	Thread  int64    `xml:"thread,attr"`
	Version int64    `xml:"version,attr"`
	ResFrom int64    `xml:"res_from,attr"`
}

// Comment is a interface of struct received on the chan returned by PostComment.
type Comment interface {
	comment()
}

// Thread is a struct of xml received immediately after connection.
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

// Chat is an xml struct of comment.
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

// ChatResult is an xml struct that returns the posting result of comment.
type ChatResult struct {
	XMLName xml.Name `xml:"chat_result"`
	Thread  int64    `xml:"thread,attr"`
	Status  int64    `xml:"status,attr"`
	No      int64    `xml:"no,attr"`
}

func (r *ChatResult) comment() {}

// CommentError is a struct containing error of StreamingComment function.
type CommentError struct{ error }

func (e *CommentError) comment() {}

// SendChat is a struct to use when posting comment.
type SendChat struct {
	XMLName xml.Name `xml:"chat"`
	Vpos    int64    `xml:"vpos,attr"`
	Mail    string   `xml:"mail,attr,omitempty"`
	UserID  string   `xml:"user_id,attr"`
	Postkey string   `xml:"postkey,attr"`
	Comment string   `xml:",chardata"`
}

type nicovideoUserResponse struct {
	Status   string        `xml:"status,attr"`
	UserInfo UserInfo      `xml:"user"`
	Error    UserInfoError `xml:"error"`
}

// UserInfo is user information.
type UserInfo struct {
	ID           int64  `xml:"id"`
	Nickname     string `xml:"nickname"`
	ThumbnailURL string `xml:"thumbnail_url"`
}

// UserInfoError is an error when acquiring user information.
type UserInfoError struct {
	Status      string
	Code        string `xml:"code"`
	Description string `xml:"description"`
}

func (e UserInfoError) Error() string {
	return fmt.Sprintf("%s: %s: %s", e.Status, e.Code, e.Description)
}

// Error code of PlayerStatus.
const (
	PlayerStatusErrorCodeFull                   = "full"
	PlayerStatusErrorCodeNotlogin               = "notlogin"
	PlayerStatusErrorCodeRequireCommunityMember = "require_community_member"
)

// PlayerStatusError is an error to return if Status of PlayerStatus is not ok.
type PlayerStatusError struct {
	Status string
	Code   string
}

func (e PlayerStatusError) Error() string {
	return fmt.Sprintf("%s: %s", e.Status, e.Code)
}
