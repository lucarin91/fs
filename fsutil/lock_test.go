package fsutil_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/codeclysm/fs/v2"
	"github.com/codeclysm/fs/v2/fsutil"
)

func TestLock(t *testing.T) {
	tmp, _ := ioutil.TempDir("", "")
	chroot := fs.Chroot{
		FS:   fs.OS{},
		Base: tmp,
	}

	spread := fs.Spread{
		FS: chroot,
	}

	io := fsutil.NewLock(spread)

	err := io.MkdirAll("/folder", 0755)
	check(t, err)

	err = io.WriteFile("/folder/file", []byte("hello"), 0644)
	check(t, err)

	_, err = io.ReadFile("/folder/file")
	check(t, err)

	_, err = io.ReadFile("/folder/filemissing")
	contains(t, err.Error(), "no such file")

	err = io.Remove("/folder/filemissing")
	contains(t, err.Error(), "no such file")

	err = io.RemoveAll("/folder/filemissing")
	contains(t, err.Error(), "no such file")

	err = io.MkdirAll("/folder/otherfolder", 0755)
	check(t, err)

	err = io.RemoveAll("/folder/otherfolder")
	check(t, err)

	files := []string{}

	err = io.Walk("/folder", func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	check(t, err)

	equal(t, fmt.Sprintf("%s", files), "[/folder /folder/file]")

}
