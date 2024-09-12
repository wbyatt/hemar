package image

import (
	"encoding/json"
	"log"
	"os"
	"path"

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

var imageDirectory = "./.hemar/images"

func NewImage(repository string, tag string) *Image {
	return &Image{
		Repository: repository,
		Tag:        tag,
	}
}

func (image *Image) Pull() {
	registry := registry.NewRegistryApi()
	latestManifest, err := registry.PullManifestsForTag(image.Repository, image.Tag)
	image.Digest = latestManifest

	imagePath := path.Join(imageDirectory, image.Repository, image.Tag)
	os.MkdirAll(imagePath, 0700)

	if err != nil {
		log.Fatalf("Failed to find a manifest: %v", err)
	}

	manifestLayers, err := registry.PullManifest(image.Repository, latestManifest)
	if err != nil {
		log.Fatalf("Failed to pull manifest: %v", err)
	}

	for _, layer := range manifestLayers {
		image.Layers = append(image.Layers, ImageLayer{Digest: layer.Digest})

		registry.PullLayer(image.Repository, layer, imagePath)
	}

	// Write the image to disk as a json file
	imageFile, err := os.Create(path.Join(imagePath, "image.json"))
	if err != nil {
		log.Fatalf("Failed to create image file: %v", err)
	}
	defer imageFile.Close()
	json.NewEncoder(imageFile).Encode(image)
}
