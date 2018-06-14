// SPDX-License-Identifier:	GPL-3.0-or-later

package dynconf

import (
	"bytes"
	"io/ioutil"
)

func ApplyToFile(r Recipe, filename string) ([]byte, error) {
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return ApplyToInput(r, input), nil
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
	modified = append(modified, '\n')

	return modified
}

func ApplyToInput(r Recipe, input []byte) []byte {
	inLen := len(input)

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
		// Handle Windows line breaks...
		if next < inLen && isNewLine(input[next]) {
			next++
		}

		// Skip line if it matches a pattern that shall be deleted.
		for _, d := range r.Delete {
			if d.SearchRegexp.Match(line) {
				goto next
			}
		}

		// Check if line matches a pattern that shall be replaced.
		for _, sr := range r.Replace {
			line = sr.SearchRegexp.ReplaceAll(line, sr.ReplaceBytes)
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
