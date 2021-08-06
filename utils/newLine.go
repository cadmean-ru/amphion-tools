package utils

import (
	"runtime"
	"strings"
)

func NewLineString() string {
	switch runtime.GOOS {
	case "windows":
		return "\r\n"
	default:
		return "\n"
	}
}

func UnixPath(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}
