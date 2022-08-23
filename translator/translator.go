package translator

import (
	"strings"
)

const (
	vowels          = "aeiou"
	consonantAppend = "ay"
	vowelAppend     = "way"
)

var (
	conAppend = strings.Split(consonantAppend, "")
	vowAppend = strings.Split(vowelAppend, "")
)

// PigLatin ...
type PigLatin struct{}

// NewPigLatin ...
func NewPigLatin() (*PigLatin, error) {
	return &PigLatin{}, nil
}

// Translate ...
func (pl PigLatin) Translate(input string) (string, error) {
	var output, currentWord []string
	data := strings.Split(input, "")

	for _, ch := range data {
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
		currentWord = []string{}
	}

	currentWord = pl.doTranslation(currentWord)
	output = append(output, currentWord...)

	return strings.Join(output, ""), nil
}

// doTranslation ...
func (pl PigLatin) doTranslation(in []string) []string {
	if len(in) == 0 {
		return in
	}

	if pl.isUpper(in[0]) {
		in[0] = strings.ToLower(in[0])
		in[1] = strings.ToUpper(in[1])
	}

	var toAppend []string

	if pl.isVowel(in[0]) {
		toAppend = vowAppend
	} else {
		toAppend = conAppend
		in = append(in[1:], in[0])
	}

	in = append(in, toAppend...)

	return in
}

// isLetter ...
func (pl PigLatin) isLetter(in string) bool {
	if in >= "a" && in <= "z" || in >= "A" && in <= "Z" {
		return true
	}
	return false
}

// isVowel  ...
func (pl PigLatin) isVowel(in string) bool {
	if in == "" {
		return false
	}
	return strings.Contains(vowels, in)
}

// isUpper ...
func (pl PigLatin) isUpper(in string) bool {
	if !pl.isLetter(in) {
		return false
	}

	return in == strings.ToUpper(in)
}

// isSpace  ...
func (pl PigLatin) isSpace(in string) bool {
	return in == " "
}
