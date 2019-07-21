// SPDX-License-Identifier:	GPL-3.0-or-later

package dynconf

import (
	"bytes"
	"fmt"
	"io/ioutil"
)

func ApplyToFile(r Recipe, filename string) ([]byte, []byte, []error) {
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, nil, []error{err}
	}

	modified, errs := ApplyToInput(r, input)
	return input, modified, errs
}

func isNewLine(c byte) bool {
	return c == '\r' || c == '\n'
}

func containsOnlyNewLines(b []byte) bool {
	for _, c := range b {
		if !isNewLine(c) {
			return false
		}
	}

	return true
}

func evaluateContext(c Context, line []byte, active bool) bool {
	// Index where the begin match ended. This is to avoid matching the same string for the end.
	beginMatch := 0
	if !active && c.BeginRegexp != nil {
		loc := c.BeginRegexp.FindIndex(line)
		if loc != nil {
			active = true
			// Remember where the match ended.
			beginMatch = loc[1]
		}
	}
	if active && c.EndRegexp != nil && c.EndRegexp.Match(line[beginMatch:]) {
		active = false
	}

	return active
}

func applyReplacement(r ReplaceEntry, line []byte) ([]byte, int) {
	s := r.SearchRegexp
	allIndexes := s.FindAllSubmatchIndex(line, -1)
	if allIndexes == nil {
		return line, 0
	}

	modified := make([]byte, 0)
	pos := 0
	for _, loc := range allIndexes {
		// Append bytes up to match.
		modified = append(modified, line[pos:loc[0]]...)
		modified = s.Expand(modified, r.ReplaceBytes, line, loc)
		pos = loc[1]
	}
	// Append rest of line.
	modified = append(modified, line[pos:]...)

	return modified, len(allIndexes)
}

func applyAppend(r Recipe, modified []byte) []byte {
	if len(r.Append) == 0 {
		// Do nothing.
		return modified
	}

	if len(modified) > 0 {
		if !isNewLine(modified[len(modified)-1]) {
			modified = append(modified, '\n')
		} else if containsOnlyNewLines(modified) {
			// There is no content in modified...
			modified = []byte{}
		}
	}

	modified = append(modified, r.Append...)
	if !isNewLine(r.Append[len(r.Append)-1]) {
		modified = append(modified, '\n')
	}

	return modified
}

func ApplyToInput(r Recipe, input []byte) ([]byte, []error) {
	inLen := len(input)

	deleteActive := []bool(nil)
	replaceActive := []bool(nil)
	if r.hasContext {
		// A rule is active iff there is no begin pattern for a context.
		deleteActive = make([]bool, len(r.Delete))
		replaceActive = make([]bool, len(r.Replace))
		for idx, d := range r.Delete {
			deleteActive[idx] = (d.Context.BeginRegexp == nil)
		}
		for idx, r := range r.Replace {
			replaceActive[idx] = (r.Context.BeginRegexp == nil)
		}
	}

	// Count number of matches for deletes and replacements.
	errs := make([]error, 0)
	deleteCount := []int(nil)
	replaceCount := []int(nil)
	if r.hasCount {
		deleteCount = make([]int, len(r.Delete))
		replaceCount = make([]int, len(r.Replace))
	}

	modified := make([]byte, 0)
	// Loop over all lines and modify input.
	idx := 0
	for idx < inLen {
		// Find the first character that introduces a newline.
		to := bytes.IndexAny(input[idx:], "\x00\r\n")
		if to == -1 {
			to = inLen
		} else {
			// 'to' is relative to the slice, so add idx.
			to += idx
		}
		line := input[idx:to]

		// Determine where the next line starts.
		next := to + 1
		// Handle Windows line breaks. Make sure that the two line breaks are different or empty lines might be skipped.
		if next < inLen && input[to] != input[next] && isNewLine(input[next]) {
			next++
		}

		if r.hasContext {
			// For each delete and replace, check if the context begins or ends.
			for idx, d := range r.Delete {
				deleteActive[idx] = evaluateContext(d.Context, line, deleteActive[idx])
			}
			for idx, sr := range r.Replace {
				replaceActive[idx] = evaluateContext(sr.Context, line, replaceActive[idx])
			}
		}

		// Skip line if it matches a pattern that shall be deleted.
		for idx, d := range r.Delete {
			if (!r.hasContext || deleteActive[idx]) && d.SearchRegexp.Match(line) {
				if r.hasCount {
					deleteCount[idx]++
				}
				goto next
			}
		}

		// Check if line matches a pattern that shall be replaced.
		for idx, sr := range r.Replace {
			if !r.hasContext || replaceActive[idx] {
				var count int
				line, count = applyReplacement(sr, line)
				if r.hasCount {
					replaceCount[idx] += count
				}
			}
		}
		modified = append(modified, line...)
		if to < inLen && input[to] != '\x00' {
			// Copy original newline characters.
			modified = append(modified, input[to:next]...)
		}

	next:
		idx = next
	}

	modified = applyAppend(r, modified)

	for idx, d := range r.Delete {
		if d.CheckCount != 0 && d.CheckCount != deleteCount[idx] {
			errs = append(errs, fmt.Errorf("Delete pattern '%s' applied %d times, expected %d!", d.Search, deleteCount[idx], d.CheckCount))
		}
	}
	for idx, r := range r.Replace {
		if r.CheckCount != 0 && r.CheckCount != replaceCount[idx] {
			errs = append(errs, fmt.Errorf("Replace pattern '%s' applied %d times, expected %d!", r.Search, replaceCount[idx], r.CheckCount))
		}
	}

	return modified, errs
}
