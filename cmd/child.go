package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"syscall"

	"github.com/codeclysm/extract"
	"github.com/spf13/cobra"
)

var Child = &cobra.Command{
	Use: "child IMAGE [COMMAND]",
	Run: func(_ *cobra.Command, args []string) {
		image := args[0]
		cmd := args[1]

		child(image, cmd)
	},
}

func child(image string, command string) {
	tar := fmt.Sprintf("./assets/%s.tar.gz", image)

	if _, err := os.Stat(tar); errors.Is(err, os.ErrNotExist) {
		panic(err)
	}

	dir := createTempDir(tar)
	defer os.RemoveAll(dir)
	unTar(tar, dir)

	containerize(dir, command)
}

func createTempDir(name string) string {
	var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

	prefix := nonAlphanumericRegex.ReplaceAllString(name, "_")
	dir, err := os.MkdirTemp("", prefix)
	if err != nil {
		log.Fatal(err)
	}

	return dir
}

func unTar(source string, destination string) error {
	tarball, err := os.Open(source)

	if err != nil {
		return err
	}

	defer tarball.Close()

	ctx := context.Background()
	return extract.Archive(ctx, tarball, destination, nil)
}

func containerize(root string, call string) {
	oldrootHandle, err := os.Open("/")
	if err != nil {
		panic(err)
	}
	defer oldrootHandle.Close()

	syscall.Sethostname([]byte("container"))
	syscall.Chdir(root)
	defer syscall.Fchdir(int(oldrootHandle.Fd()))
	syscall.Chroot(root)
	defer syscall.Chroot(".")
	syscall.Mount("proc", "proc", "proc", 0, "")

	fmt.Printf("Running %s in %s\n", call, root)
	cmd := exec.Command(call)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()

	if err != nil {
		fmt.Println(err)
	}
}
