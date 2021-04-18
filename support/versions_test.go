package support

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsAmphionVersionSupported(t *testing.T) {
	a := assert.New(t)

	ver := "v0.2.0-preview.2"
	supported := IsAmphionVersionSupported(ver)
	a.True(supported)

	ver = "v0.1.13-0.20210418113942-dc411b3c50a0"
	supported = IsAmphionVersionSupported(ver)
	a.True(supported)
}
