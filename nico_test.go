package nico

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogin(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Set-Cookie", "user_session=deleted; Max-Age=0; Expires=Wed, 07 Jun 2017 09:43:04 GMT; Path=/")
		w.Header().Set("Set-Cookie", "user_session=foobarbaz; Max-Age=2591999; Expires=Fri, 07 Jul 2017 09:43:03 GMT; Path=/; Domain=.nicovideo.jp")
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
