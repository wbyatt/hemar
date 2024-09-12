package cmd

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var Pull = &cobra.Command{
	Use:   "pull IMAGE",
	Short: "Pull an image from DockerHub",
	Long:  "Fetches a container image from DockerHub and unpacks it locally",
	Run: func(cmd *cobra.Command, args []string) {
		image := args[0]

		shellCmd := exec.Command("./pull", image)
		shellCmd.Stdin, shellCmd.Stdout, shellCmd.Stderr = os.Stdin, os.Stdout, os.Stderr

		if err := shellCmd.Run(); err != nil {
			panic("Could not fork")
		}
	},
}
