package registry

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
)

type ManifestLayer struct {
	Digest    string `json:"digest"`
	Size      int    `json:"size"`
	MediaType string `json:"mediaType"`
}

type Manifest struct {
	Layers []ManifestLayer `json:"layers"`
}

type ManifestPlatform struct {
	Architecture string `json:"architecture"`
	OS           string `json:"os"`
}

type ManifestListEntry struct {
	Digest    string           `json:"digest"`
	MediaType string           `json:"mediaType"`
	Size      int              `json:"size"`
	Platform  ManifestPlatform `json:"platform"`
}

type ManifestList struct {
	Manifests []ManifestListEntry `json:"manifests"`
}

type RegistryApi struct {
	baseUrl string
	authUrl string
	token   string
}

func NewRegistryApi() *RegistryApi {
	return &RegistryApi{
		authUrl: "https://auth.docker.io/token",
		baseUrl: "https://registry-1.docker.io/v2",
		token:   "",
	}
}

func (registry *RegistryApi) PullManifestsForTag(repository string, tag string) (string, error) {
	fmt.Println("Pulling manifest for repository:", repository)

	if registry.token == "" {
		registry.authorizeForRepo(repository)
	}

	manifestUrl := fmt.Sprintf("%s/%s/manifests/%s", registry.baseUrl, repository, tag)

	req, err := http.NewRequest("GET", manifestUrl, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", registry.token))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to get manifest: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	var manifestList ManifestList
	json.Unmarshal(body, &manifestList)

	fmt.Println("Found the following manifests:")
	for _, manifest := range manifestList.Manifests {
		if manifest.Platform.Architecture == runtime.GOARCH && manifest.Platform.OS == runtime.GOOS {
			return manifest.Digest, nil
		}
	}

	return "", errors.New("no manifest found for this architecture and OS")

}

func (registry *RegistryApi) PullManifest(repository string, digest string) ([]ManifestLayer, error) {

	manifestUrl := fmt.Sprintf("%s/%s/manifests/%s", registry.baseUrl, repository, digest)

	req, err := http.NewRequest("GET", manifestUrl, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", registry.token))
	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	req.Header.Add("Accept", "application/vnd.oci.image.manifest.v1+json")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to get manifest: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	var manifest Manifest
	json.Unmarshal(body, &manifest)

	return manifest.Layers, nil
}

func (registry *RegistryApi) PullLayer(repository string, layer ManifestLayer) {
	fmt.Println("Pulling layer:", layer.Digest)

	layerUrl := fmt.Sprintf("%s/%s/blobs/%s", registry.baseUrl, repository, layer.Digest)

	req, err := http.NewRequest("GET", layerUrl, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", registry.token))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", layer.MediaType)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to get layer: %v", err)
	}
	defer response.Body.Close()

	layerData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Failed to read layer body: %v", err)
	}

	layerName := fmt.Sprintf("./.hemar/images/%s.tar.gz", layer.Digest[len(layer.Digest)-7:])
	os.WriteFile(layerName, layerData, 0644)

	fmt.Println("Saved layer to", layerName)
}

func (registry *RegistryApi) downloadLayer(repository string, digest string) {
}

func (registry *RegistryApi) authorizeForRepo(repository string) {
	authUrl := fmt.Sprintf("%s?service=registry.docker.io&scope=repository:%s:pull", registry.authUrl, repository)

	response, err := http.Get(authUrl)
	if err != nil {
		log.Fatalf("Failed to get authorization token: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	var bodyMap map[string]interface{}

	json.Unmarshal(body, &bodyMap)

	token := bodyMap["token"]

	registry.token = token.(string)
}
