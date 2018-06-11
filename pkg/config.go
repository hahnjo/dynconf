// SPDX-License-Identifier:	GPL-3.0-or-later

package dynconf

import (
	"fmt"
	"os"
	"path"
	"regexp"

	"gopkg.in/yaml.v2"
)

type DeleteEntry struct {
	Search       string
	SearchRegexp *regexp.Regexp
}

type ReplaceEntry struct {
	Search       string
	SearchRegexp *regexp.Regexp
	Replace      string
}

type Config struct {
	File    string
	Delete  []DeleteEntry
	Replace []ReplaceEntry
	Append  string
}

func (c *Config) Read(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	dec := yaml.NewDecoder(file)
	dec.SetStrict(true)

	err = dec.Decode(c)
	if err != nil {
		return err
	}

	for _, d := range c.Delete {
		d.SearchRegexp, err = regexp.Compile(d.Search)
		if err != nil {
			return err
		}
	}

	for _, r := range c.Replace {
		r.SearchRegexp, err = regexp.Compile(r.Search)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) Validate() error {
	if len(c.File) == 0 {
		return fmt.Errorf("Cannot have empty filename!")
	}

	for _, d := range c.Delete {
		if len(d.Search) == 0 {
			return fmt.Errorf("Delete entry cannot have empty regex!")
		}
	}

	for _, r := range c.Replace {
		if len(r.Search) == 0 {
			return fmt.Errorf("Replace entry cannot have empty regex!")
		}
	}

	return nil
}

func (c *Config) Warn() []error {
	w := make([]error, 0)

	if !path.IsAbs(c.File) {
		w = append(w, fmt.Errorf("File should reference an absolute path!"))
	}

	return w
}
