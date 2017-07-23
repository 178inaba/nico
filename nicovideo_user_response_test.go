package nico

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUserInfo(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID != "2525" {
			_, err := io.WriteString(w, `<?xml version="1.0" encoding="UTF-8"?>
<nicovideo_user_response status="fail"><error><code>NOT_FOUND</code><description>user not found</description></error></nicovideo_user_response>`)
			if err != nil {
				t.Fatalf("should not be fail: %v", err)
			}
			return
		}
		_, err := io.WriteString(w, `<?xml version="1.0" encoding="UTF-8"?>
<nicovideo_user_response status="ok">
  <user>
    <id>2525</id>
    <nickname>foo</nickname>
    <thumbnail_url>http://example.com/icon.jpg</thumbnail_url>
  </user>
</nicovideo_user_response>`)
		if err != nil {
			t.Fatalf("should not be fail: %v", err)
		}
	}))
	defer ts.Close()

	c := &Client{ceBaseRawurl: ts.URL}
	_, err := c.GetUserInfo(context.Background(), 0)
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	if err.Error() != "fail: NOT_FOUND: user not found" {
		t.Fatalf("want %q but %q", "fail: NOT_FOUND: user not found", err)
	}
	uie, ok := err.(UserInfoError)
	if !ok {
		t.Fatalf("should be assertion to UserInfoError: %T", err)
	}
	if uie.Status != "fail" {
		t.Fatalf("want %q but %q", "fail", uie.Status)
	}
	if uie.Code != "NOT_FOUND" {
		t.Fatalf("want %q but %q", "NOT_FOUND", uie.Code)
	}
	if uie.Description != "user not found" {
		t.Fatalf("want %q but %q", "user not found", uie.Description)
	}

	ui, err := c.GetUserInfo(context.Background(), 2525)
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if ui.ID != 2525 {
		t.Fatalf("want %d but %d", 2525, ui.ID)
	}
	if ui.Nickname != "foo" {
		t.Fatalf("want %q but %q", "foo", ui.Nickname)
	}
	if ui.ThumbnailURL != "http://example.com/icon.jpg" {
		t.Fatalf("want %q but %q", "http://example.com/icon.jpg", ui.ThumbnailURL)
	}
}
