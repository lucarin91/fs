package fs

import (
	"os"
	"time"

	"github.com/ncw/directio"
)

// OS is a wrapper around os package. Doesn't do anything fancy
type OS struct{}

func (o OS) Abs(name string) (string, error) {
	return name, nil
}

func (o OS) Chmod(name string, mode os.FileMode) error {
	return os.Chmod(name, mode)
}

func (o OS) Chown(name string, uid, gid int) error {
	return os.Chown(name, uid, gid)
}

func (o OS) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return os.Chtimes(name, atime, mtime)
}

func (o OS) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (o OS) Lchown(name string, uid, gid int) error {
	return os.Lchown(name, uid, gid)
}

func (o OS) Link(oldname, newname string) error {
	return os.Link(oldname, newname)
}

func (o OS) Lstat(name string) (os.FileInfo, error) {
	return os.Lstat(name)
}

func (o OS) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

func (o OS) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (o OS) NewFile(fd uintptr, name string) *os.File {
	return os.NewFile(fd, name)
}

func (o OS) Open(name string) (*os.File, error) {
	return directio.OpenFile(name, os.O_RDONLY, 0)
}

func (o OS) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return directio.OpenFile(name, flag, perm)
}

func (o OS) Readlink(name string) (string, error) {
	return os.Readlink(name)
}

func (o OS) Remove(name string) error {
	return os.Remove(name)
}

func (o OS) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (o OS) Rename(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}

func (o OS) Symlink(oldname, newname string) error {
	return os.Symlink(oldname, newname)
}

func (o OS) Truncate(name string, size int64) error {
	return os.Truncate(name, size)
}
