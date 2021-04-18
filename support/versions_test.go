package support

import (
	"fmt"
	"testing"
)

func TestIsAmphionVersionSupported(t *testing.T) {
	ver := "0.2.0preview1"
	supported := IsAmphionVersionSupported(ver)
	fmt.Println(supported)
}
