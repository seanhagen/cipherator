package piglatin

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

const (
	lengthOfPerfectVowelSuffix = 3
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
func (e *Handler) DecodeString(in string) error {
	read := strings.NewReader(in)
	return e.Decode(read)
}

// Decode  ...
func (e *Handler) Decode(r io.Reader) error {
	scan := e.getScanner(r)
	return e.scanTokens(scan, e.decodeToken)
}

// decodeString  ...
func (e *Handler) decodeToken(token string) error {
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
func (e *Handler) splitToken(token string) (word, suffix []rune, hasZW bool) {
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
func (e *Handler) decodeBestGuess(word []rune) error {
	if len(word) == 1 {
		return e.writeRunes(word)
	}

	wl := len(word)
	suffix := word[wl-lengthOfPerfectVowelSuffix:]
	word = word[:wl-lengthOfPerfectVowelSuffix]

	if suffix[0] == e.suffix[0] {
		word = append([]rune{'[', suffix[0], ']'}, word...)
	} else {
		word = append([]rune{suffix[0]}, word...)
	}

	return e.writeRunes(word)
}

// decodePerfect  ...
func (e *Handler) decodePerfect(word, suffix []rune) error {
	if len(suffix) == lengthOfPerfectVowelSuffix {
		return e.writeRunes(word)
	}

	l := len(word)
	last := word[l-1]
	word = append([]rune{last}, word[:l-1]...)

	return e.writeRunes(word)
}

// writeRunes  ...
func (e *Handler) writeRunes(input []rune) error {
	b := strings.Builder{}
	for _, r := range input {
		if _, err := b.WriteRune(r); err != nil {
			return fmt.Errorf("unable to write rune to output string: %w", err)
		}
	}

	_, err := io.WriteString(e.output, b.String())
	return err
}
