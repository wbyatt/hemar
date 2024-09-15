package cmd

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"

	"github.com/wbyatt/hemar/db"
)

var Images = &cobra.Command{
	Use:   "images [SUBCOMMAND]",
	Short: "Manage images",
	Long:  "View and manage images in the hemar registry",
	Run: func(_ *cobra.Command, args []string) {
		listImages()
	},
}

var ImageDescribe = &cobra.Command{
	Use:   "describe",
	Short: "Show details of an image",
	Long:  "Dump the manifest of an image to console",
	Run: func(_ *cobra.Command, args []string) {
		describeImage(args[0])
	},
}

func init() {
	Images.AddCommand(ImageDescribe)
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

func describeImage(image string) {
	duckdb := db.DB()
	defer duckdb.Close()

	var imageRecord db.Image
	// split image into repository and tag
	repository, tag, tagExists := strings.Cut(image, ":")
	imageRecord.Repository = repository
	if tagExists {
		imageRecord.Tag = tag
	} else {
		imageRecord.Tag = "latest"
	}

	exists, err := imageRecord.ExistsByRepositoryAndTag(context.Background(), duckdb)
	if err != nil {
		log.Fatalf("Failed to check if image exists: %v", err)
	}

	if !exists {
		log.Fatalf("Image %s:%s does not exist", repository, tag)
	}

	err = imageRecord.HydrateByRepositoryAndTag(context.Background(), duckdb)
	if err != nil {
		log.Fatalf("Failed to hydrate image: %v", err)
	}

	// print entire imageRecord
	fmt.Printf("%+v\n", imageRecord)

	// fmt.Printf("%s", imageRecord.Manifest)
}
