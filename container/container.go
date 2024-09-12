package container

import (
	"crypto/rand"
	"encoding/hex"
)

var containerPath = "./.hemar/containers"

type ContainerConfig struct {
	Hostname string
}

type Container struct {
	Config *ContainerConfig
	Digest string
}

func NewContainer() *Container {
	return &Container{
		Digest: randomHex(),
	}
}

func randomHex() string {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		panic("could not generate hash")
	}

	return hex.EncodeToString(bytes)
}
