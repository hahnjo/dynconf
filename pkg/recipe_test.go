// SPDX-License-Identifier:	GPL-3.0-or-later

package dynconf

import (
	"io/ioutil"
	"os"
	"testing"
)

func writeRecipe(t *testing.T, recipe string) string {
	file, err := ioutil.TempFile("", "recipe.yml")
	if err != nil {
		t.Errorf("could not create temporary file: %s\n", err)
	}

	_, err = file.WriteString(recipe)
	if err != nil {
		t.Errorf("could not write to temporary file: %s\n", err)
	}

	file.Close()
	if err != nil {
		t.Errorf("could not close temporary file: %s\n", err)
	}

	return file.Name()
}

func TestRead(t *testing.T) {
	filename := writeRecipe(t, `
file: "test.conf"

delete:
  -
    search: "remove"

replace:
  -
    search: "pattern"
    replace: "substitution"

append: "last line"`)
	defer os.Remove(filename)

	var r Recipe
	err := r.Read(filename)
	if err != nil {
		t.Errorf("could not read recipe: %s\n", err)
	}

	if r.File != "test.conf" {
		t.Errorf("file was not read correctly: %s\n", r.File)
	}

	if len(r.Delete) != 1 {
		t.Errorf("wrong number of delete entries: %d\n", len(r.Delete))
	}
	d := r.Delete[0]
	if d.Search != "remove" || d.SearchRegexp.String() != "remove" {
		t.Errorf("delete pattern was not read correctly: %s\n", d.Search)
	} else if d.Context.BeginRegexp != nil || d.Context.EndRegexp != nil {
		t.Errorf("delete should not have context!\n")
	}

	if len(r.Replace) != 1 {
		t.Errorf("wrong number of replace entries: %d\n", len(r.Replace))
	}
	rs := r.Replace[0]
	if rs.Search != "pattern" || rs.SearchRegexp.String() != "pattern" {
		t.Errorf("replace pattern was not read correctly: %s\n", rs.Search)
	} else if rs.Replace != "substitution" {
		t.Errorf("replacement was not read correctly: %s\n", rs.Replace)
	} else if rs.Context.BeginRegexp != nil || rs.Context.EndRegexp != nil {
		t.Errorf("replacement should not have context!\n")
	}

	if r.Append != "last line" {
		t.Errorf("append was not read correctly: %s\n", r.Append)
	}
}

func TestRead_Context(t *testing.T) {
	filename := writeRecipe(t, `
delete:
  -
    context:
      begin: "begin"
      end: "end"
    search: "remove"

replace:
  -
    context:
      begin: "begin"
      end: "end"
    search: "pattern"
    replace: "substitution"

append: "last line"`)
	defer os.Remove(filename)

	var r Recipe
	err := r.Read(filename)
	if err != nil {
		t.Errorf("could not read recipe: %s\n", err)
	}

	if len(r.Delete) != 1 {
		t.Errorf("wrong number of delete entries: %d\n", len(r.Delete))
	}
	d := r.Delete[0]
	if d.Context.Begin != "begin" || d.Context.BeginRegexp.String() != "begin" {
		t.Errorf("delete context begin was not read correctly: %s\n", d.Context.Begin)
	} else if d.Context.End != "end" || d.Context.EndRegexp.String() != "end" {
		t.Errorf("delete context end was not read correctly: %s\n", d.Context.End)
	} else if d.Search != "remove" || d.SearchRegexp.String() != "remove" {
		t.Errorf("delete pattern was not read correctly: %s\n", d.Search)
	}

	if len(r.Replace) != 1 {
		t.Errorf("wrong number of replace entries: %d\n", len(r.Replace))
	}
	rs := r.Replace[0]
	if rs.Context.Begin != "begin" || rs.Context.BeginRegexp.String() != "begin" {
		t.Errorf("replace context begin was not read correctly: %s\n", d.Context.Begin)
	} else if rs.Context.End != "end" || rs.Context.EndRegexp.String() != "end" {
		t.Errorf("replace context end was not read correctly: %s\n", d.Context.End)
	} else if rs.Search != "pattern" || rs.SearchRegexp.String() != "pattern" {
		t.Errorf("replace pattern was not read correctly: %s\n", rs.Search)
	} else if rs.Replace != "substitution" {
		t.Errorf("replacement was not read correctly: %s\n", rs.Replace)
	}

	if r.Append != "last line" {
		t.Errorf("append was not read correctly: %s\n", r.Append)
	}
}

func TestValidateErrs(t *testing.T) {
	filename := writeRecipe(t, `
file: ""

delete:
  -
    search: ""

replace:
  -
    search: ""
    replace: ""`)
	defer os.Remove(filename)

	var r Recipe
	err := r.Read(filename)
	if err != nil {
		t.Errorf("could not read recipe: %s\n", err)
	}

	errs, warns := r.Validate()
	if len(errs) != 3 {
		t.Errorf("unexpected number of errors: %d\n", len(errs))
	} else if len(warns) != 0 {
		t.Errorf("unexpected number of warnings: %d\n", len(warns))
	}
}

func TestValidateWarns(t *testing.T) {
	filename := writeRecipe(t, "file: 'relative.conf'")
	defer os.Remove(filename)

	var r Recipe
	err := r.Read(filename)
	if err != nil {
		t.Errorf("could not read recipe: %s\n", err)
	}

	errs, warns := r.Validate()
	if len(errs) != 0 {
		t.Errorf("unexpected number of errors: %d\n", len(errs))
	} else if len(warns) != 1 {
		t.Errorf("unexpected number of warnings: %d\n", len(warns))
	}
}
