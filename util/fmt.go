package util

import (
	"bytes"
	"os/exec"
)

func GoFmt(path string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("gofmt", "-l", "-w", "-s", path)
	cmd.Dir = path
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	cmd.Process.Kill()
	if err != nil {
		panic(err.Error() + stderr.String())
	}
}
