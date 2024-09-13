package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/wbyatt/hemar/container"
	"github.com/wbyatt/hemar/image"
	"github.com/wbyatt/hemar/network"
)

var Run = &cobra.Command{
	Use:   "run IMAGE COMMAND",
	Short: "Runs the COMMAND against the container IMAGE",
	Long:  "Always interactive and (for now) ignores the Dockerfile CMD and ENTRYPOINT directives",
	Run: func(cmd *cobra.Command, args []string) {
		run(args)
	},
}

func run(commands []string) {
	fmt.Println("Running command", commands)

	imageName := commands[0]
	container := container.NewContainer(&container.ContainerConfig{
		Image: image.NewImage(imageName, "latest"),
	})

	network.SetupBridge("hemar0")
	defer network.TeardownBridge("hemar0")
	network.SetupNAT("hemar0", "eth0")

	unmountNetwork, err := container.SetupNetwork("hemar0")
	if err != nil {
		log.Fatalf("Failed to setup network: %v", err)
	}
	defer unmountNetwork()

	container.MountFilesystem()
	defer container.Cleanup()

	commands = append([]string{"child"}, append([]string{container.Digest}, commands[1:]...)...)

	child := exec.Command("/proc/self/exe", commands...)
	child.Stdin, child.Stdout, child.Stderr = os.Stdin, os.Stdout, os.Stderr

	child.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWIPC,
		Credential: &syscall.Credential{
			Uid:    0,
			Gid:    0,
			Groups: []uint32{0},
		},
	}

	child.Run()
}
