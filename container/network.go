package container

import (
	"fmt"
	"path/filepath"

	"github.com/wbyatt/hemar/network"
)

func (container *Container) SetupNetwork(bridge string) (func() error, error) {
	nsMountTarget := filepath.Join(containerPath, container.Digest, "netns")
	vethName := fmt.Sprintf("veth%.7s", container.Digest)
	peerName := fmt.Sprintf("P%s", vethName)

	if err := network.SetupVirtualEthernet(vethName, peerName); err != nil {
		return nil, err
	}

	if err := network.SetLinkMaster(vethName, bridge); err != nil {
		return nil, err
	}

	if err := network.SetLinkUp(vethName); err != nil {
		return nil, err
	}

	unmount, err := network.MountNewNetworkNamespace(nsMountTarget)
	if err != nil {
		return unmount, err
	}

	if err := network.SetLinkNsByFile(nsMountTarget, peerName); err != nil {
		return unmount, err
	}

	unset, err := network.SetNetNSByFile(nsMountTarget)
	if err != nil {
		return unmount, err
	}
	defer unset()

	containerEthName := "eth0"
	containerEthIPAddr := "10.69.69.169/16"

	fmt.Println("renaming link: ", peerName, "to", containerEthName)
	if err := network.RenameLink(peerName, containerEthName); err != nil {
		return unmount, err
	}

	if err := network.SetLinkAddress(containerEthName, containerEthIPAddr); err != nil {
		return unmount, err
	}

	if err := network.SetLinkUp(containerEthName); err != nil {
		return unmount, err
	}

	if err := network.AddGateway(containerEthName, "10.69.69.69"); err != nil {
		return unmount, err
	}

	if err := network.SetLinkUp("lo"); err != nil {
		return unmount, err
	}

	return unmount, nil
}
