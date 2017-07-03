package nico

import "strings"

// Mail is a structure that specifies comment options.
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
