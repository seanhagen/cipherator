package piglatin

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/scanner"
	"unicode"
)

// Decode ...
func Decode(in string) (string, error) {
	buf := bytes.NewBuffer(nil)
	err := DecodeTo(in, buf)
	return buf.String(), err
}

// DecodeTo ...
func DecodeTo(in string, wr io.Writer) error {
	rt, err := New(wr)
	if err != nil {
		return err
	}
	return rt.DecodeString(in)
}

// DecodeFromString  ...
func (e *Encoder) DecodeString(in string) error {
	read := strings.NewReader(in)
	return e.Decode(read)
}

// Decode  ...
func (e *Encoder) Decode(r io.Reader) error {
	scan := scanner.Scanner{}
	scan.Init(r)
	scan.Filename = "decoding"
	scan.Whitespace ^= 1<<'\t' | 1<<' '
	// tell the scanner to treat the zero width rune as part of a token, and not a separator
	scan.IsIdentRune = e.scannerIsIdentRune

	// start decoding tokens
	return e.scanTokens(scan, e.decodeToken)
}

// scanTokens  ...
func (e *Encoder) scanTokens(scan scanner.Scanner, process func(string) error) error {
	for {
		ch := scan.Scan()
		if ch == scanner.EOF {
			break
		}

		if err := process(scan.TokenText()); err != nil {
			return err
		}
	}
	return nil
}

// scannerIsIdentRune ...
func (e *Encoder) scannerIsIdentRune(ch rune, i int) bool {
	if i <= 1 {
		// no numbers in the first two characters, or everything will probably explode
		return (ch == zeroWidth || unicode.IsLetter(ch)) && !unicode.IsDigit(ch)
	}
	return ch == zeroWidth || unicode.IsLetter(ch)
}

// decodeString  ...
func (e *Encoder) decodeToken(token string) error {
	if len(token) == 0 {
		return nil
	}

	word, suffix, hasZW := e.splitToken(token)

	if hasZW {
		return e.decodePerfect(word, suffix)
	}
	return e.decodeBestGuess(word)
}

// splitToken ...
func (e *Encoder) splitToken(token string) (word, suffix []rune, hasZW bool) {
	// split our token into the runes making up the word, and the runes making up the suffix
	for _, r := range token {
		if e.isZeroWidth(r) {
			hasZW = true
			continue
		}

		if !hasZW {
			word = append(word, r)
		} else {
			suffix = append(suffix, r)
		}
	}
	return
}

// decodeBestGuess ...
func (e *Encoder) decodeBestGuess(word []rune) error {
	// wordFormat := "%c%s"
	wl := len(word)
	suffix := word[wl-3:]
	word = word[:wl-3]

	if suffix[0] == e.defaultSuffix[0] {
		word = append([]rune{'[', suffix[0], ']'}, word...)
	} else {
		word = append([]rune{suffix[0]}, word...)
	}

	return e.writeRunes(word)
}

// decodePerfect  ...
func (e *Encoder) decodePerfect(word, suffix []rune) error {
	if len(suffix) == 3 {
		return e.writeRunes(word)
	}

	l := len(word)
	last := word[l-1]
	word = append([]rune{last}, word[:l-1]...)

	return e.writeRunes(word)
}

// wordStartsWithVowel  ...
func (e *Encoder) suffixForWordStartingWithVowel(suffix []rune) bool {
	return len(suffix) == 3
}

// writeRunes  ...
func (e *Encoder) writeRunes(input []rune) error {
	b := strings.Builder{}
	for _, r := range input {
		if _, err := b.WriteRune(r); err != nil {
			return fmt.Errorf("unable to write rune to output string: %w", err)
		}
	}

	_, err := io.WriteString(e.output, b.String())
	return err
}
