package fs_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/codeclysm/fs/v2"
)

func TestChroot(t *testing.T) {
	tmp, _ := ioutil.TempDir("", "")
	chroot := fs.Chroot{
		FS:   fs.OS{},
		Base: tmp,
	}

	// Write a file on /
	_, err := chroot.Create("/file")
	check(t, err)

	_, err = os.Stat(filepath.Join(tmp, "file"))
	check(t, err)

	// Delete the file
	err = chroot.Remove("/file")
	check(t, err)

	_, err = os.Stat(filepath.Join(tmp, "file"))
	contains(t, err.Error(), "no such file")
}

func check(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
func contains(t *testing.T, str, substr string) {
	if !strings.Contains(str, substr) {
		t.Fatalf("expected `%s` to contain `%s`", str, substr)
	}
}
