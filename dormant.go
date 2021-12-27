package main

import (
	"os"

	"github.com/iljaSL/dormant/cmd"
	"github.com/pterm/pterm"
)

//
func main() {
	if err := cmd.Execute(); err != nil {
		pterm.Error.Println(err)
		os.Exit(1)
	}
}
