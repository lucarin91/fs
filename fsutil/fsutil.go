// Package fsutil builds on the underlying fs to add useful functions such as WriteFile or Walk
package fsutil

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/ncw/directio"

	"github.com/codeclysm/fs"
)

type Util interface {
	fs.FS
	ReadDir(dirname string) ([]os.FileInfo, error)
	ReadFile(filename string) ([]byte, error)
	WriteFile(filename string, data []byte, perm os.FileMode) error
	Walk(root string, walkFn filepath.WalkFunc) error
}

type Base struct {
	fs.FS
}

func NewBase(fs fs.FS) Base {
	return Base{FS: fs}
}

func (b Base) ReadDir(dirname string) ([]os.FileInfo, error) {
	f, err := b.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })
	return list, nil
}

func (b Base) ReadFile(filename string) ([]byte, error) {
	f, err := b.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	// It's a good but not certain bet that FileInfo will tell us exactly how much to
	// read, so let's try it but be prepared for the answer to be wrong.
	var n int64 = bytes.MinRead

	if fi, err := f.Stat(); err == nil {
		// As initial capacity for readAll, use Size + a little extra in case Size
		// is zero, and to avoid another allocation after Read has filled the
		// buffer. The readAll call will read into its allocated internal buffer
		// cheaply. If the size was wrong, we'll either waste some space off the end
		// or reallocate as needed, but in the overwhelmingly common case we'll get
		// it just right.
		if size := fi.Size() + bytes.MinRead; size > n {
			n = size
		}
	}
	return readAll(f, n)
}
func (b Base) WriteFile(filename string, data []byte, perm os.FileMode) error {
	f, err := b.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

func (b Base) Walk(root string, walkFn filepath.WalkFunc) error {
	info, err := b.FS.Lstat(root)
	if err != nil {
		err = walkFn(root, nil, err)
	} else {
		err = b.walk(root, info, walkFn)
	}
	if err == filepath.SkipDir {
		return nil
	}
	return err
}

// readAll reads from r until an error or EOF and returns the data it read
// from the internal buffer allocated with a specified capacity.

func (b Base) walk(path string, info os.FileInfo, walkFn filepath.WalkFunc) error {
	if !info.IsDir() {
		return walkFn(path, info, nil)
	}

	names, err := b.readDirNames(path)
	err1 := walkFn(path, info, err)
	// If err != nil, walk can't walk into this directory.
	// err1 != nil means walkFn want walk to skip this directory or stop walking.
	// Therefore, if one of err and err1 isn't nil, walk will return.
	if err != nil || err1 != nil {
		// The caller's behavior is controlled by the return value, which is decided
		// by walkFn. walkFn may ignore err and return nil.
		// If walkFn returns SkipDir, it will be handled by the caller.
		// So walk should return whatever walkFn returns.
		return err1
	}

	for _, name := range names {
		filename := filepath.Join(path, name)
		fileInfo, err := b.FS.Lstat(filename)
		if err != nil {
			if err := walkFn(filename, fileInfo, err); err != nil && err != filepath.SkipDir {
				return err
			}
		} else {
			err = b.walk(filename, fileInfo, walkFn)
			if err != nil {
				if !fileInfo.IsDir() || err != filepath.SkipDir {
					return err
				}
			}
		}
	}
	return nil
}

func (b Base) readDirNames(dirname string) ([]string, error) {
	f, err := b.FS.Open(dirname)
	if err != nil {
		return nil, err
	}
	names, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Strings(names)
	return names, nil
}

func readAll(r io.Reader, capacity int64) ([]byte, error) {
	var buf bytes.Buffer
	buf.Grow(int(capacity))

	// read in block to ensure memory alignment
	block := directio.AlignedBlock(directio.BlockSize)
	for {
		n, err := io.ReadFull(r, block)
		if err != nil {
			if err == io.EOF {
				break
			}
			if err == io.ErrUnexpectedEOF {
				buf.Write(block[:n])
				break
			}

			return nil, err
		}
		buf.Write(block)
	}

	return buf.Bytes(), nil
}
