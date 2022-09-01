package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/seanhagen/cipherator/cipher"
	"github.com/spf13/cobra"
)

const encodeShortHelpText = "Encode some text using the named cipher."
const encodeLongHelpText = `Use one of the built-in ciphers to encode some text.

Use the 'list-ciphers' command to see the list of built-in ciphers`

// getEncodeCommand ...
func getEncodeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "encode <cipher>",
		Short: encodeShortHelpText,
		Long:  encodeLongHelpText,
		RunE: func(cmd *cobra.Command, args []string) error {
			useCipher := args[0]
			toEncode := strings.Join(args[1:], " ")

			c, err := cipher.ParseEncoderType(useCipher)
			if err != nil {
				return fmt.Errorf("unable to parse cipher name: %w", err)
			}

			reader := strings.NewReader(toEncode)
			output := bytes.NewBuffer(nil)

			enc, err := cipher.New(c, output)
			if err != nil {
				return fmt.Errorf("unable to fetch cipher: %w", err)
			}

			err = enc.Encode(reader)
			if err != nil {
				return fmt.Errorf("unable to encode using the '%v' cipher: %w", c.String(), err)
			}

			_, err = cmd.OutOrStdout().Write(output.Bytes())
			return err
		},
	}
}

func setupEncodeCommand(root, enc *cobra.Command) {
	root.AddCommand(enc)

	// flags for encoding go here
}
