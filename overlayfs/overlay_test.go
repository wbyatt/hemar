package overlayfs

import (
	"fmt"
	"testing"
)

func TestBuildOverlayConfig(t *testing.T) {
	layers := []string{
		"sha256:1234567890abcdef",
		"sha256:0987654321fedcba",
		"sha256:fedcba0987654321",
	}

	root := "/var/lib/hemar"

	overlay, err := BuildOverlayConfig(layers, root)
	if err != nil {
		t.Fatalf("Failed to build overlay config: %v", err)
	}

	expectedUpperDir := fmt.Sprintf("%s/%s", root, layers[2])
	expectedLowerDir := fmt.Sprintf("%s:%s", fmt.Sprintf("%s/%s", root, layers[1]), fmt.Sprintf("%s/%s", root, layers[0]))
	expectedWorkDir := fmt.Sprintf("%s/work", root)

	if overlay.UpperDir != expectedUpperDir {
		t.Errorf("Expected upper dir %s, got %s", expectedUpperDir, overlay.UpperDir)
	}
	if overlay.LowerDir != expectedLowerDir {
		t.Errorf("Expected lower dir %s, got %s", expectedLowerDir, overlay.LowerDir)
	}
	if overlay.WorkDir != expectedWorkDir {
		t.Errorf("Expected work dir %s, got %s", expectedWorkDir, overlay.WorkDir)
	}
}
