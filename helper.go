package nico

import (
	"errors"
	"regexp"
)

var liveIDRE = regexp.MustCompile(`lv\d+`)

// FindLiveID find for the live id from s.
func FindLiveID(s string) (string, error) {
	liveID := liveIDRE.Copy().FindString(s)
	if liveID == "" {
		return "", errors.New("live id not found")
	}
	return liveID, nil
}
