package cipher

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPigLatin_Basics(t *testing.T) {
	var pl *PigLatin
	var err error

	pl, err = NewPigLatin()
	assert.NotNil(t, pl)
	assert.NoError(t, err)
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
	}

	pl, err := NewPigLatin()
	assert.NotNil(t, pl)
	assert.NoError(t, err)

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v to %v", tt.input, tt.output), func(t *testing.T) {
			got, err := pl.Encode(tt.input)
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

	pl, err := NewPigLatin()
	assert.NotNil(t, pl)
	assert.NoError(t, err)

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%c is letter %v", tt.input, tt.expect), func(t *testing.T) {
			got := pl.isLetter(tt.input)
			assert.Equal(t, tt.expect, got, "isLetter(%c)", tt.input)
		})
	}
}

func TestPigLatin_IsSpace(t *testing.T) {
	tests := []struct {
		input  rune
		expect bool
	}{
		{' ', true},
		{'a', false},
		{'1', false},
		{'!', false},
	}

	pl, err := NewPigLatin()
	assert.NotNil(t, pl)
	assert.NoError(t, err)

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%c is letter %v", tt.input, tt.expect), func(t *testing.T) {
			got := pl.isSpace(tt.input)
			assert.Equal(t, tt.expect, got, "isSpace('%c')", tt.input)
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

	pl, err := NewPigLatin()
	assert.NotNil(t, pl)
	assert.NoError(t, err)

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

	pl, err := NewPigLatin()
	assert.NotNil(t, pl)
	assert.NoError(t, err)

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v is uppercase %v", tt.input, tt.expect), func(t *testing.T) {
			got := pl.isUpper(tt.input)
			assert.Equal(t, tt.expect, got, "isUpper('%v')", tt.input)
		})
	}
}

func TestPigLatin_DoTranslation(t *testing.T) {
	tests := []struct {
		input  []rune
		expect []rune
	}{
		{[]rune{'h', 'e', 'l', 'l', 'o'}, []rune{'e', 'l', 'l', 'o', 'h', 'a', 'y'}},
		{[]rune{'H', 'e', 'l', 'l', 'o'}, []rune{'E', 'l', 'l', 'o', 'h', 'a', 'y'}},
		{[]rune{'e', 'a', 't'}, []rune{'e', 'a', 't', 'w', 'a', 'y'}},
	}

	pl, err := NewPigLatin()
	assert.NotNil(t, pl)
	assert.NoError(t, err)

	for _, tt := range tests {
		t.Run(fmt.Sprintf("converting '%c' to '%c'", tt.input, tt.expect), func(t *testing.T) {
			got := pl.encodeRunes(tt.input)
			assert.Equal(t, tt.expect, got, "doTranslation(%c)", tt.input)
		})
	}
}
