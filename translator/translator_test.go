package translator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTranslator_Basics(t *testing.T) {
	var pl *PigLatin
	var err error

	pl, err = NewPigLatin()
	assert.NotNil(t, pl)
	assert.NoError(t, err)
}

func TestTranslator_Words(t *testing.T) {
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
			got, err := pl.Translate(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.output, got)
		})
	}
}

func TestTranslator_IsLetter(t *testing.T) {
	tests := []struct {
		input  string
		expect bool
	}{
		{"a", true},
		{"B", true},
		{"", false},
		{" ", false},
		{"!", false},
		{"1", false},
	}

	pl, err := NewPigLatin()
	assert.NotNil(t, pl)
	assert.NoError(t, err)

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v is letter %v", tt.input, tt.expect), func(t *testing.T) {
			got := pl.isLetter(tt.input)
			assert.Equal(t, tt.expect, got, "isLetter(%v)", tt.input)
		})
	}
}

func TestTranslator_IsSpace(t *testing.T) {
	tests := []struct {
		input  string
		expect bool
	}{
		{" ", true},
		{"", false},
		{"a", false},
		{"1", false},
		{"!", false},
	}

	pl, err := NewPigLatin()
	assert.NotNil(t, pl)
	assert.NoError(t, err)

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v is letter %v", tt.input, tt.expect), func(t *testing.T) {
			got := pl.isSpace(tt.input)
			assert.Equal(t, tt.expect, got, "isSpace('%v')", tt.input)
		})
	}
}

func TestTranslator_IsVowel(t *testing.T) {
	tests := []struct {
		input  string
		expect bool
	}{
		{"a", true},
		{"e", true},
		{"i", true},
		{"o", true},
		{"u", true},
		{"y", false},
		{"b", false},
		{"!", false},
		{" ", false},
		{"", false},
	}

	pl, err := NewPigLatin()
	assert.NotNil(t, pl)
	assert.NoError(t, err)

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v is vowel %v", tt.input, tt.expect), func(t *testing.T) {
			got := pl.isVowel(tt.input)
			assert.Equal(t, tt.expect, got, "isVowel('%v')", tt.input)
		})
	}
}

func TestTranslator_IsUpper(t *testing.T) {
	tests := []struct {
		input  string
		expect bool
	}{
		{"A", true},
		{"a", false},
		{"!", false},
		{" ", false},
		{"", false},
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

func TestTranslator_DoTranslation(t *testing.T) {
	tests := []struct {
		input  []string
		expect []string
	}{
		{[]string{"h", "e", "l", "l", "o"}, []string{"e", "l", "l", "o", "h", "a", "y"}},
		{[]string{"H", "e", "l", "l", "o"}, []string{"E", "l", "l", "o", "h", "a", "y"}},
		{[]string{"e", "a", "t"}, []string{"e", "a", "t", "w", "a", "y"}},
	}

	pl, err := NewPigLatin()
	assert.NotNil(t, pl)
	assert.NoError(t, err)

	for _, tt := range tests {
		t.Run(fmt.Sprintf("converting '%v' to '%v'", tt.input, tt.expect), func(t *testing.T) {
			got := pl.doTranslation(tt.input)
			assert.Equal(t, tt.expect, got, "doTranslation(%v)", tt.input)
		})
	}
}
