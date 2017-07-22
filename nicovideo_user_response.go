package nico

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// NicovideoUserResponse is the response of the user information acquisition API.
type NicovideoUserResponse struct {
	Status      string        `xml:"status,attr"`
	UserInfo    UserInfo      `xml:"user"`
	VitaOption  VitaOption    `xml:"vita_option"`
	Additionals string        `xml:"additionals"`
	Error       UserInfoError `xml:"error"`
}

// UserInfo is user information.
type UserInfo struct {
	ID           int64  `xml:"id"`
	Nickname     string `xml:"nickname"`
	ThumbnailURL string `xml:"thumbnail_url"`
}

// VitaOption is an option of Vita.
type VitaOption struct {
	UserSecret int64 `xml:"user_secret"`
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

// GetNicovideoUserResponse gets the response of the user information API.
// No login required.
func (c *Client) GetNicovideoUserResponse(ctx context.Context, userID int64) (*NicovideoUserResponse, error) {
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

	nur := NicovideoUserResponse{}
	if err := xml.NewDecoder(resp.Body).Decode(&nur); err != nil {
		return nil, err
	}
	if nur.Status != "ok" {
		nur.Error.Status = nur.Status
		return nil, nur.Error
	}
	return &nur, nil
}

// GetUserInfo is get user's information.
// No login required.
func (c *Client) GetUserInfo(ctx context.Context, userID int64) (*UserInfo, error) {
	nur, err := c.GetNicovideoUserResponse(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &nur.UserInfo, nil
}
