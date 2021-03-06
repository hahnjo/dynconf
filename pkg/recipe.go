// SPDX-License-Identifier:	GPL-3.0-or-later

package dynconf

import (
	"fmt"
	"os"
	"path"
	"regexp"

	"gopkg.in/yaml.v2"
)

type Context struct {
	Begin       string
	BeginRegexp *regexp.Regexp
	End         string
	EndRegexp   *regexp.Regexp
}

type DeleteEntry struct {
	Context      Context
	Search       string
	SearchRegexp *regexp.Regexp
	CheckCount   int `yaml:"checkCount"`
}

type ReplaceEntry struct {
	Context      Context
	Search       string
	SearchRegexp *regexp.Regexp
	Replace      string
	ReplaceBytes []byte
	CheckCount   int `yaml:"checkCount"`
}

type Recipe struct {
	File    string
	Delete  []DeleteEntry
	Replace []ReplaceEntry
	Append  string

	hasContext bool
	hasCount   bool
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

		if d.Context.Begin != "" {
			r.hasContext = true
			r.Delete[idx].Context.BeginRegexp, err = regexp.Compile(d.Context.Begin)
			if err != nil {
				return err
			}
		}
		if d.Context.End != "" {
			r.hasContext = true
			r.Delete[idx].Context.EndRegexp, err = regexp.Compile(d.Context.End)
			if err != nil {
				return err
			}
		}

		if d.CheckCount > 0 {
			r.hasCount = true
		}
	}

	for idx, sr := range r.Replace {
		r.Replace[idx].SearchRegexp, err = regexp.Compile(sr.Search)
		if err != nil {
			return err
		}
		r.Replace[idx].ReplaceBytes = []byte(sr.Replace)

		if sr.Context.Begin != "" {
			r.hasContext = true
			r.Replace[idx].Context.BeginRegexp, err = regexp.Compile(sr.Context.Begin)
			if err != nil {
				return err
			}
		}
		if sr.Context.End != "" {
			r.hasContext = true
			r.Replace[idx].Context.EndRegexp, err = regexp.Compile(sr.Context.End)
			if err != nil {
				return err
			}
		}

		if sr.CheckCount > 0 {
			r.hasCount = true
		}
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
		if d.CheckCount < 0 {
			errs = append(errs, fmt.Errorf("Delete entry cannot have negative count!"))
		}
	}

	for _, rs := range r.Replace {
		if len(rs.Search) == 0 {
			errs = append(errs, fmt.Errorf("Replace entry cannot have empty regex!"))
		}
		if rs.CheckCount < 0 {
			errs = append(errs, fmt.Errorf("Replace entry cannot have negative count!"))
		}
	}

	return errs, warns
}
