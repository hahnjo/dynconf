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

func ApplyToInput(c Config, input []byte) []byte {
	inLen := len(input)

	modified := make([]byte, 0)
	// Loop over all lines and modify input.
	idx := 0
	for idx < inLen {
		// Find the first character that introduces a newline.
		to := bytes.IndexAny(input[idx:], "\r\n")
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
		if input[to] == '\r' && input[next] == '\n' {
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
		// Copy original newline characters.
		modified = append(modified, input[to:next]...)

	next:
		idx = next
	}

	if len(c.Append) > 0 {
		modified = append(modified, c.Append...)
		modified = append(modified, '\n')
	}

	return modified
}
