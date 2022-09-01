package main

import (
	"github.com/spf13/cobra"
)

func main() {
	root := getRootCommand()

	enc := getEncodeCommand()
	setupEncodeCommand(root, enc)

	cobra.CheckErr(root.Execute())
}
