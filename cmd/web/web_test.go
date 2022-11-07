package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWeb_Setup(t *testing.T) {
	args := []string{"-h"}
	output := bytes.NewBuffer(nil)
	expect := webLongHelpText

	cmd := getWebCommand()
	cmd.SetOutput(output)
	cmd.SetArgs(args)
	err := cmd.Execute()

	assert.NoError(t, err)
	assert.Contains(t, output.String(), expect)
}

func TestWeb_BuildRoutes(t *testing.T) {
	var expect *http.Handler
	router, err := getWebRouter()
	require.NoError(t, err)
	require.NotNil(t, router)

	assert.Implements(t, expect, router)
}

func TestWeb_Routes(t *testing.T) {
	tests := []struct {
		operation string
		cipher    string
		input     string
		expect    string
		error     bool
	}{
		{"encode", "piglatin", "hello world", "elloh\u200cay orldw\u200cay", false},
		{"decode", "piglatin", "elloh\u200cay orldw\u200cay", "hello world", false},
		{"encode", "rot13", "hello world", "uryyb jbeyq", false},
		{"decode", "rot13", "uryyb jbeyq", "hello world", false},
	}

	for i, tt := range tests {
		name := fmt.Sprintf("test %v %v %v input '%v' expect '%v'",
			i, tt.operation, tt.cipher, tt.input, tt.expect)
		t.Run(name, func(t *testing.T) {
			router, err := getWebRouter()
			require.NoError(t, err)
			require.NotNil(t, router)

			body := strings.NewReader(tt.input)
			route := fmt.Sprintf("/%v/%v", tt.operation, tt.cipher)

			req, err := http.NewRequest(http.MethodPost, route, body)
			require.NoError(t, err)
			require.NotNil(t, req)

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Result().StatusCode)
			assert.Equal(t, tt.expect, w.Body.String())
		})
	}
}

func TestWeb_SizeLimit(t *testing.T) {
	tests := []struct {
		operation, cipher string
		inputFile         string
		expectSuccess     bool
	}{
		{"encode", "piglatin", "tiny.txt", true},
		{"encode", "piglatin", "kinglear.txt", true},
		{"encode", "piglatin", "toobig.txt", false},
		{"encode", "rot13", "kinglear.txt", true},
		{"encode", "rot13", "toobig.txt", false},
	}

	for i, tt := range tests {
		var name string
		if tt.expectSuccess {
			name = fmt.Sprintf("test %v SHOULD SUCCEED in sending %v", i, tt.inputFile)
		} else {
			name = fmt.Sprintf("test %v SHOULD FAIL    in sending %v", i, tt.inputFile)
		}

		t.Run(name, func(t *testing.T) {
			router, err := getWebRouter()
			require.NoError(t, err)
			require.NotNil(t, router)

			body, err := os.Open("./testdata/" + tt.inputFile)
			require.NoError(t, err)

			route := fmt.Sprintf("/%v/%v", tt.operation, tt.cipher)

			req, err := http.NewRequest(http.MethodPost, route, body)
			require.NoError(t, err)
			require.NotNil(t, req)

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			// because of the way HTTP works, we can't write output and THEN
			// change the status code; we'd have to buffer the whole output
			// to see if an error occurs, then we could set the status code.
			//
			// basically, the status code will always be "200 OK" if the server
			// gets to the point where the cipher has started writing output back
			// to the client.
			assert.Equal(t, http.StatusOK, w.Result().StatusCode)

			if tt.expectSuccess {
				assert.NotContains(t, w.Body.String(), "error during operation")
			} else {
				assert.Equal(t, http.StatusOK, w.Result().StatusCode)
				assert.Contains(t, w.Body.String(), "error during operation")
			}
		})
	}
}
