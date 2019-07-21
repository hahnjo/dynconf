// SPDX-License-Identifier:	GPL-3.0-or-later

package dynconf

import (
	"bytes"
	"io/ioutil"
)

func ApplyToFile(r Recipe, filename string) ([]byte, []byte, error) {
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, nil, err
	}

	return input, ApplyToInput(r, input), nil
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

func applyReplacement(r ReplaceEntry, line []byte) []byte {
	s := r.SearchRegexp
	allIndexes := s.FindAllSubmatchIndex(line, -1)
	if allIndexes == nil {
		return line
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

	return modified
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

func ApplyToInput(r Recipe, input []byte) []byte {
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
				c := d.Context
				if !deleteActive[idx] && c.BeginRegexp != nil && c.BeginRegexp.Match(line) {
					deleteActive[idx] = true
				}
				if deleteActive[idx] && c.EndRegexp != nil && c.EndRegexp.Match(line) {
					deleteActive[idx] = false
				}
			}
			for idx, sr := range r.Replace {
				c := sr.Context
				if !replaceActive[idx] && c.BeginRegexp != nil && c.BeginRegexp.Match(line) {
					replaceActive[idx] = true
				}
				if replaceActive[idx] && c.EndRegexp != nil && c.EndRegexp.Match(line) {
					replaceActive[idx] = false
				}
			}
		}

		// Skip line if it matches a pattern that shall be deleted.
		for idx, d := range r.Delete {
			if (!r.hasContext || deleteActive[idx]) && d.SearchRegexp.Match(line) {
				goto next
			}
		}

		// Check if line matches a pattern that shall be replaced.
		for idx, sr := range r.Replace {
			if !r.hasContext || replaceActive[idx] {
				line = applyReplacement(sr, line)
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

	return modified
}
