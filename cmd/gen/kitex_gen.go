package gen

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/carlos-yuan/cargen/util/fileUtil"
)

func KitexGen(name, path string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("kitex", "-module", name, name+".proto")
	cmd.Dir = fileUtil.FixPathSeparator(path + "/biz/" + name + "/rpc")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		if strings.Contains(err.Error(), "executable file not found") {
			installKitex()
			KitexGen(name, path)
		} else {
			panic(err.Error() + stderr.String())
		}
	} else {
		cmd.Process.Kill()
	}
}

func installKitex() {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("go", "install", "github.com/cloudwego/kitex/tool/cmd/kitex@v0.6.2")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		panic(err.Error() + stderr.String())
	} else {
		cmd.Process.Kill()
	}
}
