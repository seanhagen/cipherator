//go:generate go-enum -f=$GOFILE --marshal

package cipher

import (
	"fmt"
	"io"

	"github.com/seanhagen/cipherator/cipher/piglatin"
	"github.com/seanhagen/cipherator/cipher/rot13"
)

// EncoderType ...
// ENUM(piglatin, reverse)
type EncoderType int32

// Encoder ...
type Encoder interface {
	EncodeFromString(string) error
	Encode(io.Reader) error
}

// Decoder ...
type Decoder interface {
	DecodeString(string) error
	Decode(io.Reader) error
}

// Handler ...
type Handler interface {
	Encoder
	Decoder
}

// New ...
func New(t EncoderType, wr io.Writer) (Handler, error) {
	switch t {
	case EncoderTypePiglatin:
		return piglatin.New(wr)
	case EncoderTypeRot13:
		return rot13.New(wr)
	}

	return nil, fmt.Errorf("%v is an unknown encoder type", t.String())
}
