package gen

import (
	"bytes"
	"github.com/carlos-yuan/cargen/util"
	"os/exec"
)

func KitexGen(name, path string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("kitex", "-module", name, name+".proto")
	cmd.Dir = util.FixPathSeparator(path + "/biz/" + name + "/rpc")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	cmd.Process.Kill()
	if err != nil {
		panic(err.Error() + stderr.String())
	}
}
