package nico

import "testing"

func TestFindLiveID(t *testing.T) {
	_, err := FindLiveID("http://live.nicovideo.jp/")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
	liveID, err := FindLiveID("http://live.nicovideo.jp/watch/lv1234567?ref=community")
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if liveID != "lv1234567" {
		t.Fatalf("want %s but %s", "lv1234567", liveID)
	}
}
