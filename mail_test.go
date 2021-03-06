package nico

import "testing"

func TestMail_String(t *testing.T) {
	tests := []struct {
		in  Mail
		out string
	}{
		{Mail{}, ""},
		{Mail{Is184: true}, "184"},
		{Mail{CommentColor: CommentColorRed}, "red"},
		{Mail{CommentColor: "fail"}, ""},
		{Mail{Size: SizeSmall}, "small"},
		{Mail{Size: "fail"}, ""},
		{Mail{Position: PositionUe}, "ue"},
		{Mail{Position: "fail"}, ""},
		{Mail{Is184: true, CommentColor: CommentColorPink}, "184 pink"},
		{Mail{Is184: true, CommentColor: CommentColorPink, Size: SizeBig}, "184 pink big"},
		{Mail{Is184: true, CommentColor: CommentColorPink, Size: SizeBig, Position: PositionShita}, "184 pink big shita"},
	}
	for _, tt := range tests {
		if tt.in.String() != tt.out {
			t.Fatalf("want %q but %q", tt.out, tt.in.String())
		}
	}
}
