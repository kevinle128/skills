//go:build windows

package update

import (
	"errors"
	"os"
)

func replaceExecutable(source, destination string) error {
	old := destination + ".kevinkit-old"
	if err := os.Remove(old); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if err := os.Rename(destination, old); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if err := os.Rename(source, destination); err != nil {
		_ = os.Rename(old, destination)
		return err
	}
	_ = os.Remove(old)
	return nil
}
