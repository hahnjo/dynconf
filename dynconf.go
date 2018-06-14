// SPDX-License-Identifier:	GPL-3.0-or-later

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/hahnjo/dynconf/internal"
)

const usage = `
Usage:	dynconf command [arguments]

The commands are:

	check	Validate a recipe
	show	Apply a recipe and output the result
`

func printUsage() {
	fmt.Println(strings.TrimSpace(usage))
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		printUsage()
		os.Exit(0)
	}

	switch args[0] {
	case "check":
		internal.Check(args[1:])
	case "show":
		internal.Show(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", args[0])
		printUsage()
		os.Exit(1)
	}
}
