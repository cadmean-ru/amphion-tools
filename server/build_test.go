package server

import (
	"runtime"
	"testing"
)

func TestGoBuild(t *testing.T) {
	err := goBuild("/Users/alex/Projects/AmphionEngine/amphion-tools", "/Users/alex/Projects/AmphionEngine/amphion-tools/build", "test", runtime.GOOS, runtime.GOARCH)
	if err != nil {
		t.Fail()
	}
}
