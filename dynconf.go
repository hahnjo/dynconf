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

	apply	Apply a recipe and commit the result
	check	Validate a recipe
	show	Apply a recipe and output the result

	help	Print this help message
	version	Show version information
`

const version = `
DynConf 1.0.2
Copyright (C) 2018  Jonas Hahnfeld
License GPLv3+: GNU GPL version 3 or later <https://gnu.org/licenses/gpl.html>
This program comes with ABSOLUTELY NO WARRANTY. This is free software, and you
are welcome to redistribute it under certain conditions.
See COPYING or above link for defails.
`

func printUsage() {
	fmt.Println(strings.TrimSpace(usage))
}

func printVersion() {
	fmt.Println(strings.TrimSpace(version))
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}

	switch args[0] {
	case "apply":
		internal.Apply(args[1:])
	case "check":
		internal.Check(args[1:])
	case "show":
		internal.Show(args[1:])

	case "help":
		printUsage()
	case "version":
		printVersion()

	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", args[0])
		printUsage()
		os.Exit(1)
	}
	os.Exit(0)
}
