package cipher

import (
	"bytes"
	"io"
	"testing"

	"github.com/seanhagen/cipherator/cipher/piglatin"
	"github.com/stretchr/testify/assert"
)

func TestCipher_GetEncoder(t *testing.T) {
	var enc Encoder
	var err error
	var buf io.Writer

	buf = bytes.NewBuffer(nil)
	enc, err = New(EncoderTypePiglatin, buf)

	assert.NoError(t, err)
	assert.IsType(t, &piglatin.Handler{}, enc)
}
