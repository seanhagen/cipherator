package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/seanhagen/cipherator/cipher"
	"github.com/stretchr/testify/assert"
)

func TestCmd_Encode_Basics(t *testing.T) {
	output := bytes.NewBuffer(nil)
	expect := encodeLongHelpText

	cmd := getEncodeCommand()
	cmd.SetArgs([]string{"-h"})
	cmd.SetOutput(output)
	err := cmd.Execute()

	assert.NoError(t, err)
	assert.Contains(t, output.String(), expect)
}

func TestCmd_EncodeNoFlags(t *testing.T) {
	encPig := cipher.EncoderTypePiglatin.String()

	tests := []struct {
		cipher string
		input  []string
		expect string
		error  bool
	}{
		{encPig, []string{"hello world"}, "ellohay orldway", false},
		{encPig, []string{"hello", " ", "world"}, "ellohay   orldway", false},
		{"nope", []string{"hello world"}, "hello world", true},
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("test %v cipher %s input %s expect %s error %v", i, tt.cipher, tt.input, tt.expect, tt.error),
			func(t *testing.T) {
				output := bytes.NewBuffer(nil)

				cmd := getEncodeCommand()
				args := append([]string{tt.cipher}, tt.input...)

				cmd.SetArgs(args)
				cmd.SetOutput(output)
				err := cmd.Execute()

				if tt.error {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tt.expect, output.String())
				}
			},
		)
	}
}
