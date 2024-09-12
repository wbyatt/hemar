package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wbyatt/hemar/image"
)

var Pull = &cobra.Command{
	Use:   "pull IMAGE",
	Short: "Pull an image from DockerHub",
	Long:  "Fetches a container image from DockerHub and unpacks it locally",
	Run: func(cmd *cobra.Command, args []string) {
		image := image.NewImage(args[0], "latest")
		image.Pull()
	},
}
