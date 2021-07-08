package comtools

import (
	"bytes"
	"os/exec"
)

func RunCmd(cmdStr string) (content string, err error) {
	content = ""
	cmd := exec.Command("mybash", "-c", cmdStr)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err == nil {
		content = out.String()
	}
	return
}
