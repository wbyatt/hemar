package container

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/codeclysm/extract"
	"github.com/wbyatt/hemar/image"
	"github.com/wbyatt/hemar/overlayfs"
)

var containerPath = "/home/byatt/hemar/.hemar/containers"
var imagePath = "/home/byatt/hemar/.hemar/images"
var layerPath = "/home/byatt/hemar/.hemar/layers"

type ContainerConfig struct {
	Hostname string
	Image    *image.Image
}

type Container struct {
	Config *ContainerConfig
	Digest string
}

func NewContainer(config *ContainerConfig) *Container {
	container := &Container{
		Config: config,
		Digest: randomHex(),
	}

	containerPath := path.Join(containerPath, container.Digest)
	os.MkdirAll(containerPath, 0700)

	return container
}

func (container *Container) MountFilesystem() (func() error, error) {
	containerPath := path.Join(containerPath, container.Digest)
	mountPath := path.Join(containerPath, "rootfs")

	layers := []string{}

	for _, layer := range container.Config.Image.Layers {
		tarballPath := path.Join(layerPath, fmt.Sprintf("%s.tar.gz", layer.Digest))
		tarball, err := os.Open(tarballPath)
		if err != nil {
			log.Fatalf("Failed to open tarball: %v", err)
		}
		defer tarball.Close()

		err = extract.Archive(context.Background(), tarball, mountPath, nil)
		if err != nil {
			log.Fatalf("Failed to extract tarball: %v", err)
		}

		layers = append(layers, path.Join(layerPath, layer.Digest))
	}

	overlay, err := overlayfs.BuildOverlayConfig(layers, containerPath)
	if err != nil {
		log.Fatalf("Failed to build overlay config: %v", err)
	}

	unmounter, err := overlayfs.MountOverlay(mountPath, overlay)
	if err != nil {
		log.Fatalf("Failed to mount overlay: %v", err)
	}

	return unmounter, nil
}

// func (container *Container) MountFilesystem() {
// 	// imagePath := path.Join(imagePath, container.Config.Image.Repository, container.Config.Image.Tag)
// 	containerPath := path.Join(containerPath, container.Digest)
// 	mountPath := path.Join(containerPath, "rootfs")

// 	// open image.json at imagePath
// 	imageJson, err := os.Open(path.Join(imagePath, "image.json"))
// 	if err != nil {
// 		log.Fatalf("Failed to open image.json: %v", err)
// 	}
// 	defer imageJson.Close()

// 	// read image.json
// 	imageJsonData, err := io.ReadAll(imageJson)
// 	if err != nil {
// 		log.Fatalf("Failed to read image.json: %v", err)
// 	}

// 	// unmarshal image.json
// 	var imageData image.Image
// 	err = json.Unmarshal(imageJsonData, &imageData)
// 	if err != nil {
// 		log.Fatalf("Failed to unmarshal image.json: %v", err)
// 	}

// 	for _, layer := range imageData.Layers {
// 		// open the tarballs at imagePath/layer.Digest
// 		tarballFileName := fmt.Sprintf("%s.tar.gz", layer.Digest)
// 		tarball, err := os.Open(path.Join(imagePath, tarballFileName))
// 		if err != nil {
// 			log.Fatalf("Failed to open tarball: %v", err)
// 		}
// 		defer tarball.Close()

// 		// extract the tarball to mountPath
// 		ctx := context.Background()
// 		err = extract.Archive(ctx, tarball, mountPath, nil)
// 		if err != nil {
// 			log.Fatalf("Failed to extract tarball: %v", err)
// 		}

// 	}
// }

func (container *Container) Cleanup() {
	containerPath := path.Join(containerPath, container.Digest)
	os.RemoveAll(containerPath)
}

func randomHex() string {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		panic("could not generate hash")
	}

	return hex.EncodeToString(bytes)
}
