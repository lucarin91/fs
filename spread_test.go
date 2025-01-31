package fs_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/codeclysm/fs/v2"
)

func TestSpreadChroot(t *testing.T) {
	tmp, _ := ioutil.TempDir("", "")
	chroot := fs.Chroot{
		FS:   fs.OS{},
		Base: tmp,
	}
	spread := fs.Spread{
		FS: chroot,
	}

	// Write a file on /
	_, err := spread.Create("/file")
	check(t, err)

	_, err = os.Stat(filepath.Join(tmp, "fi", "le", "file"))
	check(t, err)

	// Delete the file
	err = spread.Remove("/file")
	check(t, err)

	_, err = os.Stat(filepath.Join(tmp, "fi", "le", "file"))
	contains(t, err.Error(), "no such file")
}
