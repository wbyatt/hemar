package image

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/wbyatt/hemar/db"
	"github.com/wbyatt/hemar/registry"
)

type Image struct {
	Repository string
	Tag        string
	Digest     string
	Layers     []ImageLayer
}

type ImageLayer struct {
	Digest string
}

var imagesPath = "./.hemar/images"
var layersPath = "./.hemar/layers"

func NewImage(repository string, tag string) *Image {
	return &Image{
		Repository: repository,
		Tag:        tag,
	}
}

func (image *Image) Pull() {
	var imageRecord db.Image
	duckdb := db.DB()
	defer duckdb.Close()

	registry := registry.NewRegistryApi()
	latestManifest, err := registry.PullManifestsForTag(image.Repository, image.Tag)
	if err != nil {
		log.Fatalf("Failed to find a manifest: %v", err)
	}
	jsonManifest, err := json.Marshal(latestManifest)
	if err != nil {
		log.Fatalf("Failed to marshal manifest: %v", err)
	}
	// print the json manifest
	fmt.Printf("MANIFEST: %s\n", jsonManifest)

	image.Digest = latestManifest.Digest
	imageRecord.Repository = image.Repository
	imageRecord.Tag = image.Tag
	imageRecord.Digest = image.Digest
	imageRecord.CreatedAt = time.Now()
	imageRecord.Manifest = jsonManifest

	imageRecord.Insert(context.Background(), duckdb)

	manifestLayers, err := registry.PullManifest(image.Repository, latestManifest.Digest)
	if err != nil {
		log.Fatalf("Failed to pull manifest: %v", err)
	}

	for _, layer := range manifestLayers {
		var layerRecord db.Layer
		layerRecord.Digest = layer.Digest
		layerRecord.CreatedAt = time.Now()

		exists, err := layerRecord.Exists(context.Background(), duckdb)
		if err != nil {
			log.Fatalf("Failed to check if layer exists: %v", err)
		}

		if exists {
			// check for the existence of the file
			layerPath := path.Join(layersPath, fmt.Sprintf("%s.tar.gz", layer.Digest))
			if _, err := os.Stat(layerPath); !os.IsNotExist(err) {
				fmt.Printf("Layer %s already exists, skipping download\n", layer.Digest[12:])
				continue
			}
		}

		registry.PullLayer(image.Repository, layer, layersPath)

		if err := layerRecord.Insert(context.Background(), duckdb); err != nil {
			log.Fatalf("Failed to insert layer: %v", err)
		}
	}
}
