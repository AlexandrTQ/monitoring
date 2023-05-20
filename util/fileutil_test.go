package util

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFileAsLinesSlice(t *testing.T) {
	result, err := GetFileAsLinesSlice("../resourses/server_list_test.txt")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	assert.True(t, strings.Contains(result[0], "google.com"))
	assert.True(t, strings.Contains(result[1], "https://youtube.com"))
	assert.Equal(t, 3, len(result))
}
