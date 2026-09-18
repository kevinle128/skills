//go:build !windows

package update

import "os"

func replaceExecutable(source, destination string) error {
	return os.Rename(source, destination)
}
