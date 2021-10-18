package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCommandLine(t *testing.T) {
	const address = ":1234"

	args := []string{
		"--address", address,
		"--allow-dotfiles",
		"--version",
		"a/b/c",
		"d",
		"e/f",
	}

	expected := &CommandLine{
		Address:       address,
		AllowDotFiles: true,
		Paths:         []string{"a/b/c", "d", "e/f"},
		Version:       true,
	}

	cl, err := ParseCommandLine("", args)
	assert.NoError(t, err)
	assert.Equal(t, expected, cl)
}
