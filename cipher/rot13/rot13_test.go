package rot13

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRot13(t *testing.T) {
	tests := []struct {
		input, expect string
	}{
		{"hello", "uryyb"},
		{"eat", "rng"},
		{"by", "ol"},
		{"world", "jbeyq"},
		{"apples", "nccyrf"},
		{"hello world", "uryyb jbeyq"},
		{"Hello world", "Uryyb jbeyq"},
		{"Hello, world!", "Uryyb, jbeyq!"},
	}

	t.Run("Encoding", func(t *testing.T) {
		t.Run("Encode(string)", func(t *testing.T) {
			for i, tt := range tests {
				t.Run(fmt.Sprintf("test %v encode '%s' to '%s'", i, tt.input, tt.expect), func(t *testing.T) {
					got, err := Encode(tt.input)
					require.NoError(t, err)
					assert.Equal(t, tt.expect, got)
				})
			}
		})

		t.Run("EncodeTo(string, io.Writer)", func(t *testing.T) {
			for i, tt := range tests {
				t.Run(fmt.Sprintf("test %v encode '%s' to '%s'", i, tt.input, tt.expect), func(t *testing.T) {
					buf := bytes.NewBuffer(nil)
					err := EncodeTo(tt.input, buf)
					require.NoError(t, err)
					assert.Equal(t, tt.expect, buf.String())
				})
			}
		})

		t.Run("(*Encoder).EncodeFromString(string)", func(t *testing.T) {
			for i, tt := range tests {
				t.Run(fmt.Sprintf("test %v encode '%s' to '%s'", i, tt.input, tt.expect), func(t *testing.T) {

					buf := bytes.NewBuffer(nil)
					rt, err := New(buf)
					require.NoError(t, err)
					require.NotNil(t, rt)

					err = rt.EncodeFromString(tt.input)
					assert.NoError(t, err)
					assert.Equal(t, tt.expect, buf.String())
				})
			}
		})

		t.Run("(*Encoder).Encode(io.Writer)", func(t *testing.T) {
			for i, tt := range tests {
				t.Run(fmt.Sprintf("test %v encode '%s' to '%s'", i, tt.input, tt.expect), func(t *testing.T) {
					buf := bytes.NewBuffer(nil)
					rt, err := New(buf)
					require.NoError(t, err)
					require.NotNil(t, rt)
					read := strings.NewReader(tt.input)

					err = rt.Encode(read)
					assert.NoError(t, err)
					assert.Equal(t, tt.expect, buf.String(), "(*PigLatin).Encode(io.Reader)")
				})
			}
		})
	})

	t.Run("Decoding", func(t *testing.T) {
		t.Run("Decode(string)", func(t *testing.T) {
			for i, tt := range tests {
				t.Run(fmt.Sprintf("test %v decode '%s' to '%s'", i, tt.expect, tt.input), func(t *testing.T) {
					got, err := Decode(tt.expect)
					require.NoError(t, err)
					assert.Equal(t, tt.input, got)
				})
			}
		})

		t.Run("DecodeTo(string, io.Writer)", func(t *testing.T) {
			for i, tt := range tests {
				t.Run(fmt.Sprintf("test %v decode '%s' to '%s'", i, tt.expect, tt.input), func(t *testing.T) {
					buf := bytes.NewBuffer(nil)
					err := DecodeTo(tt.expect, buf)
					require.NoError(t, err)
					assert.Equal(t, tt.input, buf.String())
				})
			}
		})

		t.Run("(*Decoder).DecodeFromString(string)", func(t *testing.T) {
			for i, tt := range tests {
				t.Run(fmt.Sprintf("test %v decode '%s' to '%s'", i, tt.expect, tt.input), func(t *testing.T) {

					buf := bytes.NewBuffer(nil)
					rt, err := New(buf)
					require.NoError(t, err)
					require.NotNil(t, rt)

					err = rt.DecodeString(tt.expect)
					assert.NoError(t, err)
					assert.Equal(t, tt.input, buf.String())
				})
			}
		})

		t.Run("(*Decoder).Decode(io.Writer)", func(t *testing.T) {
			for i, tt := range tests {
				t.Run(fmt.Sprintf("test %v decode '%s' to '%s'", i, tt.expect, tt.input), func(t *testing.T) {
					buf := bytes.NewBuffer(nil)
					rt, err := New(buf)
					require.NoError(t, err)
					require.NotNil(t, rt)
					read := strings.NewReader(tt.expect)

					err = rt.Decode(read)
					assert.NoError(t, err)
					assert.Equal(t, tt.input, buf.String(), "(*PigLatin).Decode(io.Reader)")
				})
			}
		})
	})
}

func TestRot13_EncodeFromRunes(t *testing.T) {
	tests := []struct {
		input, expect string
	}{
		{"hello world", "uryyb jbeyq"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %v encodeFromBytes %v", i, tt.input), func(t *testing.T) {
			out := bytes.NewBuffer(nil)
			rt, err := New(out)
			require.NoError(t, err)

			r := strings.NewReader(tt.input)

			err = rt.encodeFromRunes(r)
			require.NoError(t, err)

			assert.Equal(t, tt.expect, out.String())
		})
	}
}

func TestRot13_EncodeFromBytes(t *testing.T) {
	tests := []struct {
		input, expect string
	}{
		{"hello world", "uryyb jbeyq"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %v encodeFromBytes %v", i, tt.input), func(t *testing.T) {
			out := bytes.NewBuffer(nil)
			rt, err := New(out)
			require.NoError(t, err)

			r := strings.NewReader(tt.input)

			err = rt.encodeFromBytes(r)
			require.NoError(t, err)

			assert.Equal(t, tt.expect, out.String())
		})
	}
}

func TestRot13_EncodeFromReader(t *testing.T) {
	tests := []struct {
		input, expect string
	}{
		{"hello world", "uryyb jbeyq"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %v encodeFromBytes %v", i, tt.input), func(t *testing.T) {
			out := bytes.NewBuffer(nil)
			rt, err := New(out)
			require.NoError(t, err)

			r := strings.NewReader(tt.input)

			err = rt.encodeFromReader(r)
			require.NoError(t, err)

			assert.Equal(t, tt.expect, out.String())
		})
	}
}
