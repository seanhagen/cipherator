package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCmd_Root(t *testing.T) {
	output := bytes.NewBuffer(nil)
	expect := rootHelpText

	rootCmd := getRootCommand()
	rootCmd.SetOutput(output)
	err := rootCmd.Execute()

	assert.NoError(t, err)
	assert.Contains(t, output.String(), expect)
}
