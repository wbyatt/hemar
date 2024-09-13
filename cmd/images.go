package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/wbyatt/hemar/db"
)

var Images = &cobra.Command{
	Use: "images",
	Run: func(_ *cobra.Command, args []string) {
		listImages()
	},
}

func listImages() {
	duckdb := db.DB()
	defer duckdb.Close()

	images, err := db.ListImages(context.Background(), duckdb)
	if err != nil {
		log.Fatalf("Failed to list images: %v", err)
	}

	for _, image := range images {
		fmt.Printf("%s:%s\n", image.Repository, image.Tag)
	}
}
