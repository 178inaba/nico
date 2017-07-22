package nico_test

import (
	"fmt"
	"log"

	"github.com/178inaba/nico"
)

func ExampleFindLiveID() {
	liveID, err := nico.FindLiveID("http://live.nicovideo.jp/watch/lv1234567?ref=community")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(liveID)
	// Output: lv1234567
}
