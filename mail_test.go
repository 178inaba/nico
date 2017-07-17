package nico

import "testing"

func TestMail_String(t *testing.T) {
	m := Mail{}
	if m.String() != "" {
		t.Fatalf("want %q but %q", "", m.String())
	}

	m = Mail{Is184: true}
	if m.String() != "184" {
		t.Fatalf("want %q but %q", "184", m.String())
	}

	m = Mail{CommentColor: CommentColorRed}
	if m.String() != "red" {
		t.Fatalf("want %q but %q", "red", m.String())
	}

	m = Mail{Is184: true, CommentColor: CommentColorPink}
	if m.String() != "184 pink" {
		t.Fatalf("want %q but %q", "184 pink", m.String())
	}
}
