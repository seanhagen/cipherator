package rot13

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

// Encoder ...
type Encoder struct {
	wr io.Writer
}

// New ...
func New(wr io.Writer) (*Encoder, error) {
	return &Encoder{wr: wr}, nil
}

// Encode ...
func Encode(in string) (string, error) {
	buf := bytes.NewBuffer(nil)
	err := EncodeTo(in, buf)
	return buf.String(), err
}

// EncodeTo  ...
func EncodeTo(in string, wr io.Writer) error {
	rt := &Encoder{wr}
	return rt.EncodeFromString(in)
}

// EncodeFromString  ...
func (e *Encoder) EncodeFromString(in string) error {
	read := strings.NewReader(in)
	return e.Encode(read)
}

// Encode ...
func (e *Encoder) Encode(r io.Reader) error {
	// fmt.Printf("reader: %v\n", spew.Sdump(r))

	if rr, ok := r.(io.RuneReader); ok {
		return e.encodeFromRunes(rr)
	}

	if rr, ok := r.(io.ByteReader); ok {
		return e.encodeFromBytes(rr)
	}

	return e.encodeFromReader(r)
}

// encodeFromRunes ...
func (e *Encoder) encodeFromRunes(rr io.RuneReader) error {
	for {
		r, _, err := rr.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("unable to read rune: %w", err)
		}

		b := utf8.AppendRune(nil, rot13(r))
		if _, err = e.wr.Write(b); err != nil {
			return fmt.Errorf("unable to write to output: %w", err)
		}
	}

	return nil
}

// encodeFromBytes ...
func (e *Encoder) encodeFromBytes(rr io.ByteReader) error {
	buf := []byte{}

	for {
		b, err := rr.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("unable to read next byte: %w", err)
		}

		buf = append(buf, b)
		r, s := utf8.DecodeRune(buf)
		if r == utf8.RuneError && s == 0 {
			continue
		}

		out := utf8.AppendRune(nil, rot13(r))
		_, err = e.wr.Write(out)
		if err != nil {
			return fmt.Errorf("unable to write byte to output: %w", err)
		}
		buf = []byte{}
	}
	return nil
}

// encodeFromReader  ...
func (e *Encoder) encodeFromReader(r io.Reader) error {
	buf := make([]byte, 8)

	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("unable to encode from reader: %w", err)
		}

		bits := buf[:n]
		out := []byte{}

		for chmp := 0; chmp < len(bits); {
			r, n := utf8.DecodeRune(bits[chmp:])
			if r == utf8.RuneError {
				return fmt.Errorf("unable to decode rune in string '%s'", bits)
			}

			r = rot13(r)

			out = utf8.AppendRune(out, r)
			chmp += n
		}

		_, err = e.wr.Write(out)
		if err != nil {
			return fmt.Errorf("unable to write bytes to output: %w", err)
		}
	}

	return nil
}

// shamelessly ripped from https://stackoverflow.com/questions/31669266/tour-of-go-exercise-23-rot13reader
func rot13(r rune) rune {
	capital := r >= 'A' && r <= 'Z'
	if !capital && (r < 'a' || r > 'z') {
		return r // Not a letter
	}

	r += 13
	if capital && r > 'Z' || !capital && r > 'z' {
		r -= 26
	}
	return r
}
