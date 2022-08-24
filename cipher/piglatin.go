package cipher

import (
	"strings"
	"text/scanner"
	"unicode"
)

const (
	vowels    = "aeiou"
	conAppend = "ay"
	vowAppend = "way"
)

// PigLatin ...
type PigLatin struct {
	s scanner.Scanner
}

// NewPigLatin ...
func NewPigLatin() (*PigLatin, error) {
	return &PigLatin{}, nil
}

// Encode ...
func (pl PigLatin) Encode(input string) (string, error) {
	var output []string

	pl.s.Init(strings.NewReader(input))
	pl.s.Filename = "translation"
	pl.s.Whitespace ^= 1<<'\t' | 1<<' '

	for {
		ch := pl.s.Peek()
		tok := pl.s.Scan()
		if tok == scanner.EOF {
			break
		}

		if pl.isLetter(ch) {
			currentWord := pl.encodeStr(pl.s.TokenText())
			output = append(output, currentWord)
			continue
		} else {
			output = append(output, pl.s.TokenText())
		}
	}

	return strings.Join(output, ""), nil
}

// encodeStr ...
func (pl PigLatin) encodeStr(in string) string {
	var d []rune
	for _, ch := range in {
		d = append(d, ch)
	}
	d = pl.encodeRunes(d)
	out := strings.Builder{}
	for _, ch := range d {
		out.WriteRune(ch)
	}
	return out.String()
}

// encodeRunes ...
func (pl PigLatin) encodeRunes(in []rune) []rune {
	if len(in) == 0 {
		return in
	}

	if pl.isUpper(in[0]) {
		in[0] = unicode.ToLower(in[0])
		in[1] = unicode.ToUpper(in[1])
	}

	var toAppend string

	if pl.isVowel(in[0]) {
		toAppend = vowAppend
	} else {
		toAppend = conAppend
		in = append(in[1:], in[0])
	}

	for _, r := range toAppend {
		in = append(in, r)
	}
	return in
}

// isLetter ...
func (pl PigLatin) isLetter(in rune) bool {
	return unicode.IsLetter(in)
}

// isVowel  ...
func (pl PigLatin) isVowel(in rune) bool {
	for _, v := range vowels {
		if in == v {
			return true
		}
	}

	return false
}

// isUpper ...
func (pl PigLatin) isUpper(in rune) bool {
	if !pl.isLetter(in) {
		return false
	}

	return unicode.IsUpper(in)
}

// isSpace  ...
func (pl PigLatin) isSpace(in rune) bool {
	return unicode.IsSpace(in)
}
