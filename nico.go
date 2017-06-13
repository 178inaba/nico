package nico

import (
	"bufio"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
)

// Client is a API client for niconico.
type Client struct {
	http.Client
	loginRawurl    string
	liveBaseRawurl string
	userSession    string
}

// NewClient return new niconico client.
func NewClient() *Client {
	return &Client{
		loginRawurl:    "https://secure.nicovideo.jp/secure/login",
		liveBaseRawurl: "http://live.nicovideo.jp",
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

	if resp.StatusCode != http.StatusFound {
		return errors.New(resp.Status)
	}

	cookies := resp.Cookies()
	for i := len(cookies) - 1; i >= 0; i-- {
		if cookies[i].Name == "user_session" {
			c.userSession = cookies[i].Value
			return nil
		}
	}

	return errors.New("user session not found")
}

// PlayerStatus is niconico live player status.
type PlayerStatus struct {
	Status string `xml:"status,attr"`
	Time   int64  `xml:"time,attr"`
	Stream Stream `xml:"stream"`
	User   User   `xml:"user"`
	Rtmp   Rtmp   `xml:"rtmp"`
	Ms     Ms     `xml:"ms"`

	// TODO
	TidList interface{} `xml:"tid_list"`

	Twitter Twitter `xml:"twitter"`
	Player  Player  `xml:"player"`
	Marquee Marquee `xml:"marquee"`
	Error   Error   `xml:"error"`
}

// Stream is niconico live player status in player status.
type Stream struct {
	ID                       string `xml:"id"`
	Title                    string `xml:"title"`
	Description              string `xml:"description"`
	ProviderType             string `xml:"provider_type"`
	DefaultCommunity         string `xml:"default_community"`
	International            int64  `xml:"international"`
	IsOwner                  int64  `xml:"is_owner"`
	OwnerID                  int64  `xml:"owner_id"`
	OwnerName                string `xml:"owner_name"`
	IsReserved               int64  `xml:"is_reserved"`
	IsNiconicoEnqueteEnabled int64  `xml:"is_niconico_enquete_enabled"`
	WatchCount               int64  `xml:"watch_count"`
	CommentCount             int64  `xml:"comment_count"`
	BaseTime                 int64  `xml:"base_time"`
	OpenTime                 int64  `xml:"open_time"`
	StartTime                int64  `xml:"start_time"`
	EndTime                  int64  `xml:"end_time"`
	IsRerunStream            int64  `xml:"is_rerun_stream"`

	// TODO
	BourbonURL   interface{} `xml:"bourbon_url"`
	FullVideo    interface{} `xml:"full_video"`
	AfterVideo   interface{} `xml:"after_video"`
	BeforeVideo  interface{} `xml:"before_video"`
	KickoutVideo interface{} `xml:"kickout_video"`

	TwitterTag       string `xml:"twitter_tag"`
	DanjoCommentMode int64  `xml:"danjo_comment_mode"`
	InfinityMode     int64  `xml:"infinity_mode"`
	Archive          int64  `xml:"archive"`
	Press            Press  `xml:"press"`

	// TODO
	PluginDelay interface{} `xml:"plugin_delay"`
	PluginURL   interface{} `xml:"plugin_url"`
	PluginURLs  interface{} `xml:"plugin_urls"`

	AllowNetduetto               int64 `xml:"allow_netduetto"`
	NgScoring                    int64 `xml:"ng_scoring"`
	IsNonarchiveTimeshiftEnabled int64 `xml:"is_nonarchive_timeshift_enabled"`
	IsTimeshiftReserved          int64 `xml:"is_timeshift_reserved"`
	HeaderComment                int64 `xml:"header_comment"`
	FooterComment                int64 `xml:"footer_comment"`
	SplitBottom                  int64 `xml:"split_bottom"`
	SplitTop                     int64 `xml:"split_top"`
	BackgroundComment            int64 `xml:"background_comment"`

	// TODO
	FontScale interface{} `xml:"font_scale"`

	CommentLock  int64        `xml:"comment_lock"`
	Telop        Telop        `xml:"telop"`
	ContentsList ContentsList `xml:"contents_list"`
	PictureURL   string       `xml:"picture_url"`
	ThumbURL     string       `xml:"thumb_url"`

	// TODO
	IsPriorityPrefecture interface{} `xml:"is_priority_prefecture"`
}

type Press struct {
	DisplayLines int64 `xml:"display_lines"`
	DisplayTime  int64 `xml:"display_time"`

	// TODO
	StyleConf interface{} `xml:"style_conf"`
}

type Telop struct {
	Enable int64 `xml:"enable"`
}

type ContentsList struct {
	// TODO slice?
	Contents Contents `xml:"contents"`
}

type Contents struct {
	ID           string `xml:"id,attr"`
	DisableAudio int64  `xml:"disableAudio,attr"`
	DisableVideo int64  `xml:"disableVideo,attr"`
	StartTime    int64  `xml:"start_time,attr"`
	Contents     string `xml:",chardata"`
}

// User is niconico user data in player status.
type User struct {
	UserID         int64  `xml:"user_id"`
	Nickname       string `xml:"nickname"`
	IsPremium      int64  `xml:"is_premium"`
	UserAge        int64  `xml:"userAge"`
	UserSex        int64  `xml:"userSex"`
	UserDomain     string `xml:"userDomain"`
	UserPrefecture int64  `xml:"userPrefecture"`
	UserLanguage   string `xml:"userLanguage"`
	RoomLabel      string `xml:"room_label"`
	RoomSeetno     int64  `xml:"room_seetno"`

	// TODO
	IsJoin interface{} `xml:"is_join"`

	TwitterInfo TwitterInfo `xml:"twitter_info"`
}

// TwitterInfo is user's twitter info in user.
type TwitterInfo struct {
	Status          string      `xml:"status"`
	ScreenName      interface{} `xml:"screen_name"`
	FollowersCount  int64       `xml:"followers_count"`
	IsVip           int64       `xml:"is_vip"`
	ProfileImageURL string      `xml:"profile_image_url"`
	AfterAuth       int64       `xml:"after_auth"`
	TweetToken      string      `xml:"tweet_token"`
}

type Rtmp struct {
	IsFms     int64  `xml:"is_fms,attr"`
	RtmptPort int64  `xml:"rtmpt_port,attr"`
	URL       string `xml:"url"`
	Ticket    string `xml:"ticket"`
}

type Ms struct {
	Addr   string `xml:"addr"`
	Port   int64  `xml:"port"`
	Thread int64  `xml:"thread"`
}

func (m *Ms) StreamingComment(ctx context.Context, resFrom int64) (chan Comment, error) {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", m.Addr, m.Port))
	if err != nil {
		return nil, err
	}
	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	b, err := xml.Marshal(SendThread{Thread: m.Thread, Version: 20061206, ResFrom: resFrom})
	if err != nil {
		return nil, err
	}
	b = append(b, 0)
	io.WriteString(conn, string(b))

	r := bufio.NewReader(conn)
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
			var rt ReceiveThread
			if err := xml.Unmarshal(rb, &rt); err == nil {
				ch <- &rt
			}
			var chat Chat
			if err := xml.Unmarshal(rb, &chat); err == nil {
				ch <- &chat
			}
		}
	}()
	return ch, nil
}

type Twitter struct {
	LiveEnabled  int64  `xml:"live_enabled"`
	VipModeCount int64  `xml:"vip_mode_count"`
	LiveApiUrl   string `xml:"live_api_url"`
}

type Player struct {
	QosAnalytics                 int64       `xml:"qos_analytics"`
	DialogImage                  DialogImage `xml:"dialog_image"`
	IsNoticeViewerBalloonEnabled int64       `xml:"is_notice_viewer_balloon_enabled"`
	ErrorReport                  int64       `xml:"error_report"`
}

type DialogImage struct {
	Oidashi string `xml:"oidashi"`
}

type Marquee struct {
	Category         string `xml:"category"`
	GameKey          string `xml:"game_key"`
	GameTime         int64  `xml:"game_time"`
	ForceNicowariOff int64  `xml:"force_nicowari_off"`
}

type Error struct {
	Code string `xml:"code"`
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
	req.AddCookie(&http.Cookie{Name: "user_session", Value: c.userSession})

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
		return nil, fmt.Errorf("%s: %s", ps.Status, ps.Error.Code)
	}

	return &ps, nil
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

type ReceiveThread struct {
	XMLName    xml.Name `xml:"thread"`
	Resultcode int64    `xml:"resultcode,attr"`
	Thread     int64    `xml:"thread,attr"`
	LastRes    int64    `xml:"last_res,attr"`
	Ticket     string   `xml:"ticket,attr"`
	Revision   int64    `xml:"revision,attr"`
	ServerTime int64    `xml:"server_time,attr"`
}

func (t *ReceiveThread) comment() {}

type Chat struct {
	XMLName   xml.Name `xml:"chat"`
	Thread    int64    `xml:"thread,attr"`
	No        int64    `xml:"no,attr"`
	Vpos      int64    `xml:"vpos,attr"`
	Date      int64    `xml:"date,attr"`
	DateUsec  int64    `xml:"date_usec,attr"`
	Mail      string   `xml:"mail,attr"`
	UserID    string   `xml:"user_id,attr"`
	Premium   int64    `xml:"premium,attr"`
	Anonymity int64    `xml:"anonymity,attr"`
	Locale    string   `xml:"locale,attr"`
	Comment   string   `xml:",chardata"`
}

func (c *Chat) comment() {}

type CommentError struct{ error }

func (e *CommentError) comment() {}

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
	req.AddCookie(&http.Cookie{Name: "user_session", Value: c.userSession})

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
