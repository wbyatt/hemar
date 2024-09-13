package network

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"

	"golang.org/x/sys/unix"
)

func SetupBridge(name string) error {
	cmd := exec.Command("ip", "link", "add", name, "type", "bridge")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create bridge %s: %w", name, err)
	}

	cmd = exec.Command("ip", "addr", "add", "10.69.69.69/24", "dev", name)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add ip address to bridge %s: %w", name, err)
	}

	cmd = exec.Command("ip", "link", "set", "dev", name, "up")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to bring up bridge %s: %w", name, err)
	}

	return nil
}

func SetupNAT(bridge string, externalInterface string) error {
	commands := []string{
		"echo 1 > /proc/sys/net/ipv4/ip_forward",
		fmt.Sprintf("iptables -t nat -A POSTROUTING -o %s -j MASQUERADE", externalInterface),
		fmt.Sprintf("iptables -A FORWARD -i %s -o %s -j ACCEPT", bridge, externalInterface),
		fmt.Sprintf("iptables -A FORWARD -i %s -o %s -j ACCEPT", externalInterface, bridge),
		"iptables-save > /etc/iptables/rules.v4",
	}

	for _, command := range commands {
		if err := exec.Command("sh", "-c", command).Run(); err != nil {
			return fmt.Errorf("failed to run command: %s: %w", command, err)
		}
	}

	return nil

}

func TeardownBridge(name string) error {
	cmd := exec.Command("ip", "link", "del", name)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to delete bridge %s: %w", name, err)
	}

	return nil
}

func SetupVirtualEthernet(vethName string, peerName string) error {
	cmd := exec.Command("ip", "link", "add", vethName, "type", "veth", "peer", "name", peerName)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create veth pair %s and %s: %w", vethName, peerName, err)
	}

	return nil
}

func SetLinkMaster(vethName string, bridge string) error {
	cmd := exec.Command("ip", "link", "set", "dev", vethName, "master", bridge)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set bridge master for %s: %w", vethName, err)
	}

	return nil
}

func MountNewNetworkNamespace(target string) (func() error, error) {
	// I think this is just checking that the target exists
	if _, err := os.OpenFile(target, syscall.O_RDONLY|syscall.O_CREAT|syscall.O_EXCL, 0644); err != nil {
		return nil, fmt.Errorf("unable to create target file %s: %w", target, err)
	}

	// store current network namespace
	nsFile, err := os.OpenFile("/proc/self/ns/net", os.O_RDONLY, 0)
	if err != nil {
		return nil, fmt.Errorf("unable to open current network namespace: %w", err)
	}
	defer nsFile.Close()

	// Split current process into new namespace
	if err := syscall.Unshare(syscall.CLONE_NEWNET); err != nil {
		return nil, fmt.Errorf("unable to unshare network namespace: %w", err)
	}

	// Mount the network namespace
	if err := syscall.Mount("/proc/self/ns/net", target, "bind", syscall.MS_BIND, ""); err != nil {
		return nil, fmt.Errorf("unable to bind mount network namespace: %w", err)
	}

	unmount := func() error {
		if err := syscall.Unmount(target, 0); err != nil {
			return fmt.Errorf("unable to unmount network namespace: %w", err)
		}

		if err := os.Remove(target); err != nil {
			return fmt.Errorf("unable to remove network namespace: %w", err)
		}

		return nil
	}

	// reset the previous network namespace
	if err := unix.Setns(int(nsFile.Fd()), unix.CLONE_NEWNET); err != nil {
		return nil, fmt.Errorf("unable to setns to previous network namespace: %w", err)
	}

	return unmount, nil
}

func SetLinkNsByFile(target string, link string) error {

	nsFile, err := os.OpenFile(target, os.O_RDONLY, 0)

	if err != nil {
		return fmt.Errorf("unable to open network namespace: %w", err)
	}

	strFd := strconv.FormatUint(uint64(nsFile.Fd()), 10)
	fmt.Println("ip link set dev", link, "netns", strFd)
	cmd := exec.Command("ip", "link", "set", "dev", link, "netns", target)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("unable to set link namespace: %w\n%s", err, string(out))
	}

	return nil
}

func SetNetNSByFile(target string) (func() error, error) {
	currentNS, err := os.OpenFile("/proc/self/ns/net", os.O_RDONLY, 0)
	if err != nil {
		return nil, fmt.Errorf("unable to open current network namespace: %w", err)
	}

	unsetFunc := func() error {
		defer currentNS.Close()
		if err != nil {
			return err
		}
		return unix.Setns(int(currentNS.Fd()), unix.CLONE_NEWNET)
	}

	netnsFile, err := os.OpenFile(target, os.O_RDONLY, 0)
	if err != nil {
		return unsetFunc, fmt.Errorf("unable to open network namespace: %w", err)
	}
	defer netnsFile.Close()

	if err := unix.Setns(int(netnsFile.Fd()), unix.CLONE_NEWNET); err != nil {
		return unsetFunc, fmt.Errorf("unable to setns to network namespace: %w", err)
	}

	return unsetFunc, nil
}

func RenameLink(oldName string, newName string) error {
	cmd := exec.Command("ip", "link", "set", "dev", oldName, "name", newName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("unable to rename link %s to %s: %w", oldName, newName, err)
	}

	return nil
}

func SetLinkAddress(link string, address string) error {
	cmd := exec.Command("ip", "addr", "add", address, "dev", link)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("unable to set link address %s to %s: %w", link, address, err)
	}

	return nil
}

func SetLinkUp(link string) error {
	cmd := exec.Command("ip", "link", "set", "dev", link, "up")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("unable to set link %s up: %w", link, err)
	}

	return nil
}

func AddGateway(link string, gateway string) error {
	cmd := exec.Command("ip", "route", "add", "default", "via", gateway, "dev", link)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("unable to add gateway %s to %s: %w", gateway, link, err)
	}

	return nil
}
