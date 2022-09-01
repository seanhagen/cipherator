package piglatin

import (
	"bytes"
	"io"
	"strings"
	"text/scanner"
	"unicode"
)

const (
	vowels    = "aeiou"
	conAppend = "ay"
	vowAppend = "way"
)

// Encoder ...
type Encoder struct {
	output io.Writer
}

// New  ...
func New(wr io.Writer) (*Encoder, error) {
	return &Encoder{wr}, nil
}

// Encode ...
func Encode(in string) (string, error) {
	out := bytes.NewBuffer(nil)
	err := EncodeTo(in, out)
	return out.String(), err
}

// EncodeTo  ...
func EncodeTo(in string, wr io.Writer) error {
	pl, err := New(wr)
	if err != nil {
		return err
	}
	return pl.EncodeFromString(in)
}

// EncodeFromString  ...
func (spl *Encoder) EncodeFromString(in string) error {
	read := strings.NewReader(in)
	return spl.readInto(read, spl.output)
}

// Encode  ...
func (spl *Encoder) Encode(r io.Reader) error {
	return spl.readInto(r, spl.output)
}

// readInto  ...
func (spl *Encoder) readInto(r io.Reader, w io.Writer) error {
	// set up the scanner
	scan := scanner.Scanner{}
	scan.Init(r)
	scan.Filename = "translation"

	// include spaces and tabs as 'tokens'
	scan.Whitespace ^= 1<<'\t' | 1<<' '

	for {
		ch := scan.Scan()
		if ch == scanner.EOF {
			break
		}

		if err := spl.encode(scan.TokenText(), w); err != nil {
			return err
		}
	}

	return nil
}

// encode ...
func (spl *Encoder) encode(token string, writeTo io.Writer) error {
	if len(token) == 0 {
		return nil
	}

	// if our token is only one character, just write it to the output.
	// if it's a character we don't want to do anything, and if it's not we ALSO
	// don't want to do anything.
	if len(token) == 1 {
		_, err := io.WriteString(writeTo, token)
		return err
	}

	// a builder to hold the encoded string we're building
	var build strings.Builder

	// some handy variables we'll be using
	var swapCase bool
	var vowelCase bool = true
	var holdRune rune

	for i, r := range token {
		// if the first character is uppercase, downcase it and toggle swapCase
		if i == 0 && spl.isUpper(r) {
			swapCase = true
			r = unicode.ToLower(r)
		}
		// if swapCase is true and we're on the second character, upcase it
		if i == 1 && swapCase {
			r = unicode.ToUpper(r)
		}

		// is the first character a vowel?
		if i == 0 && !spl.isVowel(r) {
			// nope!
			vowelCase = false
			// hold onto this rune so we can move it to the end
			holdRune = r
			// and continue with the rest of the string
			continue
		}

		// write the current rune to our string
		if _, err := build.WriteRune(r); err != nil {
			return err
		}
	}

	// we need to append 'way' if the first letter is a vowel
	append := []rune{'w', 'a', 'y'}
	// but if the first letter isn't a vowel, replace the 'w' in 'way' with the first
	// character of our original string
	if !vowelCase {
		append[0] = holdRune
	}

	// then append our suffix to the string
	for _, r := range append {
		if _, err := build.WriteRune(r); err != nil {
			return err
		}
	}

	// then write the string to our output
	_, err := io.WriteString(writeTo, build.String())
	return err
}

// isLetter  ...
func (spl *Encoder) isLetter(l rune) bool {
	return unicode.IsLetter(l)
}

// isVowel  ...
func (spl *Encoder) isVowel(in rune) bool {
	// look into using `unicode.In` to test this: https://pkg.go.dev/unicode#In
	for _, v := range vowels {
		if in == v {
			return true
		}
	}

	return false
}

// isUpper ...
func (spl *Encoder) isUpper(in rune) bool {
	if !spl.isLetter(in) {
		return false
	}

	return unicode.IsUpper(in)
}
