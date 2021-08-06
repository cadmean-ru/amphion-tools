package utils

import "runtime"

func NewLineString() string {
	switch runtime.GOOS {
	case "windows":
		return "\r\n"
	default:
		return "\n"
	}
}

func PathSeparator() string {
	switch runtime.GOOS {
	case "windows":
		return "\\"
	default:
		return "/"
	}
}