package overlayfs

import (
	"fmt"
	"slices"
	"strings"
	"syscall"
)

// LowerDir does a lot of the heavy lifting for the overlay filesystem.
// It takes a list of layers and returns a string that can be used to
// mount an overlay filesystem
// The string is a colon separated list of directories
// The order is important, the first directory is the "highest" layer,
// and the last directory is the "lowest" layer
type OverlayConfig struct {
	UpperDir string
	LowerDir string
	WorkDir  string
}

type Layer interface {
	Digest() string
}

func BuildOverlayConfig(layers []string, root string) (OverlayConfig, error) {
	if len(layers) == 0 {
		return OverlayConfig{}, fmt.Errorf("no layers provided")
	}

	reverseLayers := safeReverse(layers)

	// split layer into a head and tail
	uppermostLayer := reverseLayers[0]
	reverseLayers = reverseLayers[1:]

	lowerLayers := []string{}
	for _, layer := range reverseLayers {
		lowerLayers = append(lowerLayers, fmt.Sprintf("%s/%s", root, layer))
	}

	return OverlayConfig{
		UpperDir: fmt.Sprintf("%s/%s", root, uppermostLayer),
		LowerDir: strings.Join(lowerLayers, ":"),
		WorkDir:  fmt.Sprintf("%s/work", root),
	}, nil
}

func MountOverlay(target string, config OverlayConfig) (func() error, error) {
	unmounter := func() error {
		return syscall.Unmount(target, 0)
	}

	err := syscall.Mount("", target, "overlay", 0, config.mountString())
	if err != nil {
		return nil, err
	}

	return unmounter, nil
}

func (o OverlayConfig) mountString() string {
	return fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", o.LowerDir, o.UpperDir, o.WorkDir)
}

func safeReverse[S any](items []S) []S {
	copy := slices.Clone(items)
	slices.Reverse(copy)
	return copy
}
