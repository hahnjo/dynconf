// SPDX-License-Identifier:	GPL-3.0-or-later

package internal

import (
	"fmt"
	"os"

	"github.com/hahnjo/dynconf/pkg"
)

func Apply(args []string) {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Command 'apply' requires a recipe")
		os.Exit(1)
	}

	file := args[0]
	var r dynconf.Recipe
	err := r.Read(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading recipe '%s': %s\n", file, err)
		os.Exit(1)
	}

	errs, _ := r.Validate()
	if len(errs) > 0 {
		fmt.Fprintf(os.Stderr, "Recipe '%s' is invalid:\n", file)
		for _, e := range errs {
			fmt.Printf("error: %s\n", e)
		}
		os.Exit(1)
	}

	c := dynconf.NewConfig(r.File)
	input := c.GetInput()

	orig, modified, errs := dynconf.ApplyToFile(r, input)
	if len(errs) != 0 {
		fmt.Fprintf(os.Stderr, "Recipe '%s' could not be applied:\n", file)
		for _, e := range errs {
			fmt.Printf("error: %s\n", e)
		}
		os.Exit(1)
	}

	err = c.Commit(orig, modified)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error commiting %s: %s\n", r.File, err)
		os.Exit(1)
	}

	os.Exit(0)
}
