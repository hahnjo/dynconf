// SPDX-License-Identifier:	GPL-3.0-or-later

package dynconf

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func createTempFiles(t *testing.T, prefix string, files []string) (string, []string) {
	dir, err := ioutil.TempDir("", prefix)
	if err != nil {
		t.Errorf("could not create temporary directory: %s\n", err)
	}

	filenames := make([]string, 0)
	for _, f := range files {
		filename := path.Join(dir, f)
		file, err := os.Create(filename)
		if err != nil {
			os.RemoveAll(dir)
			t.Errorf("could not create %s: %s\n", filename, err)
		}
		file.Close()

		filenames = append(filenames, filename)
	}

	return dir, filenames
}

func writeTempFile(t *testing.T, filename string, data []byte) {
	err := ioutil.WriteFile(filename, data, 066)
	if err != nil {
		t.Errorf("could not write to %s: %s\n", filename, err)
	}
}

var origData = []byte("orig")
var oldOrigData = []byte("old-orig")
var modifiedData = []byte("modified")
var newData = []byte("new")

func commit(t *testing.T, c *Config) {
	err := c.Commit(origData, modifiedData)
	if err != nil {
		t.Errorf("could not commit: %s\n", err)
	}
}

func checkContent(t *testing.T, filename string, expected []byte) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Errorf("could not read '%s': %s\n", filename, err)
	} else if !bytes.Equal(data, expected) {
		t.Errorf("%s should contain '%s': %s\n", filename, expected, data)
	}
}

func TestBase(t *testing.T) {
	dir, filenames := createTempFiles(t, "base", []string{
		"test.conf",
	})
	defer os.RemoveAll(dir)
	base := filenames[0]

	c := NewConfig(base)
	i := c.GetInput()
	if i != base {
		t.Errorf("getInput should return base: %s\n", i)
	}

	commit(t, c)
	checkContent(t, base, modifiedData)
	orig := path.Join(dir, "test.conf.orig")
	checkContent(t, orig, origData)
}

func TestOrig(t *testing.T) {
	dir, filenames := createTempFiles(t, "orig", []string{
		"test.conf",
		"test.conf.orig",
	})
	defer os.RemoveAll(dir)
	base := filenames[0]
	orig := filenames[1]
	writeTempFile(t, orig, oldOrigData)

	c := NewConfig(base)
	i := c.GetInput()
	if i != orig {
		t.Errorf("getInput should return .orig: %s\n", i)
	}

	commit(t, c)
	checkContent(t, base, modifiedData)
	checkContent(t, orig, oldOrigData)
}

func testNew(t *testing.T, newSuffix string) {
	dir, filenames := createTempFiles(t, newSuffix, []string{
		"test.conf",
		"test.conf.orig",
		"test.conf." + newSuffix,
	})
	defer os.RemoveAll(dir)
	base := filenames[0]
	orig := filenames[1]
	writeTempFile(t, orig, oldOrigData)
	pacnew := filenames[2]
	writeTempFile(t, pacnew, newData)

	c := NewConfig(base)
	i := c.GetInput()
	if i != pacnew {
		t.Errorf("getInput should return .%s: %s\n", newSuffix, i)
	}

	commit(t, c)
	checkContent(t, base, modifiedData)
	checkContent(t, orig, newData)
	if exists(pacnew) {
		t.Errorf(".%s should have been deleted\n", newSuffix)
	}
}

func TestPacnew(t *testing.T) {
	testNew(t, "pacnew")
}

func TestRpmnew(t *testing.T) {
	testNew(t, "rpmnew")
}

func TestNew_Unrecognized(t *testing.T) {
	dir, filenames := createTempFiles(t, "new_unrecognized", []string{
		"test.conf",
		"test.conf.new",
	})
	defer os.RemoveAll(dir)
	base := filenames[0]

	c := NewConfig(base)
	i := c.GetInput()
	if i != base {
		t.Errorf("getInput should return base: %s\n", i)
	}
}

func testNew_NoOrig(t *testing.T, newSuffix string) {
	dir, filenames := createTempFiles(t, newSuffix + "_noOrig", []string{
		"test.conf",
		"test.conf." + newSuffix,
	})
	defer os.RemoveAll(dir)
	base := filenames[0]
	pacnew := filenames[1]
	writeTempFile(t, pacnew, newData)

	c := NewConfig(base)
	i := c.GetInput()
	if i != pacnew {
		t.Errorf("getInput should return .%s: %s\n", newSuffix, i)
	}

	commit(t, c)
	checkContent(t, base, modifiedData)
	orig := path.Join(dir, "test.conf.orig")
	checkContent(t, orig, newData)
	if exists(pacnew) {
		t.Errorf(".%s should have been deleted\n", newSuffix)
	}
}

func testPacnew_NoOrig(t *testing.T) {
	testNew_NoOrig(t, "pacnew")
}

func testRpmnew_NoOrig(t *testing.T) {
	testNew_NoOrig(t, "rpmnew")
}
