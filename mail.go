package nico

import "strings"

// Comment color.
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

// Comment size.
const (
	SizeMedium = "medium"
	SizeBig    = "big"
	SizeSmall  = "small"
)

var validateCommentColorMap = map[string]bool{
	CommentColorWhite:  true,
	CommentColorRed:    true,
	CommentColorPink:   true,
	CommentColorOrange: true,
	CommentColorYellow: true,
	CommentColorGreen:  true,
	CommentColorCyan:   true,
	CommentColorBlue:   true,
	CommentColorPurple: true,
}

var validateSizeMap = map[string]bool{
	SizeMedium: true,
	SizeBig:    true,
	SizeSmall:  true,
}

// Mail is a structure that specifies comment options.
type Mail struct {
	Is184        bool
	CommentColor string
	Size         string
}

func (m Mail) String() string {
	var strs []string
	if m.Is184 {
		strs = append(strs, "184")
	}
	if validateCommentColorMap[m.CommentColor] {
		strs = append(strs, m.CommentColor)
	}
	if validateSizeMap[m.Size] {
		strs = append(strs, m.Size)
	}
	return strings.Join(strs, " ")
}
