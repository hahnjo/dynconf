// SPDX-License-Identifier:	GPL-3.0-or-later

// +build linux

package dynconf

import (
	"fmt"
	"os"
	"syscall"
	"path"
)

// Write data and synchronize file.
func writeData(file *os.File, data []byte) error {
	_, err := file.Write(data)
	if err != nil {
		return err
	}

	err = file.Sync()
	if err != nil {
		return err
	}

	// All operations succeeded.
	return nil
}

// Apply file mode and file owners.
func applyStat(file *os.File, stat os.FileInfo) error {
	err := file.Chmod(stat.Mode())
	if err != nil {
		return err
	}

	sys, ok := stat.Sys().(*syscall.Stat_t)
	if !ok {
		return fmt.Errorf("can not get owner of %s\n", stat.Name())
	}
	uid, gid := int(sys.Uid), int(sys.Gid)
	err = file.Chown(uid, gid)
	if err != nil {
		return err
	}

	return nil
}

func writeFile(filename string, data []byte, stat os.FileInfo) error {
	var err error

	dir, file := path.Split(filename)
	if dir == "" {
		// Assume filename is relative and use the current directory.
		dir, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	// Create temporary file to write data and atomically move to destination.
	tmpFilename := path.Join(dir, ".dynconf." + file)
	f, err := os.Create(tmpFilename)
	if err != nil {
		return err
	}
	defer os.Remove(tmpFilename)

	// Write data and apply file mode as well as file owners.
	err = writeData(f, data)
	if err != nil {
		return err
	}
	err = applyStat(f, stat)
	if err != nil {
		return err
	}

	// Close file after all operations are done.
	err = f.Close()
	if err != nil {
		return err
	}

	// Move file, this is atomic and will replace the current data.
	err = os.Rename(tmpFilename, filename)
	if err != nil {
		return err
	}

	return nil
}
