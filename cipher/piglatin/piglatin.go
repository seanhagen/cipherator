package piglatin

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/scanner"
	"unicode"
	"unicode/utf8"
)

const (
	vowels         = "aeiou"
	zeroWidth rune = '\u200c'
)

// Encoder ...
type Encoder struct {
	output        io.Writer
	defaultSuffix []rune
}

// New  ...
func New(wr io.Writer) (*Encoder, error) {
	return &Encoder{output: wr, defaultSuffix: []rune{'w', 'a', 'y'}}, nil
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
func (e *Encoder) EncodeFromString(in string) error {
	read := strings.NewReader(in)
	return e.Encode(read)
}

// Encode  ...
func (e *Encoder) Encode(r io.Reader) error {
	return e.encodeReaderIntoWriter(r, e.output)
}

// encodeReaderIntoWriter  ...
func (e *Encoder) encodeReaderIntoWriter(r io.Reader, w io.Writer) error {
	// set up the scanner
	scan := scanner.Scanner{}
	scan.Init(r)
	scan.Filename = "encoding"

	// include spaces and tabs as 'tokens'
	scan.Whitespace ^= 1<<'\t' | 1<<' '

	return e.scanTokens(scan, e.encodeToken)
	// for {
	// 	ch := scan.Scan()
	// 	if ch == scanner.EOF {
	// 		break
	// 	}

	// 	if err := e.encodeToken(scan.TokenText(), w); err != nil {
	// 		return err
	// 	}
	// }

	// return nil
}

// encodeToken ...
func (e *Encoder) encodeToken(token string) error {
	if len(token) == 0 {
		return nil
	}

	// the suffix to append to the word
	toAppend := make([]rune, len(e.defaultSuffix))
	copy(toAppend, e.defaultSuffix)

	// toAppend := []rune{'w', 'a', 'y'}
	// toAppend := []rune{'h', 'a', 'y'}
	// toAppend := []rune{'y', 'a', 'y'}

	// a builder to hold the encoded string we're building
	var build strings.Builder

	// cache how long our token is
	tokenLen := len(token)

	if tokenLen == 1 {
		r, s := utf8.DecodeRuneInString(token)
		if r == utf8.RuneError || s == 0 {
			return fmt.Errorf("unable to decode rune in string: '%s'", token)
		}

		// it's only a single character, write this rune to our output
		if _, err := build.WriteRune(r); err != nil {
			return fmt.Errorf("unable to write rune: %w", err)
		}

		// if the rune isn't a letter, we're done
		if !e.isLetter(r) {
			_, err := io.WriteString(e.output, build.String())
			return err
		}

		// add our non-printable character
		if _, err := build.WriteRune(0x200C); err != nil {
			return fmt.Errorf("unable to add OxAD to string: %w", err)
		}

		// if it IS a letter, then also append 'way', as the only
		// single-letter words in English we're going to handle are "I"
		// and "a".
		for _, r := range toAppend {
			if _, err := build.WriteRune(r); err != nil {
				return fmt.Errorf("unable to append suffix: %w", err)
			}
		}

		_, err := io.WriteString(e.output, build.String())
		return err
	}

	// some handy variables we'll be using
	var swapCase bool
	var holdRune rune
	vowelCase := true

	for i, r := range token {
		// if the first character is uppercase, downcase it and toggle swapCase
		if i == 0 && e.isUpper(r) {
			swapCase = true
			r = unicode.ToLower(r)
		}
		// if swapCase is true and we're on the second character, upcase it
		if i == 1 && swapCase {
			r = unicode.ToUpper(r)
		}

		// is the first character a vowel?
		if i == 0 && !e.isVowel(r) {
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

	// if the first letter isn't a vowel, replace the 'w' in 'way' with the first
	// character of our original string
	if !vowelCase {
		toAppend = append([]rune{holdRune, zeroWidth}, toAppend[1:]...)
		// toAppend[0] = holdRune
	} else {
		toAppend = append([]rune{zeroWidth}, toAppend...)
	}

	// then append our suffix to the string
	for _, r := range toAppend {
		if _, err := build.WriteRune(r); err != nil {
			return err
		}
	}

	// then write the string to our output
	_, err := io.WriteString(e.output, build.String())
	return err
}

// isLetter  ...
func (e *Encoder) isLetter(l rune) bool {
	return unicode.IsLetter(l)
}

// isVowel  ...
func (e *Encoder) isVowel(in rune) bool {
	// look into using `unicode.In` to test this: https://pkg.go.dev/unicode#In
	for _, v := range vowels {
		if in == v {
			return true
		}
	}

	return false
}

// isUpper ...
func (e *Encoder) isUpper(in rune) bool {
	if !e.isLetter(in) {
		return false
	}

	return unicode.IsUpper(in)
}

// isZeroWidth ...
func (e *Encoder) isZeroWidth(in rune) bool {
	return in == zeroWidth
}
