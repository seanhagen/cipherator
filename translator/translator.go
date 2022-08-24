package translator

import (
	"strings"
	"unicode"
)

const (
	vowels    = "aeiou"
	conAppend = "ay"
	vowAppend = "way"
)

// PigLatin ...
type PigLatin struct{}

// NewPigLatin ...
func NewPigLatin() (*PigLatin, error) {
	return &PigLatin{}, nil
}

// Translate ...
func (pl PigLatin) Translate(input string) (string, error) {
	var output, currentWord []rune

	for _, ch := range input {
		// if it's a letter, append to currentWord
		if pl.isLetter(ch) {
			currentWord = append(currentWord, ch)
			continue
		}

		// if there isn't anything in current word and we're not on a letter, just append
		// the letter to output and continue on
		if len(currentWord) == 0 && !pl.isLetter(ch) {
			output = append(output, ch)
			continue
		}

		// translate the word
		currentWord = pl.doTranslation(currentWord)

		// append currentWord to output
		output = append(output, currentWord...)

		// add the current character ( ie, not a letter, like spaces or punctuation )
		output = append(output, ch)

		// and reset currentWord to an empty slice
		currentWord = []rune{}
	}

	currentWord = pl.doTranslation(currentWord)
	output = append(output, currentWord...)

	outWr := strings.Builder{}
	for _, r := range output {
		outWr.WriteRune(r)
	}

	return outWr.String(), nil
}

// doTranslation ...
func (pl PigLatin) doTranslation(in []rune) []rune {
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
