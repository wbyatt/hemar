package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
)

var imagesPath = "/home/byatt/hemar/.hemar/images"
var containersPath = "/home/byatt/hemar/.hemar/containers"
var layersPath = "/home/byatt/hemar/.hemar/layers"

func init() {
	os.MkdirAll(imagesPath, 0700)
	os.MkdirAll(layersPath, 0700)
}

func NewHemarCommand() *cobra.Command {
	return &cobra.Command{
		Use:                   "hemar COMMAND",
		Short:                 "A tool for creating and running containers",
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,

		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			if os.Getuid() != 0 {
				return errors.New("must be root to run hemar")
			}

			return nil
		},
	}
}
