package nico

import "strings"

type Mail struct {
	Is184 bool
}

func (m Mail) String() string {
	var strs []string
	if m.Is184 {
		strs = append(strs, "184")
	}
	return strings.Join(strs, " ")
}
