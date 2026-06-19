package config

import "os"

func readFileFromOS(path string) ([]byte, error) {
	return os.ReadFile(path)
}
