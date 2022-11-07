package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/seanhagen/cipherator/cipher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestCmd_EncodeReadFromStdin(t *testing.T) {
	encPig := cipher.EncoderTypePiglatin.String()
	encRot13 := cipher.EncoderTypeRot13.String()

	tests := []struct {
		cipher string
		input  []string
		expect string
		error  bool
	}{
		{encPig, []string{"hello world"}, "elloh\u200cay orldw\u200cay", false},
		{encPig, []string{"hello", " ", "world"}, "elloh\u200cay   orldw\u200cay", false},
		{"nope", []string{"hello world"}, "hello world", true},
		{encRot13, []string{"hello"}, "uryyb", false},
		{encRot13, []string{"hello world"}, "uryyb jbeyq", false},
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

func TestCmd_EncodeReadFromFile(t *testing.T) {
	tests := []struct {
		input, output string
		cipher        string
		success       bool
	}{
		{"helloworld-input.txt", "piglatin-output.txt", "piglatin", true},
		{"helloworld-input.txt", "piglatin-output.txt", "rot13", false},
		{"helloworld-input.txt", "rot13-output.txt", "rot13", true},
		{"helloworld-input.txt", "rot13-output.txt", "piglatin", false},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %v file %v with cipher %v expect success %v", i, tt.input, tt.cipher, tt.success), func(t *testing.T) {
			output := bytes.NewBuffer(nil)

			inputFile := "./testdata/" + tt.input
			outputFile := "./testdata/" + tt.output

			expect, err := os.Open(outputFile)
			require.NoError(t, err)

			cmd := getEncodeCommand()
			args := []string{tt.cipher, "-f", inputFile}
			cmd.SetArgs(args)
			cmd.SetOutput(output)
			err = cmd.Execute()
			require.NoError(t, err)

			expectData, err := io.ReadAll(expect)
			require.NoError(t, err)

			gotData, err := io.ReadAll(output)
			require.NoError(t, err)

			if tt.success {
				assert.Equal(t, string(expectData), string(gotData))
			} else {
				assert.NotEqual(t, string(expectData), string(gotData))
			}
		})
	}
}

func TestCmd_EncodeWriteOutputToFile(t *testing.T) {
	tests := []struct {
		cipher string
		input  string
		expect string
	}{
		{"piglatin", "hello", "elloh\u200cay"},
		{"rot13", "hello", "uryyb"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %v encode %v to file", i, tt.input), func(t *testing.T) {
			path, err := os.MkdirTemp("", "")
			require.NoError(t, err)

			outputFileName := path + "/output.txt"
			args := []string{tt.cipher, "-o", outputFileName, tt.input}

			cmd := getEncodeCommand()
			cmd.SetArgs(args)
			err = cmd.Execute()
			require.NoError(t, err)

			gotFile, err := os.Open(outputFileName)
			require.NotNil(t, gotFile)
			require.NoError(t, err)

			got, err := io.ReadAll(gotFile)
			require.NoError(t, err)
			assert.Equal(t, tt.expect, string(got))
		})
	}
}
