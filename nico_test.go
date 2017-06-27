package nico

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestLogin(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "user_session", Value: "deleted", Path: "/", Expires: time.Now()})
		http.SetCookie(w, &http.Cookie{Name: "user_session", Value: "foobarbaz", Path: "/", Expires: time.Now().AddDate(0, 1, 0), Domain: ".nicovideo.jp"})
		http.Redirect(w, r, "http://example.com", http.StatusFound)
	}))
	defer ts.Close()

	c := &Client{loginRawurl: ts.URL}
	us, err := c.Login(context.Background(), "foo@foo.com", "bar")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if us != "foobarbaz" {
		t.Fatalf("want %q but %q", "foobarbaz", us)
	}
	if c.UserSession != "foobarbaz" {
		t.Fatalf("want %q but %q", "foobarbaz", c.UserSession)
	}
}

func TestGetPlayerStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		us, err := r.Cookie("user_session")
		if err != nil {
			t.Fatalf("should not be fail: %v", err)
		}
		if us.Value != "user-session" {
			t.Fatalf("want %q but %q", "user-session", us.Value)
		}
		v := r.URL.Query().Get("v")
		if v != "lv123456789" {
			t.Fatalf("want %q but %q", "lv123456789", v)
		}
		io.WriteString(w, `<?xml version="1.0" encoding="utf-8"?><getplayerstatus status="ok" time="1234567890"><stream><title>test-title</title></stream></getplayerstatus>`)
	}))
	defer ts.Close()

	c := &Client{
		liveBaseRawurl: ts.URL,
		UserSession:    "user-session",
	}
	ps, err := c.GetPlayerStatus(context.Background(), "lv123456789")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if ps.Status != "ok" {
		t.Fatalf("want %q but %q", "ok", ps.Status)
	}
	if ps.Time != 1234567890 {
		t.Fatalf("want %d but %d", 1234567890, ps.Time)
	}
	if ps.Stream.Title != "test-title" {
		t.Fatalf("want %q but %q", "test-title", ps.Stream.Title)
	}
}

func TestGetCommunityIDFromLiveID(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "lv123456789") {
			t.Fatalf("%q should contain %q", r.URL.Path, "lv123456789")
		}
		io.WriteString(w, `<div class="shosai"><a href="http://com.nicovideo.jp/community/co1234567">foo</a></div>`)
	}))
	defer ts.Close()

	c := &Client{liveBaseRawurl: ts.URL}
	communityID, err := c.GetCommunityIDFromLiveID(context.Background(), "lv123456789")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if communityID != "co1234567" {
		t.Fatalf("want %q but %q", "co1234567", communityID)
	}
}
