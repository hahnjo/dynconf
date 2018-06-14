// SPDX-License-Identifier:	GPL-3.0-or-later

package dynconf

import (
	"os"
)

type Config struct {
	base string
	orig *string
	new  *string
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func getOrig(filename string) string {
	return filename + ".orig"
}

func (c *Config) findOrig() {
	orig := getOrig(c.base)
	if exists(orig) {
		c.orig = &orig
	}
}

var NewSuffixes = []string{
	".pacnew", // pacman (Arch Linux)
}

func (c *Config) findNew() {
	for _, suffix := range NewSuffixes {
		newFilename := c.base + suffix
		if exists(newFilename) {
			c.new = &newFilename
			return
		}
	}
}

func NewConfig(filename string) *Config {
	c := Config{base: filename}
	c.findOrig()
	c.findNew()
	return &c
}

func (c *Config) GetInput() string {
	if c.new != nil {
		return *c.new
	} else if c.orig != nil {
		return *c.orig
	}

	// Fall back to using the base file.
	return c.base
}
