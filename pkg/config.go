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

func (c *Config) getOrig() string {
	return c.base + ".orig"
}

func (c *Config) findOrig() {
	orig := c.getOrig()
	if exists(orig) {
		c.orig = &orig
	}
}

var newSuffixes = []string{
	".pacnew", // pacman (Arch Linux)
	".rpmnew", // rpm (RHEL, Fedora)
}

func (c *Config) findNew() {
	for _, suffix := range newSuffixes {
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

func (c *Config) Commit(origData []byte, modified []byte) error {
	// Get FileInfo of the configuration file.
	stat, err := os.Stat(c.base)
	if err != nil {
		return err
	}

	if c.new == nil && c.orig == nil {
		// Copy the unmodified file to allow idempotence.
		origFile := c.getOrig()
		c.orig = &origFile
		err = writeFile(origFile, origData, stat)
		if err != nil {
			return err
		}
	}

	err = writeFile(c.base, modified, stat)
	if err != nil {
		return err
	}

	if c.new != nil {
		// Move the new file to allow idempotence.
		orig := c.getOrig()
		c.orig = &orig
		err = os.Rename(*c.new, orig)
		if err != nil {
			return err
		}

		// The new file doesn't exist anymore.
		c.new = nil
	}

	return nil
}
