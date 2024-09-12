package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
)

func init() {

}

func NewHemarCommand() *cobra.Command {
	return &cobra.Command{
		Use:                   "hemar COMMAND",
		Short:                 "A tool for creating and running containers",
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,

		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			if os.Getuid() != 0 {
				return errors.New("Must be root to run hemar")
			}

			return nil
		},
	}
}
