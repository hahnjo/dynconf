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

	s := apply(r, "line1\nremove\n\nline2\ntest remove test\n")
	if s != "line1\n\nline2\n" {
		t.Errorf("'remove' lines should have been removed: %s", s)
	}
}

func TestApply_DeleteContext(t *testing.T) {
	r := Recipe{
		Delete: []DeleteEntry{
			{Context: Context{Begin: "\\[begin\\]", End: "\\[end\\]"}, Search: "removeContext"},
			{Context: Context{Begin: "\\[begin\\]"}, Search: "removeBegin"},
			{Context: Context{End: "\\[end\\]"}, Search: "removeEnd"},
			{Context: Context{End: "notThere"}, Search: "removeAlways"},
			{Context: Context{Begin: "notThere"}, Search: "removeNever"},
		},
	}
	r.Compile()

	s := apply(r, `
removeContext
removeBegin
removeEnd
removeAlways
removeNever

--[begin]--
removeContext
removeBegin
removeEnd
removeAlways
removeNever
--[end]--

removeContext
removeBegin
removeEnd
removeAlways
removeNever`)

	if s != `
removeContext
removeBegin
removeNever

--[begin]--
removeNever
--[end]--

removeContext
removeEnd
removeNever` {
		t.Errorf("some lines should have been removed: %s", s)
	}
}

func TestApply_DeleteContextSpecial(t *testing.T) {
	r := Recipe{
		Delete: []DeleteEntry{
			{Context: Context{Begin: "begin", End: "end"}, Search: "remove"},
		},
	}
	r.Compile()

	sameLine := "remove\nbegin end\nremove\n"
	s := apply(r, sameLine)
	if s != sameLine {
		t.Errorf("input should not have been modified: %s", s)
	}

	s = apply(r, "remove\nbegin\nremove\nend\nremove\nbegin\nremove\nend\nremove\n")
	if s != "remove\nbegin\nend\nremove\nbegin\nend\nremove\n" {
		t.Errorf("some lines should have been removed: %s", s)
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
func TestApply_ReplaceContext(t *testing.T) {
	r := Recipe{
		Replace: []ReplaceEntry{
			{Context: Context{Begin: "\\[begin\\]", End: "\\[end\\]"}, Search: "searchContext", Replace: "replaceContext"},
			{Context: Context{Begin: "\\[begin\\]"}, Search: "searchBegin", Replace: "replaceBegin"},
			{Context: Context{End: "\\[end\\]"}, Search: "searchEnd", Replace: "replaceEnd"},
			{Context: Context{End: "notThere"}, Search: "searchAlways", Replace: "replaceAlways"},
			{Context: Context{Begin: "notThere"}, Search: "searchNever", Replace: "replaceNever"},
		},
	}
	r.Compile()

	s := apply(r, `
searchContext
searchBegin
searchEnd
searchAlways
searchNever

--[begin]--
searchContext
searchBegin
searchEnd
searchAlways
searchNever
--[end]--

searchContext
searchBegin
searchEnd
searchAlways
searchNever`)

	if s != `
searchContext
searchBegin
replaceEnd
replaceAlways
searchNever

--[begin]--
replaceContext
replaceBegin
replaceEnd
replaceAlways
searchNever
--[end]--

searchContext
replaceBegin
searchEnd
replaceAlways
searchNever` {
		t.Errorf("some lines should have been replaced: %s", s)
	}
}

func TestApply_ReplaceContextSpecial(t *testing.T) {
	r := Recipe{
		Replace: []ReplaceEntry{
			{Context: Context{Begin: "begin", End: "end"}, Search: "search", Replace: "replace"},
		},
	}
	r.Compile()

	sameLine := "search\nbegin end\nsearch\n"
	s := apply(r, sameLine)
	if s != sameLine {
		t.Errorf("input should not have been modified: %s", s)
	}

	s = apply(r, "search\nbegin\nsearch\nend\nsearch\nbegin\nsearch\nend\nsearch\n")
	if s != "search\nbegin\nreplace\nend\nsearch\nbegin\nreplace\nend\nsearch\n" {
		t.Errorf("some lines should have been replaced: %s", s)
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

	// Apply should add a newline after the input (if there is none).
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

func TestApply_AppendNewLine(t *testing.T) {
	r := Recipe{
		Append: "line1\nline2\n",
	}

	// There should be no additiona newline after the append.
	s := apply(r, "original\n")
	if s != "original\nline1\nline2\n" {
		t.Errorf("lines should have been appended: %s", s)
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
