package build

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func runError(cmd string, args ...string) ([]byte, error) {
	ecmd := exec.Command(cmd, args...)
	bs, err := ecmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	return bytes.TrimSpace(bs), nil
}

//prints command to the terminal and executes it
func runPrint(cmd string, args ...string) {
	log.Println(cmd, strings.Join(args, " "))
	ecmd := exec.Command(cmd, args...)
	ecmd.Stdout = os.Stdout
	ecmd.Stderr = os.Stderr
	err := ecmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func logAndClose(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Println("error closing:", err)
	}
}

func md5File(file string) error {
	// Can ignore gosec G304 because this function is not used in Grafana, only in the build process.
	//nolint:gosec
	fd, err := os.Open(file)
	if err != nil {
		return err
	}
	defer logAndClose(fd)

	h := md5.New()
	_, err = io.Copy(h, fd)
	if err != nil {
		return err
	}

	out, err := os.Create(file + ".md5")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(out, "%x\n", h.Sum(nil))
	if err != nil {
		return err
	}

	return out.Close()
}

// basically `rm -r`s the list of files provided
func rmr(paths ...string) {
	for _, path := range paths {
		log.Println("rm -r", path)
		if err := os.RemoveAll(path); err != nil {
			log.Println("error deleting folder", path, "error:", err)
		}
	}
}
