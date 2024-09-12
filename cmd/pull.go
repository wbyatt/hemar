package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/wbyatt/hemar/registry"
)

var Pull = &cobra.Command{
	Use:   "pull IMAGE",
	Short: "Pull an image from DockerHub",
	Long:  "Fetches a container image from DockerHub and unpacks it locally",
	Run: func(cmd *cobra.Command, args []string) {
		registry := registry.NewRegistryApi()
		repository := args[0]

		latestManifest, err := registry.PullManifestsForTag(repository, "latest")
		if err != nil {
			log.Fatalf("Failed to find a manifest: %v", err)
		}

		manifestLayers, err := registry.PullManifest(repository, latestManifest)
		if err != nil {
			log.Fatalf("Failed to pull manifest: %v", err)
		}

		for _, layer := range manifestLayers {
			registry.PullLayer(repository, layer)
		}

		// registry.PullLayer(repository, blobReference)
		// image := args[0]

		// shellCmd := exec.Command("./pull", image)
		// shellCmd.Stdin, shellCmd.Stdout, shellCmd.Stderr = os.Stdin, os.Stdout, os.Stderr

		// if err := shellCmd.Run(); err != nil {
		// 	panic("Could not fork")
		// }
	},
}
