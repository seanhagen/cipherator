package piglatin

import (
	"bufio"
	"bytes"
	"io"
	"unicode"
	"unicode/utf8"
)

// getScanner ...
func (e *Handler) getScanner(r io.Reader) *bufio.Scanner {
	reader := bufio.NewScanner(r)
	reader.Split(scannerFn)
	return reader
}

// scanTokens  ...
func (e *Handler) scanTokens(scan *bufio.Scanner, process func(string) error) error {
	for scan.Scan() {
		if err := process(scan.Text()); err != nil {
			return err
		}
	}
	err := scan.Err()
	return err
}

func scannerFn(data []byte, atEOF bool) (advance int, token []byte, err error) {
	buf := bytes.NewBuffer(nil)

	// we get in a big chunk in 'data':
	//    []byte{`DUKE OF ALBANY, her husband`}
	//
	// next step should be to start moving through pulling out one rune at a time
	//
	// but we have to deal with words vs non-words
	//   and words like "I'm" or "you're"
	//
	// i think the play is to just break words up by spaces and special characters -- but keep the
	// spaces, so []byte{`DUKE OF ALBANY, her husband`} would become:
	//
	//   - DUKE
	//   - ' '
	//   - OF
	//   - ' '
	//   - ALBANY
	//   - ','
	//   - ' '
	//   - her
	//   - ' '
	//   - husband
	//
	// and contracted words should get returned as a single word, so "can't" should be a full token

	for len(data) > 0 {
		// if len(data) == 0 {
		// 	break
		// }

		r, s := utf8.DecodeRune(data)

		// now do some checking: is this the first character?
		if buf.Len() == 0 {
			// what's special about us being on the first character?
			//
			// if it's a space or other non-letter character, we break after writing it to the buffer
			//  - excpet apostrophes; those can be at the start of a word
			//  - not zero-width though, our "spec" for Pig Latin says the zw char
			//    comes before the suffix, so it can't be the first character in a word
			if !unicode.IsLetter(r) || r == '\'' {
				if _, err := buf.WriteRune(r); err != nil {
					return 0, nil, err
				}
				break
			}
		} else {
			// we're not on the first character

			// so hitting a space or non-letter character other than the
			// zero-width or apostrophy should kick us out of the loop
			if !(unicode.IsLetter(r) || (r == '\'' || r == '\u200c')) {
				break
			}
		}

		// lastly, write the rune to our  buffer
		if _, err := buf.WriteRune(r); err != nil {
			return 0, nil, err
		}

		data = data[s:]
	}

	return buf.Len(), buf.Bytes(), nil
}
