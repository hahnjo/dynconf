// SPDX-License-Identifier:	GPL-3.0-or-later

package dynconf

import (
	"bytes"
	"io/ioutil"
)

func Apply(c Config) ([]byte, error) {
	input, err := ioutil.ReadFile(c.File)
	if err != nil {
		return nil, err
	}

	return ApplyToInput(c, input), nil
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

func applyAppend(c Config, modified []byte) []byte {
	if len(c.Append) == 0 {
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

	modified = append(modified, c.Append...)
	modified = append(modified, '\n')

	return modified
}

func ApplyToInput(c Config, input []byte) []byte {
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
		for _, d := range c.Delete {
			if d.SearchRegexp.Match(line) {
				goto next
			}
		}

		// Check if line matches a pattern that shall be replaced.
		for _, r := range c.Replace {
			line = r.SearchRegexp.ReplaceAll(line, r.ReplaceBytes)
		}
		modified = append(modified, line...)
		if to < inLen && input[to] != '\x00' {
			// Copy original newline characters.
			modified = append(modified, input[to:next]...)
		}

	next:
		idx = next
	}

	modified = applyAppend(c, modified)

	return modified
}
