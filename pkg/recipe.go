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
	ReplaceBytes []byte
}

type Recipe struct {
	File    string
	Delete  []DeleteEntry
	Replace []ReplaceEntry
	Append  string
}

func (r *Recipe) Read(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	dec := yaml.NewDecoder(file)
	dec.SetStrict(true)

	err = dec.Decode(r)
	if err != nil {
		return err
	}

	return r.Compile()
}

func (r *Recipe) Compile() error {
	var err error

	for idx, d := range r.Delete {
		r.Delete[idx].SearchRegexp, err = regexp.Compile(d.Search)
		if err != nil {
			return err
		}
	}

	for idx, sr := range r.Replace {
		r.Replace[idx].SearchRegexp, err = regexp.Compile(sr.Search)
		if err != nil {
			return err
		}
		r.Replace[idx].ReplaceBytes = []byte(sr.Replace)
	}

	return nil
}

func (r *Recipe) Validate() ([]error, []error) {
	errs := make([]error, 0)
	warns := make([]error, 0)

	if len(r.File) == 0 {
		errs = append(errs, fmt.Errorf("Cannot have empty filename!"))
	} else if !path.IsAbs(r.File) {
		warns = append(warns, fmt.Errorf("File should reference an absolute path!"))
	}

	for _, d := range r.Delete {
		if len(d.Search) == 0 {
			errs = append(errs, fmt.Errorf("Delete entry cannot have empty regex!"))
		}
	}

	for _, rs := range r.Replace {
		if len(rs.Search) == 0 {
			errs = append(errs, fmt.Errorf("Replace entry cannot have empty regex!"))
		}
	}

	return errs, warns
}
