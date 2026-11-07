//go:build !unix

package gpscan

import "os"

func mapFile(path string) ([]byte, func(), error) {
	data, err := os.ReadFile(path)
	return data, func() {}, err
}
