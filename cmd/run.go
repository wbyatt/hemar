package cmd

import (
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"
)

var Run = &cobra.Command{
	Use:   "run IMAGE COMMAND",
	Short: "Runs the COMMAND against the container IMAGE",
	Long:  "Always interactive and (for now) ignores the Dockerfile CMD and ENTRYPOINT directives",
	Run: func(cmd *cobra.Command, args []string) {
		commands := append([]string{"child"}, args...)

		run(commands)
	},
}

func run(commands []string) {
	child := exec.Command("/proc/self/exe", commands...)
	child.Stdin = os.Stdin
	child.Stdout = os.Stdout
	child.Stderr = os.Stderr
	child.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
		Credential: &syscall.Credential{
			Uid:    0,
			Gid:    0,
			Groups: []uint32{0},
		},
	}

	if err := child.Start(); err != nil {
		log.Fatal(err)
	}

	child.Wait()
}
