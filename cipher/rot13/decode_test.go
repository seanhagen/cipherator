package rot13

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRot13_Decode(t *testing.T) {
	tests := []struct {
		input, expect string
	}{
		{"uryyb jbeyq", "hello world"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %v decode '%s'", i, tt.input), func(t *testing.T) {
			t.Run("Decode(string)", func(t *testing.T) {
				got, err := Decode(tt.input)
				require.NoError(t, err)
				assert.Equal(t, tt.expect, got)
			})

			t.Run("DecodeTo(string, io.Writer)", func(t *testing.T) {
				buf := bytes.NewBuffer(nil)
				err := DecodeTo(tt.input, buf)
				require.NoError(t, err)
				assert.Equal(t, tt.expect, buf.String())
			})

			t.Run("(*Encoder).DecodeString(string)", func(t *testing.T) {
				buf := bytes.NewBuffer(nil)
				rt, err := New(buf)
				require.NoError(t, err)
				err = rt.DecodeString(tt.input)
				require.NoError(t, err)
				assert.Equal(t, tt.expect, buf.String())
			})

			t.Run("(*Encoder).Decode(io.Reader)", func(t *testing.T) {
				buf := bytes.NewBuffer(nil)
				rt, err := New(buf)
				require.NoError(t, err)
				read := strings.NewReader(tt.input)
				err = rt.Decode(read)
				require.NoError(t, err)
				assert.Equal(t, tt.expect, buf.String())
			})
		})
	}
}
