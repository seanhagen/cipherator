package piglatin

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPigLatin_Basics(t *testing.T) {
	var pl *Handler
	var err error
	buf := bytes.NewBuffer(nil)

	pl, err = New(buf)
	assert.NotNil(t, pl)
	assert.NoError(t, err)
}

func TestPigLatinEncoding(t *testing.T) {
	tests := []struct {
		input, expect string
	}{
		{"hello", "elloh\u200cay"},
		{"eat", "eat\u200cway"},
		{"by", "yb\u200cay"},
		{"you", "ouy\u200cay"},
		{"at", "at\u200cway"},
		{"to", "ot\u200cay"},
		{"world", "orldw\u200cay"},
		{"apples", "apples\u200cway"},
		{"hello world", "elloh\u200cay orldw\u200cay"},
		{"Hello world", "Elloh\u200cay orldw\u200cay"},
		{"Hello, world!", "Elloh\u200cay, orldw\u200cay!"},
		{"I", "I\u200cway"},
		{"hello\n", "elloh\u200cay\n"},
		{"DUKE OF ALBANY", "UKED\u200cAY OF\u200cWAY ALBANY\u200cWAY"},
	}

	t.Run("Encode(string)", func(t *testing.T) {
		for i, tt := range tests {
			t.Run(fmt.Sprintf("test %v encode '%s' to '%s'", i, tt.input, tt.expect), func(t *testing.T) {
				got, err := Encode(tt.input)
				require.NoError(t, err)
				assert.Equal(t, tt.expect, got)
			})
		}
	})

	t.Run("EncodeTo(string, io.Writer)", func(t *testing.T) {
		for i, tt := range tests {
			t.Run(fmt.Sprintf("test %v encode '%s' to '%s'", i, tt.input, tt.expect), func(t *testing.T) {
				buf := bytes.NewBuffer(nil)
				err := EncodeTo(tt.input, buf)
				require.NoError(t, err)
				assert.Equal(t, tt.expect, buf.String())
			})
		}
	})

	t.Run("(*Encoder).EncodeFromString(string)", func(t *testing.T) {
		for i, tt := range tests {
			t.Run(fmt.Sprintf("test %v encode '%s' to '%s'", i, tt.input, tt.expect), func(t *testing.T) {

				buf := bytes.NewBuffer(nil)
				pl, err := New(buf)
				require.NoError(t, err)
				require.NotNil(t, pl)

				err = pl.EncodeFromString(tt.input)
				assert.NoError(t, err)
				assert.Equal(t, tt.expect, buf.String())
			})
		}
	})

	t.Run("(*Encoder).Encode(io.Reader)", func(t *testing.T) {
		for i, tt := range tests {
			t.Run(fmt.Sprintf("test %v encode '%s' to '%s'", i, tt.input, tt.expect), func(t *testing.T) {
				buf := bytes.NewBuffer(nil)
				pl, err := New(buf)
				require.NoError(t, err)
				require.NotNil(t, pl)
				read := strings.NewReader(tt.input)

				err = pl.Encode(read)
				assert.NoError(t, err)
				assert.Equal(t, tt.expect, buf.String(), "(*PigLatin).Encode(io.Reader)")
			})
		}
	})
}

func TestPigLatin_IsLetter(t *testing.T) {
	tests := []struct {
		input  rune
		expect bool
	}{
		{'a', true},
		{'B', true},
		{' ', false},
		{'!', false},
		{'1', false},
	}

	pl := &Handler{}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%c is letter %v", tt.input, tt.expect), func(t *testing.T) {
			got := pl.isLetter(tt.input)
			assert.Equal(t, tt.expect, got, "isLetter(%c)", tt.input)
		})
	}
}

func TestPigLatin_IsVowel(t *testing.T) {
	tests := []struct {
		input  rune
		expect bool
	}{
		{'a', true},
		{'e', true},
		{'i', true},
		{'o', true},
		{'u', true},
		{'y', false},
		{'b', false},
		{'!', false},
		{' ', false},
	}

	pl := &Handler{}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%c is vowel %v", tt.input, tt.expect), func(t *testing.T) {
			got := pl.isVowel(tt.input)
			assert.Equal(t, tt.expect, got, "isVowel('%c')", tt.input)
		})
	}
}

func TestPigLatin_IsUpper(t *testing.T) {
	tests := []struct {
		input  rune
		expect bool
	}{
		{'A', true},
		{'a', false},
		{'!', false},
		{' ', false},
	}

	pl := &Handler{}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v is uppercase %v", tt.input, tt.expect), func(t *testing.T) {
			got := pl.isUpper(tt.input)
			assert.Equal(t, tt.expect, got, "isUpper('%v')", tt.input)
		})
	}
}
