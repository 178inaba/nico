package nico

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

// Press is unknown data.
type Press struct {
	DisplayLines int64 `xml:"display_lines"`
	DisplayTime  int64 `xml:"display_time"`

	// TODO
	StyleConf interface{} `xml:"style_conf"`
}

// Telop is unknown data.
type Telop struct {
	Enable int64 `xml:"enable"`
}

// ContentsList is a list of contents.
type ContentsList struct {
	// TODO slice?
	Contents Contents `xml:"contents"`
}

// Contents is detailed information of contents such as URL of RTMP etc.
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

// Rtmp is information on RTMP.
type Rtmp struct {
	IsFms     int64  `xml:"is_fms,attr"`
	RtmptPort int64  `xml:"rtmpt_port,attr"`
	URL       string `xml:"url"`
	Ticket    string `xml:"ticket"`
}

// Ms is comment server information.
type Ms struct {
	Addr   string `xml:"addr"`
	Port   int64  `xml:"port"`
	Thread int64  `xml:"thread"`
}

// Twitter is Twitter setting information of the niconico live.
type Twitter struct {
	LiveEnabled  int64  `xml:"live_enabled"`
	VipModeCount int64  `xml:"vip_mode_count"`
	LiveAPIURL   string `xml:"live_api_url"`
}

// Player is the setting information of the player.
type Player struct {
	QosAnalytics                 int64       `xml:"qos_analytics"`
	DialogImage                  DialogImage `xml:"dialog_image"`
	IsNoticeViewerBalloonEnabled int64       `xml:"is_notice_viewer_balloon_enabled"`
	ErrorReport                  int64       `xml:"error_report"`
}

// DialogImage is the URL of the dialog image displayed to the player.
type DialogImage struct {
	Oidashi string `xml:"oidashi"`
}

// Marquee is information related to the game etc.
type Marquee struct {
	Category         string `xml:"category"`
	GameKey          string `xml:"game_key"`
	GameTime         int64  `xml:"game_time"`
	ForceNicowariOff int64  `xml:"force_nicowari_off"`
}

// Error stores the error code if Status of PlayerStatus is not ok.
type Error struct {
	Code string `xml:"code"`
}
