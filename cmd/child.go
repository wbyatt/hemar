package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"
)

var Child = &cobra.Command{
	Use: "child IMAGE [COMMAND]",
	Run: func(_ *cobra.Command, args []string) {

		container := args[0]
		cmd := args[1]

		child(container, cmd)
	},
}

func child(container string, command string) {

	containerize(container, command)
}

func containerize(container string, call string) {
	dir := fmt.Sprintf("%s/%s/rootfs", containersPath, container)
	fmt.Println("Launching container from", dir)

	oldrootHandle, err := os.Open("/")
	if err != nil {
		panic(err)
	}
	defer oldrootHandle.Close()

	must(syscall.Sethostname([]byte(container[:8])))
	must(syscall.Chdir(dir))
	defer must(syscall.Fchdir(int(oldrootHandle.Fd())))
	must(syscall.Chroot(dir))
	defer must(syscall.Chroot("."))
	must(syscall.Mount("proc", "proc", "proc", 0, ""))

	fmt.Printf("Running %s in %s\n", call, dir)
	cmd := exec.Command(call)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()

	if err != nil {
		fmt.Println(err)
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
