package fsutil_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/codeclysm/fs/v2"
	"github.com/codeclysm/fs/v2/fsutil"
)

func TestChrootWalk(t *testing.T) {
	tmp, _ := ioutil.TempDir("", "")
	chroot := fs.Chroot{
		FS:   fs.OS{},
		Base: tmp,
	}

	io := fsutil.NewBase(chroot)

	defer os.RemoveAll(tmp)

	err := io.WriteFile("/file", []byte("hello"), 0644)
	check(t, err)
	err = io.Mkdir("/folder", 0755)
	check(t, err)
	err = io.WriteFile("/folder/file", []byte("hello"), 0644)
	check(t, err)

	files := []string{}

	err = io.Walk("/", func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	check(t, err)

	equal(t, fmt.Sprintf("%s", files), "[/ /file /folder /folder/file]")
}

func TestSpreadWalk(t *testing.T) {
	tmp, _ := ioutil.TempDir("", "")
	chroot := fs.Chroot{
		FS:   fs.OS{},
		Base: tmp,
	}

	spread := fs.Spread{
		FS: chroot,
	}

	io := fsutil.NewBase(spread)

	defer os.RemoveAll(tmp)

	err := io.Mkdir("/folder", 0755)
	check(t, err)
	err = io.WriteFile("/folder/file", []byte("hello"), 0644)
	check(t, err)
	err = io.Mkdir("/folder/folder1", 0755)
	check(t, err)
	err = io.WriteFile("/folder/folder1/file1", []byte("hello"), 0644)
	check(t, err)

	files := []string{}

	err = io.Walk("/folder", func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	check(t, err)

	equal(t, fmt.Sprintf("%s", files), "[/folder /folder/file /folder/folder1 /folder/folder1/file1]")
}

func check(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func equal(t *testing.T, got, expected interface{}) {
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("expected `%v`, got `%v`", expected, got)
	}
}

func contains(t *testing.T, str, substr string) {
	if !strings.Contains(str, substr) {
		t.Fatalf("expected `%s` to contain `%s`", str, substr)
	}
}
