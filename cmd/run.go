package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/wbyatt/hemar/container"
	"github.com/wbyatt/hemar/image"
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
	container.MountFilesystem()
	commands = append([]string{"child"}, append([]string{container.Digest}, commands[1:]...)...)

	child := exec.Command("/proc/self/exe", commands...)
	child.Stdin, child.Stdout, child.Stderr = os.Stdin, os.Stdout, os.Stderr

	child.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
		Credential: &syscall.Credential{
			Uid:    0,
			Gid:    0,
			Groups: []uint32{0},
		},
	}

	child.Run()
}
