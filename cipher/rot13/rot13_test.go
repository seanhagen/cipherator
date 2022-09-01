package rot13

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRot13_Cipher_Encode(t *testing.T) {
	tests := []struct {
		input, expect string
	}{
		{"hello world", "uryyb jbeyq"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %v Encode(%v) => '%v'", i, tt.input, tt.expect), func(t *testing.T) {
			got, err := Encode(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expect, got)
		})
	}
}

func TestRot13_Cipher_EncodeTo(t *testing.T) {
	tests := []struct {
		input, expect string
	}{
		{"hello world", "uryyb jbeyq"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %v EncodeTo(%v, io.Writer) => '%v'", i, tt.input, tt.expect), func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			err := EncodeTo(tt.input, buf)
			assert.NoError(t, err)
			assert.Equal(t, tt.expect, buf.String())
		})
	}
}

func TestRot13_Cipher_ObjEncodeFromString(t *testing.T) {
	tests := []struct {
		input, expect string
	}{
		{"hello world", "uryyb jbeyq"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %v rot13.EncodeFromString(%v) => %v", i, tt.input, tt.expect), func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			rt, err := New(buf)
			require.NotNil(t, rt)
			require.NoError(t, err)

			err = rt.EncodeFromString(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expect, buf.String())
		})
	}
}

func TestRot13_Cipher_ObjEncode(t *testing.T) {
	tests := []struct {
		input, expect string
	}{
		{"hello world", "uryyb jbeyq"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %v rot13.EncodeFromString(%v) => %v", i, tt.input, tt.expect), func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			rt, err := New(buf)
			require.NotNil(t, rt)
			require.NoError(t, err)

			input := strings.NewReader(tt.input)

			err = rt.Encode(input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expect, buf.String())
		})
	}
}
