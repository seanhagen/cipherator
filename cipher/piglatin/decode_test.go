package piglatin

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPigLatin_IsZeroWidth(t *testing.T) {
	tests := []struct {
		inputRune rune
		expect    bool
	}{
		{'\u200c', true},
		{'a', false},
		{' ', false},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %v run '%c' expect %v", i, tt.inputRune, tt.expect), func(t *testing.T) {
			enc := &Encoder{}

			got := enc.isZeroWidth(tt.inputRune)
			assert.Equal(t, tt.expect, got, "input rune '%c'", tt.inputRune)
		})
	}
}

func TestPigLatinDecoding(t *testing.T) {
	tests := []struct {
		input, expect string
	}{
		{"ellohay", "hello"},
		{"orldway", "[w]orld"},
		{"orldw\u200Cay", "world"},
		{"andway", "[w]and"},
		{"andw\u200Cay", "wand"},
		{"and\u200Cway", "and"},
	}

	t.Run("Decode(string)", func(t *testing.T) {
		for i, tt := range tests {
			t.Run(fmt.Sprintf("test %v decode '%s' to '%s'", i, tt.input, tt.expect), func(t *testing.T) {
				got, err := Decode(tt.input)
				require.NoError(t, err)
				assert.Equal(t, tt.expect, got)
			})
		}
	})

	t.Run("DecodeTo(string, io.Writer)", func(t *testing.T) {
		for i, tt := range tests {
			t.Run(fmt.Sprintf("test %v decode '%s' to '%s'", i, tt.input, tt.expect), func(t *testing.T) {
				buf := bytes.NewBuffer(nil)
				err := DecodeTo(tt.input, buf)
				require.NoError(t, err)
				assert.Equal(t, tt.expect, buf.String())
			})
		}
	})

	t.Run("(*Decoder).DecodeFromString(string)", func(t *testing.T) {
		for i, tt := range tests {
			t.Run(fmt.Sprintf("test %v decode '%s' to '%s'", i, tt.input, tt.expect), func(t *testing.T) {

				buf := bytes.NewBuffer(nil)
				pl, err := New(buf)
				require.NoError(t, err)
				require.NotNil(t, pl)

				err = pl.DecodeString(tt.input)
				assert.NoError(t, err)
				assert.Equal(t, tt.expect, buf.String())
			})
		}
	})

	t.Run("(*Decoder).Decode(io.Writer)", func(t *testing.T) {
		for i, tt := range tests {
			t.Run(fmt.Sprintf("test %v decode '%s' to '%s'", i, tt.input, tt.expect), func(t *testing.T) {
				buf := bytes.NewBuffer(nil)
				pl, err := New(buf)
				require.NoError(t, err)
				require.NotNil(t, pl)
				read := strings.NewReader(tt.input)

				err = pl.Decode(read)
				assert.NoError(t, err)
				assert.Equal(t, tt.expect, buf.String(), "(*PigLatin).Decode(io.Reader)")
			})
		}
	})
}
