package main

import (
	_ "github.com/marcboeker/go-duckdb"
	"github.com/wbyatt/hemar/cmd"
)

func main() {
	rootCmd := cmd.NewHemarCommand()
	rootCmd.AddCommand(cmd.Pull)
	rootCmd.AddCommand(cmd.Run)
	rootCmd.AddCommand(cmd.Child)
	rootCmd.AddCommand(cmd.Images)

	rootCmd.Execute()
}
