package fsutil

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/codeclysm/fs/v2"
	"github.com/gofrs/flock"
)

// Lock wraps ioutil to provide locked read and write
type Lock struct {
	Util
}

func NewLock(fs fs.FS) Lock {
	io := Base{FS: fs}

	return Lock{Util: io}
}

func (l Lock) ReadFile(filename string) ([]byte, error) {
	// Check if file exists
	_, err := l.Lstat(filename)
	if err != nil {
		return nil, err
	}

	// get shared lock
	lockpath, err := l.Abs(filename)
	if err != nil {
		return nil, err
	}

	fileLock, err := rlock(lockpath)
	if err != nil {
		return nil, err
	}

	defer unlock(fileLock)

	// ReadFile
	return l.Util.ReadFile(filename)
}

func (l Lock) WriteFile(filename string, data []byte, perm os.FileMode) error {
	// Get exclusive lock
	lockpath, err := l.Abs(filename)
	if err != nil {
		return err
	}
	fileLock, err := lock(lockpath)
	if err != nil {
		return err
	}
	defer unlock(fileLock)

	// Writefile
	return l.Util.WriteFile(filename, data, perm)
}

func (l Lock) Remove(name string) error {
	// Check if file exists
	_, err := l.Lstat(name)
	if err != nil {
		return err
	}

	// Get exclusive lock
	lockpath, err := l.Abs(name)
	if err != nil {
		return err
	}
	fileLock, err := lock(lockpath + ".lock")
	if err != nil {
		return err
	}
	defer unlock(fileLock)
	defer os.RemoveAll(lockpath + ".lock")

	// Remove
	return l.Util.Remove(name)
}

func (l Lock) RemoveAll(path string) error {
	// Check if file exists
	_, err := l.Lstat(path)
	if err != nil {
		return err
	}
	// Get exclusive lock
	lockpath, err := l.Abs(path)
	if err != nil {
		return err
	}
	fileLock, err := lock(lockpath + ".lock")
	if err != nil {
		return err
	}
	defer unlock(fileLock)
	defer os.RemoveAll(lockpath + ".lock")

	// RemoveAll
	return l.Util.RemoveAll(path)
}

// lock attempts to lock the file for 10 seconds. If 10 seconds pass without success it returns an error
func lock(filename string) (*flock.Flock, error) {
	fileLock := flock.NewFlock(filename)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	locked, err := fileLock.TryLockContext(ctx, 100*time.Millisecond)
	if err != nil {
		return nil, err
	}

	if !locked {
		return nil, errors.New("could not lock")
	}

	return fileLock, nil
}

// rlock attempts to gain a shared lock for the given file for 10 seconds.
// If 10 seconds pass without success it returns an error
func rlock(filename string) (*flock.Flock, error) {
	fileLock := flock.NewFlock(filename)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	locked, err := fileLock.TryRLockContext(ctx, 100*time.Millisecond)
	if err != nil {
		return nil, err
	}

	if !locked {
		return nil, errors.New("could not rlock")
	}

	return fileLock, nil
}

// unlock panics so that the error is not swallowed by the defer
// if it panics it means something is very wrong indeed anyway
func unlock(lock *flock.Flock) {
	err := lock.Unlock()
	if err != nil {
		panic(err)
	}
}
