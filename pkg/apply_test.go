// SPDX-License-Identifier:	GPL-3.0-or-later

package dynconf

import (
	"testing"
)

func applyNoErrors(t *testing.T, r Recipe, input string) string {
	modified, errs := ApplyToInput(r, []byte(input))
	if len(errs) != 0 {
		t.Errorf("did not expect %d errors", len(errs))
	}
	return string(modified)
}

func TestApply_Delete(t *testing.T) {
	r := Recipe{
		Delete: []DeleteEntry{
			{Search: "remove"},
		},
	}
	r.Compile()

	s := applyNoErrors(t, r, "line1\nremove\n\nline2\ntest remove test\n")
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

	s := applyNoErrors(t, r, `
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
	s := applyNoErrors(t, r, sameLine)
	if s != sameLine {
		t.Errorf("input should not have been modified: %s", s)
	}

	s = applyNoErrors(t, r, "remove\nbegin\nremove\nend\nremove\nbegin\nremove\nend\nremove\n")
	if s != "remove\nbegin\nend\nremove\nbegin\nend\nremove\n" {
		t.Errorf("some lines should have been removed: %s", s)
	}
}

func TestApply_DeleteCheckCount(t *testing.T) {
	r := Recipe{
		Delete: []DeleteEntry{
			{Search: "removeAlways", CheckCount: 3},
			{Context: Context{Begin: "\\[begin\\]", End: "\\[end\\]"}, Search: "removeContext", CheckCount: 1},
		},
	}
	r.Compile()

	i := `
removeAlways
removeContext

--[begin]--
removeAlways
removeContext
--[end]--

removeAlways
removeContext`
	s := applyNoErrors(t, r, i)

	if s != `
removeContext

--[begin]--
--[end]--

removeContext` {
		t.Errorf("some lines should have been removed: %s", s)
	}

	r.Delete[0].CheckCount = 2
	r.Delete[1].CheckCount = 2
	_, errs := ApplyToInput(r, []byte(i))
	if len(errs) != 2 {
		t.Errorf("unexpected number of errors: %d\n", len(errs))
	}
}

func TestApply_Replace(t *testing.T) {
	r := Recipe{
		Replace: []ReplaceEntry{
			{Search: "search", Replace: "replace"},
		},
	}
	r.Compile()

	s := applyNoErrors(t, r, "line1\nsearch\nline2\ntest search replace\n")
	if s != "line1\nreplace\nline2\ntest replace replace\n" {
		t.Errorf("'search' lines should have been replaced: %s", s)
	}
}

func TestApply_ReplaceSubmatch(t *testing.T) {
	r := Recipe{
		Replace: []ReplaceEntry{
			{Search: "key = (.*), search", Replace: "key = $1, replaced"},
		},
	}
	r.Compile()

	s := applyNoErrors(t, r, "key = value1, search\n")
	if s != "key = value1, replaced\n" {
		t.Errorf("'search' should have been replaced: %s", s)
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

	s := applyNoErrors(t, r, `
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
	s := applyNoErrors(t, r, sameLine)
	if s != sameLine {
		t.Errorf("input should not have been modified: %s", s)
	}

	s = applyNoErrors(t, r, "search\nbegin\nsearch\nend\nsearch\nbegin\nsearch\nend\nsearch\n")
	if s != "search\nbegin\nreplace\nend\nsearch\nbegin\nreplace\nend\nsearch\n" {
		t.Errorf("some lines should have been replaced: %s", s)
	}
}

func TestApply_ReplaceCheckCount(t *testing.T) {
	r := Recipe{
		Replace: []ReplaceEntry{
			{Search: "searchAlways", Replace: "replaceAlways", CheckCount: 3},
			{Context: Context{Begin: "\\[begin\\]", End: "\\[end\\]"}, Search: "searchContext", Replace: "replaceContext", CheckCount: 1},
		},
	}
	r.Compile()

	i := `
searchAlways
searchContext

--[begin]--
searchAlways
searchContext
--[end]--

searchAlways
searchContext`
	s := applyNoErrors(t, r, i)

	if s != `
replaceAlways
searchContext

--[begin]--
replaceAlways
replaceContext
--[end]--

replaceAlways
searchContext` {
		t.Errorf("some lines should have been replaced: %s", s)
	}

	r.Replace[0].CheckCount = 2
	r.Replace[1].CheckCount = 2
	_, errs := ApplyToInput(r, []byte(i))
	if len(errs) != 2 {
		t.Errorf("unexpected number of errors: %d\n", len(errs))
	}
}

func TestApply_ReplaceCheckCountMulti(t *testing.T) {
	r := Recipe{
		Replace: []ReplaceEntry{
			{Search: "searchAlways", Replace: "replaceAlways", CheckCount: 3},
			{Context: Context{Begin: "\\[begin\\]", End: "\\[end\\]"}, Search: "searchContext", Replace: "replaceContext", CheckCount: 1},
		},
	}
	r.Compile()

	i := `
searchAlways
searchContext

--[begin]--
searchAlways
searchContext
--[end]--

searchAlways
searchContext`
	s := applyNoErrors(t, r, i)

	if s != `
replaceAlways
searchContext

--[begin]--
replaceAlways
replaceContext
--[end]--

replaceAlways
searchContext` {
		t.Errorf("some lines should have been replaced: %s", s)
	}

	r.Replace[0].CheckCount = 2
	r.Replace[1].CheckCount = 2
	_, errs := ApplyToInput(r, []byte(i))
	if len(errs) != 2 {
		t.Errorf("unexpected number of errors: %d\n", len(errs))
	}
}

func TestApply_Append(t *testing.T) {
	r := Recipe{
		Append: "append",
	}

	s := applyNoErrors(t, r, "line\n")
	if s != "line\nappend\n" {
		t.Errorf("line should have been appended: %s", s)
	}

	// Apply should add a newline after the input (if there is none).
	s = applyNoErrors(t, r, "line")
	if s != "line\nappend\n" {
		t.Errorf("line should have been appended after newline: %s", s)
	}

	// Apply should also handle empty files.
	s = applyNoErrors(t, r, "")
	if s != "append\n" {
		t.Errorf("line should have been appended: %s", s)
	}

	s = applyNoErrors(t, r, "\n")
	if s != "append\n" {
		t.Errorf("line should have been appended: %s", s)
	}
}

func TestApply_AppendNewLine(t *testing.T) {
	r := Recipe{
		Append: "line1\nline2\n",
	}

	// There should be no additiona newline after the append.
	s := applyNoErrors(t, r, "original\n")
	if s != "original\nline1\nline2\n" {
		t.Errorf("lines should have been appended: %s", s)
	}
}

func TestApply_Empty(t *testing.T) {
	r := Recipe{}

	s := applyNoErrors(t, r, "")
	if s != "" {
		t.Errorf("result should be empty: %s", s)
	}

	s = applyNoErrors(t, r, "\n")
	if s != "\n" {
		t.Errorf("result should be empty: %s", s)
	}
}

func TestApply_Newline(t *testing.T) {
	r := Recipe{}

	i := "line1\nline2\r"
	s := applyNoErrors(t, r, i)
	if s != i {
		t.Errorf("newlines should have been copied: %s", s)
	}

	i = "line1\r\nline2\n\r"
	s = applyNoErrors(t, r, i)
	if s != i {
		t.Errorf("newlines should have been copied: %s", s)
	}
}
