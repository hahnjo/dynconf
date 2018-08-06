// SPDX-License-Identifier:	GPL-3.0-or-later

package dynconf

import (
	"os"
)

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
