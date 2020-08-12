package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObfuscate(t *testing.T) {
	smallInput := "1234"
	largeInput := "1234567890"
	assert.Equal(t, "***", Obfuscate(smallInput, 3))
	assert.Equal(t, "123***890", Obfuscate(largeInput, 3))
}

func TestSplitMethodAndPackage(t *testing.T) {
	realGRPCPath := "/package.Service/Method"
	packageAndService, methodName := SplitMethodAndPackage(realGRPCPath)
	assert.Equal(t, "/package.Service", packageAndService)
	assert.Equal(t, "Method", methodName)
	badPath := "badPath"
	packageAndService, methodName = SplitMethodAndPackage(badPath)
	assert.Empty(t, packageAndService)
	assert.Empty(t, methodName)
	unknowns := "/"
	packageAndService, methodName = SplitMethodAndPackage(unknowns)
	assert.Equal(t, "unknown", packageAndService)
	assert.Equal(t, "unknown", methodName)
}
