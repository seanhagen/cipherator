//go:generate go-enum -f=$GOFILE --marshal

package cipher

import (
	"fmt"
	"io"

	"github.com/seanhagen/cipherator/cipher/piglatin"
	"github.com/seanhagen/cipherator/cipher/rot13"
)

// EncoderType ...
// ENUM(piglatin, rot13)
type EncoderType int32

// Encoder
type Encoder interface {
	EncodeFromString(string) error
	Encode(io.Reader) error
}

// New ...
func New(t EncoderType, wr io.Writer) (Encoder, error) {
	switch t {
	case EncoderTypePiglatin:
		return piglatin.New(wr)
	case EncoderTypeRot13:
		return rot13.New(wr)
	}

	return nil, fmt.Errorf("%v is an unknown encoder type", t.String())
}
