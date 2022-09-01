package main

import (
	"github.com/spf13/cobra"
)

func main() {
	root := getWebCommand()
	cobra.CheckErr(root.Execute())
}
