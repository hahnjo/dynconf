// SPDX-License-Identifier:	GPL-3.0-or-later

package internal

import (
	"fmt"
	"os"

	"github.com/hahnjo/dynconf/pkg"
)

func Check(args []string) {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Command 'check' requires a recipe")
		os.Exit(1)
	}

	file := args[0]
	var r dynconf.Recipe
	err := r.Read(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading recipe '%s': %s\n", file, err)
		os.Exit(1)
	}

	err = r.Validate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Recipe '%s' invalid: %s\n", file, err)
		os.Exit(1)
	}

	fmt.Printf("Recipe '%s' is valid.\n", file)

	warns := r.Warn()
	if len(warns) != 0 {
		fmt.Println()
		for _, w := range warns {
			fmt.Printf("warning: %s\n", w)
		}
	}

	os.Exit(0)
}
