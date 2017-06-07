package nico

import (
	"context"
	"net/http"
	"net/http/httptest"
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
	if err := c.Login(context.Background(), "foo@foo.com", "bar"); err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if c.userSession != "foobarbaz" {
		t.Fatalf("want %q but %q", "foobarbaz", c.userSession)
	}
}
