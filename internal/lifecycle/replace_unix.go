//go:build !windows

package lifecycle

import "os"

func replacePath(source, destination string) error {
	return os.Rename(source, destination)
}
