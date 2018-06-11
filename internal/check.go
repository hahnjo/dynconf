// SPDX-License-Identifier:	GPL-3.0-or-later

package internal

import (
	"fmt"
	"os"

	"github.com/hahnjo/dynconf/pkg"
)

func Check(args []string) {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Command 'check' requires a config file")
		os.Exit(1)
	}

	file := args[0]
	var c dynconf.Config
	err := c.Read(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading configuration '%s': %s\n", file, err)
		os.Exit(1)
	}

	err = c.Validate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Configuration '%s' invalid: %s\n", file, err)
		os.Exit(1)
	}

	fmt.Printf("Configuration '%s' is valid.\n", file)

	warns := c.Warn()
	if len(warns) != 0 {
		fmt.Println()
		for _, w := range warns {
			fmt.Printf("warning: %s\n", w)
		}
	}

	os.Exit(0)
}
