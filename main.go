package main

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
)

func child(image string, command string) {
	tar := fmt.Sprintf("./assets/%s.tar.gz", image)

	if _, err := os.Stat(tar); errors.Is(err, os.ErrNotExist) {
		panic(err)
	}

	dir := createTempDir(tar)
	defer os.RemoveAll(dir)
	must(unTar(tar, dir))

	containerize(dir, command)

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

func main() {
	switch os.Args[1] {
	case "child":
		image := os.Args[2]
		cmd := ""

		if len(os.Args) > 3 {
			cmd = os.Args[3]
		} else {
			buf, err := os.ReadFile(fmt.Sprintf("./assets/%s-cmd", image))
			if err != nil {
				panic(err)
			}
			cmd = string(buf)
		}

		child(image, cmd)
	case "run":
		commands := append([]string{"child"}, os.Args[2:]...)

		run(commands)
	case "pull":
		image := os.Args[2]
		pull(image)
	default:
		panic("invalid command")
	}
}

func pull(image string) {
	cmd := exec.Command("./pull", image)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	must(cmd.Run())
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

	must(syscall.Sethostname([]byte("container")))
	must(syscall.Chdir(root))
	defer must(syscall.Fchdir(int(oldrootHandle.Fd())))
	must(syscall.Chroot(root))
	defer must(syscall.Chroot("."))
	must(syscall.Mount("proc", "proc", "proc", 0, ""))

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

func must(err error) {
	if err != nil {
		panic(err)
	}
}
