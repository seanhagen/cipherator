//go:generate go-enum -f=$GOFILE --marshal

package cipher

// EncoderType ...
// ENUM(piglatin, rot13)
type EncoderType int32
