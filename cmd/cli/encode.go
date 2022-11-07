package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/seanhagen/cipherator/cipher"
	"github.com/spf13/cobra"
)

const encodeShortHelpText = "Encode some text using the named cipher."
const encodeLongHelpText = `Use one of the built-in ciphers to encode some text.

Use the 'list-ciphers' command to see the list of built-in ciphers`

// getEncodeCommand ...
func getEncodeCommand() *cobra.Command {
	var inputFile string
	var outputFile string

	enc := &cobra.Command{
		Use:   "encode <cipher>",
		Short: encodeShortHelpText,
		Long:  encodeLongHelpText,
		RunE: func(cmd *cobra.Command, args []string) error {
			useCipher := args[0]
			c, err := cipher.ParseEncoderType(useCipher)
			if err != nil {
				return fmt.Errorf("unable to parse cipher name: %w", err)
			}

			var reader io.Reader
			if inputFile != "" {
				file, err := os.Open(inputFile)
				if err != nil {
					return fmt.Errorf("unable to %w", err)
				}
				reader = file
				defer file.Close()

			} else {
				toEncode := strings.Join(args[1:], " ")
				reader = strings.NewReader(toEncode)
			}

			var output io.Writer = cmd.OutOrStdout()
			if outputFile != "" {
				file, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
				if err != nil {
					return fmt.Errorf("unable to %w", err)
				}
				output = file
				defer file.Close()
			}

			enc, err := cipher.New(c, output)
			if err != nil {
				return fmt.Errorf("unable to fetch cipher: %w", err)
			}

			err = enc.Encode(reader)
			if err != nil {
				return fmt.Errorf("unable to encode using the '%v' cipher: %w", c.String(), err)
			}

			return nil
		},
	}

	// flags go here
	enc.Flags().StringVarP(&inputFile, "file", "f", "", "file to read from instead of STDIN")
	enc.Flags().StringVarP(&outputFile, "output", "o", "", "file to write output to instead of STDOUT")

	return enc
}

func setupEncodeCommand(root, enc *cobra.Command) {
	root.AddCommand(enc)
}
