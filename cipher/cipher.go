//go:generate go-enum -f=$GOFILE --marshal

package cipher

import (
	"fmt"
	"io"

	"github.com/seanhagen/cipherator/cipher/piglatin"
	"github.com/seanhagen/cipherator/cipher/rot13"
)

// EncoderType defines what types of encoders the system has available.
// ENUM(piglatin, rot13)
type EncoderType int32

// Encoder defines a type that can encode English text into encoded text
type Encoder interface {
	EncodeFromString(string) error
	Encode(io.Reader) error
}

// Decoder defines a type that decode text that was encoded by a specific cipher,
// so a Pig Latin decoder can decode text into Enligsh text.
type Decoder interface {
	DecodeString(string) error
	Decode(io.Reader) error
}

// Handler is a type that can handle both encoding and decoding for a cipher.
type Handler interface {
	Encoder
	Decoder
}

// New returns a Handler for the specified EncoderType. It writes the output of either
// the encoding or decoding operations to the io.Writer provided.
func New(t EncoderType, wr io.Writer) (Handler, error) {
	switch t {
	case EncoderTypePiglatin:
		return piglatin.New(wr)
	case EncoderTypeRot13:
		return rot13.New(wr)
	}

	return nil, fmt.Errorf("%v is an unknown encoder type", t.String())
}
