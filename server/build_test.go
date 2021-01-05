package server

import "testing"

func TestGoBuild(t *testing.T) {
	err := goBuild("/Users/alex/Projects/AmphionEngine/amphion-tools", "/Users/alex/Projects/AmphionEngine/amphion-tools/build", "test")
	if err != nil {
		t.Fail()
	}
}
