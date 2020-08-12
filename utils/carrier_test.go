package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMDTraceCarrier_Set(t *testing.T) {
	md := MDTraceCarrier{
		"one": []string{"two", "three"},
	}
	md.Set("FOUR", "five")
	assert.Contains(t, md, "four")
}

func TestMDTraceCarrier_ForeachKey(t *testing.T) {
	var md = MDTraceCarrier{}
	md.Set("one", "two")
	md.Set("three", "four")

	counter := 0
	md.ForeachKey(func(k, v string) error {
		counter++
		return nil
	})
	assert.Equal(t, 2, counter)
}

func TestMDTraceCarrier_ForeachKeyWithError(t *testing.T) {
	var md = MDTraceCarrier{}
	md.Set("one", "two")
	err := md.ForeachKey(func(k, v string) error {
		return fmt.Errorf("bad handler")
	})
	assert.EqualError(t, err, "bad handler")
}
