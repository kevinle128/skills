//go:build windows

package lifecycle

func processAlive(_ int) bool {
	return true
}
