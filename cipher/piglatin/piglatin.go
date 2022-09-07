package piglatin

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	vowels         = "aeiou"
	zeroWidth rune = '\u200c'
)

var (
	defaultSuffix = []rune{'w', 'a', 'y'}
)

// Handler ...
type Handler struct {
	output io.Writer
	suffix []rune
}

// New  ...
func New(wr io.Writer) (*Handler, error) {
	return &Handler{output: wr, suffix: defaultSuffix}, nil
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
func (e *Handler) EncodeFromString(in string) error {
	read := strings.NewReader(in)
	return e.Encode(read)
}

// Encode  ...
func (e *Handler) Encode(r io.Reader) error {
	scan := e.getScanner(r)
	return e.scanTokens(scan, e.encodeToken)
}

// encodeToken ...
func (e *Handler) encodeToken(token string) error {
	if len(token) == 0 {
		return nil
	}

	// if the token is only a single character it's probably either 'I',
	// 'a', 'A', or a special character.
	if len(token) == 1 {
		return e.encodeSingleCharToken(token)
	}

	return e.encodeLongToken(token)
}

// isLetter  ...
func (e *Handler) isLetter(l rune) bool {
	return unicode.IsLetter(l)
}

// isVowel  ...
func (e *Handler) isVowel(in rune) bool {
	// look into using `unicode.In` to test this: https://pkg.go.dev/unicode#In
	for _, v := range vowels {
		if in == v || in == unicode.ToUpper(v) {
			return true
		}
	}

	return false
}

// isUpper ...
func (e *Handler) isUpper(in rune) bool {
	if !e.isLetter(in) {
		return false
	}

	return unicode.IsUpper(in)
}

// isZeroWidth ...
func (e *Handler) isZeroWidth(in rune) bool {
	return in == zeroWidth
}

// encodeSingleCharToken  ...
func (e *Handler) encodeSingleCharToken(token string) error {
	// a builder to hold the encoded string we're building
	var build strings.Builder

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
	for _, r := range e.suffix {
		if _, err := build.WriteRune(r); err != nil {
			return fmt.Errorf("unable to append suffix: %w", err)
		}
	}

	_, err := io.WriteString(e.output, build.String())
	return err
}

// encodeLongToken  ...
func (e *Handler) encodeLongToken(token string) error {
	// some handy variables we'll be using:
	// a builder to hold the encoded string we're building
	var build strings.Builder
	// a variable that we set to true if we need to change the capitalization of a word
	var needToSwapCase bool
	// a variable to signal that the whole WORD is capitalized, not just the first letter
	var wordIsUpCase bool
	// a variable to hold the first character of the word if we need to move it
	var holdRune rune
	// assume the first letter of the word is a vowel for now
	tokenStartsWithVowel := true

	for i, r := range token {
		// if the first character is uppercase, downcase it and toggle swapCase
		if i == 0 && e.isUpper(r) {
			needToSwapCase = true
			// r = unicode.ToLower(r)
		}

		// is the first character a vowel?
		if i == 0 && !e.isVowel(r) {
			// nope!
			tokenStartsWithVowel = false
			// hold onto this rune so we can move it to the end
			holdRune = r
			// and continue with the rest of the string
			continue
		}

		// second letter is also capitalized, so assume the whole word is, don't need
		// to swap the case of the second letter
		if i == 1 && e.isUpper(r) {
			needToSwapCase = false
			wordIsUpCase = true
		}

		// if swapCase is true and we're on the second character, upcase it
		if i == 1 && needToSwapCase {
			r = unicode.ToUpper(r)
		}

		// write the current rune to our string
		if _, err := build.WriteRune(r); err != nil {
			return err
		}
	}

	if needToSwapCase && !wordIsUpCase {
		holdRune = unicode.ToLower(holdRune)
	}

	// the suffix to append to the word
	toAppend := make([]rune, len(e.suffix))
	copy(toAppend, e.suffix)

	// if the first letter isn't a vowel, replace the 'w' in 'way' with the first
	// character of our original string
	if !tokenStartsWithVowel {
		toAppend = append([]rune{holdRune, zeroWidth}, toAppend[1:]...)
		// toAppend[0] = holdRune
	} else {
		toAppend = append([]rune{zeroWidth}, toAppend...)
	}

	if wordIsUpCase {
		for i, r := range toAppend {
			toAppend[i] = unicode.ToUpper(r)
		}
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
