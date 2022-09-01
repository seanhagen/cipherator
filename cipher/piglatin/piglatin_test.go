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
	var pl *Encoder
	var err error
	buf := bytes.NewBuffer(nil)

	pl, err = New(buf)
	assert.NotNil(t, pl)
	assert.NoError(t, err)
}

func TestPigLatin_Methods(t *testing.T) {
	input := "hello world"
	expect := "ellohay orldway"

	got, err := Encode(input)
	require.NoError(t, err)
	assert.Equal(t, expect, got, "Encode(string)")

	buf := bytes.NewBuffer(nil)
	err = EncodeTo(input, buf)
	require.NoError(t, err)
	assert.Equal(t, expect, buf.String(), "EncodeTo(string, io.Writer)")

	buf.Reset()
	pl, err := New(buf)
	require.NoError(t, err)
	require.NotNil(t, pl)

	err = pl.EncodeFromString(input)
	assert.NoError(t, err)
	assert.Equal(t, expect, buf.String(), "(*PigLatin).EncodeFromString(string)")

	buf.Reset()
	pl, err = New(buf)
	require.NoError(t, err)
	require.NotNil(t, pl)
	read := strings.NewReader(input)
	err = pl.Encode(read)
	assert.NoError(t, err)
	assert.Equal(t, expect, buf.String(), "(*PigLatin).Encode(io.Reader)")
}

func TestPigLatin_Writer(t *testing.T) {
	input := "hello world"
	expect := "ellohay orldway"

	inRead := strings.NewReader(input)
	output := bytes.NewBuffer(nil)

	pl, err := New(output)
	require.NotNil(t, pl)
	require.NoError(t, err)

	err = pl.Encode(inRead)
	require.NoError(t, err)

	assert.Equal(t, expect, output.String())
}

func TestPigLatin_Words(t *testing.T) {
	tests := []struct {
		input, output string
	}{
		{"hello", "ellohay"},
		{"eat", "eatway"},
		{"world", "orldway"},
		{"apples", "applesway"},
		{"hello world", "ellohay orldway"},
		{"Hello world", "Ellohay orldway"},
		{"Hello, world!", "Ellohay, orldway!"},
		// add test cases for the weird edge cases -- single letter words like "I" or "a", etc
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v to %v", tt.input, tt.output), func(t *testing.T) {
			got, err := Encode(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.output, got)
		})
	}
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

	pl := &Encoder{}
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

	pl := &Encoder{}

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

	pl := &Encoder{}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v is uppercase %v", tt.input, tt.expect), func(t *testing.T) {
			got := pl.isUpper(tt.input)
			assert.Equal(t, tt.expect, got, "isUpper('%v')", tt.input)
		})
	}
}
