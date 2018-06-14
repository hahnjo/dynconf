// SPDX-License-Identifier:	GPL-3.0-or-later

package dynconf

import (
	"testing"
)

func apply(r Recipe, input string) string {
	return string(ApplyToInput(r, []byte(input)))
}

func TestApply_Delete(t *testing.T) {
	r := Recipe{
		Delete: []DeleteEntry{
			{Search: "remove"},
		},
	}
	r.Compile()

	s := apply(r, "line1\nremove\nline2\ntest remove test\n")
	if s != "line1\nline2\n" {
		t.Errorf("'remove' lines should have been removed: %s", s)
	}
}

func TestApply_Replace(t *testing.T) {
	r := Recipe{
		Replace: []ReplaceEntry{
			{Search: "search", Replace: "replace"},
		},
	}
	r.Compile()

	s := apply(r, "line1\nsearch\nline2\ntest search replace\n")
	if s != "line1\nreplace\nline2\ntest replace replace\n" {
		t.Errorf("'search' lines should have been replaced: %s", s)
	}
}

func TestApply_Append(t *testing.T) {
	r := Recipe{
		Append: "append",
	}

	s := apply(r, "line\n")
	if s != "line\nappend\n" {
		t.Errorf("line should have been appended: %s", s)
	}

	// Apply should add a newline after the input.
	s = apply(r, "line")
	if s != "line\nappend\n" {
		t.Errorf("line should have been appended after newline: %s", s)
	}

	// Apply should also handle empty files.
	s = apply(r, "")
	if s != "append\n" {
		t.Errorf("line should have been appended: %s", s)
	}

	s = apply(r, "\n")
	if s != "append\n" {
		t.Errorf("line should have been appended: %s", s)
	}
}

func TestApply_Empty(t *testing.T) {
	r := Recipe{}

	s := apply(r, "")
	if s != "" {
		t.Errorf("result should be empty: %s", s)
	}

	s = apply(r, "\n")
	if s != "\n" {
		t.Errorf("result should be empty: %s", s)
	}
}

func TestApply_Newline(t *testing.T) {
	r := Recipe{}

	i := "line1\nline2\r"
	s := apply(r, i)
	if s != i {
		t.Errorf("newlines should have been copied: %s", s)
	}

	i = "line1\r\nline2\n\r"
	s = apply(r, i)
	if s != i {
		t.Errorf("newlines should have been copied: %s", s)
	}
}
