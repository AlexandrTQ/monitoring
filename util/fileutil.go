package util

import (
	"os"
	"strings"
)

func GetFileAsLinesSlice(path string) ([]string, error) {
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(fileBytes), "\n")

	return lines, nil
}