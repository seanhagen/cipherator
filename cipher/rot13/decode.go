package rot13

import (
	"bytes"
	"io"
	"strings"
)

// Decode ...
func Decode(in string) (string, error) {
	buf := bytes.NewBuffer(nil)
	err := DecodeTo(in, buf)
	return buf.String(), err
}

// DecodeTo ...
func DecodeTo(in string, wr io.Writer) error {
	rt := Encoder{wr}
	return rt.DecodeString(in)
}

// DecodeFromString  ...
func (e *Encoder) DecodeString(in string) error {
	read := strings.NewReader(in)
	return e.Decode(read)
}

// Decode  ...
func (e *Encoder) Decode(r io.Reader) error {
	return e.Encode(r)
}
