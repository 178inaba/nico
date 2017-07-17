package nico

import "strings"

const (
	CommentColorWhite  = "white"
	CommentColorRed    = "red"
	CommentColorPink   = "pink"
	CommentColorOrange = "orange"
	CommentColorYellow = "yellow"
	CommentColorGreen  = "green"
	CommentColorCyan   = "cyan"
	CommentColorBlue   = "blue"
	CommentColorPurple = "purple"
)

// Mail is a structure that specifies comment options.
type Mail struct {
	Is184        bool
	CommentColor string
}

func (m Mail) String() string {
	var strs []string
	if m.Is184 {
		strs = append(strs, "184")
	}
	if m.CommentColor != "" {
		strs = append(strs, m.CommentColor)
	}
	return strings.Join(strs, " ")
}
