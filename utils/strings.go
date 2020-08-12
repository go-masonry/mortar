package utils

import (
	"fmt"
	"strings"
)

// SplitMethodAndPackage is a helper method to split gRPC `package.service/method`
func SplitMethodAndPackage(fullMethodName string) (packageAndService string, methodName string) {
	if index := strings.LastIndex(fullMethodName, "/"); index >= 0 && index+1 <= len(fullMethodName) {
		if packageAndService = fullMethodName[:index]; len(packageAndService) == 0 {
			packageAndService = "unknown"
		}
		if methodName = fullMethodName[index+1:]; len(methodName) == 0 {
			methodName = "unknown"
		}
	}
	return
}

// Obfuscate is a helper function to obfuscate a string, used mostly to hide passwords, etc
func Obfuscate(input string, edgesLength int) string {
	if len(input) > edgesLength*3 { // obfuscated string (middle part) should be at least edge length
		return fmt.Sprintf("%s%s%s", input[:edgesLength], strings.Repeat("*", edgesLength), input[len(input)-edgesLength:])
	}
	return strings.Repeat("*", edgesLength)
}
