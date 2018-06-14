// SPDX-License-Identifier:	GPL-3.0-or-later

package dynconf

import (
	"testing"
)

func apply(c Config, input string) string {
	return string(ApplyToInput(c, []byte(input)))
}

func TestApply_Delete(t *testing.T) {
	c := Config{
		Delete: []DeleteEntry{
			{Search: "remove"},
		},
	}
	c.Compile()

	s := apply(c, "line1\nremove\nline2\ntest remove test\n")
	if s != "line1\nline2\n" {
		t.Errorf("'remove' lines should have been removed: %s", s)
	}
}

func TestApply_Replace(t *testing.T) {
	c := Config{
		Replace: []ReplaceEntry{
			{Search: "search", Replace: "replace"},
		},
	}
	c.Compile()

	s := apply(c, "line1\nsearch\nline2\ntest search replace\n")
	if s != "line1\nreplace\nline2\ntest replace replace\n" {
		t.Errorf("'search' lines should have been replaced: %s", s)
	}
}

func TestApply_Append(t *testing.T) {
	c := Config{
		Append: "append",
	}

	s := apply(c, "line\n")
	if s != "line\nappend\n" {
		t.Errorf("line should have been appended: %s", s)
	}

	// Apply should add a newline after the input.
	s = apply(c, "line")
	if s != "line\nappend\n" {
		t.Errorf("line should have been appended after newline: %s", s)
	}

	// Apply should also handle empty files.
	s = apply(c, "")
	if s != "append\n" {
		t.Errorf("line should have been appended: %s", s)
	}

	s = apply(c, "\n")
	if s != "append\n" {
		t.Errorf("line should have been appended: %s", s)
	}
}

func TestApply_Empty(t *testing.T) {
	c := Config{}

	s := apply(c, "")
	if s != "" {
		t.Errorf("result should be empty: %s", s)
	}

	s = apply(c, "\n")
	if s != "\n" {
		t.Errorf("result should be empty: %s", s)
	}
}

func TestApply_Newline(t *testing.T) {
	c := Config{}

	i := "line1\nline2\r"
	s := apply(c, i)
	if s != i {
		t.Errorf("newlines should have been copied: %s", s)
	}

	i = "line1\r\nline2\n\r"
	s = apply(c, i)
	if s != i {
		t.Errorf("newlines should have been copied: %s", s)
	}
}
