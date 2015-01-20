// Using command line util `file` to get mime type.
package mimetype

import (
	"os/exec"
	"strings"
)

type execChecker struct {
	command string
}

func newExecChecker() *execChecker {
	return &execChecker{"file"}
}

func (c *execChecker) GetMIMEFromFilepath(filepath string) (string, error) {
	t, err := exec.Command(c.command, "--mime-type", filepath).CombinedOutput()
	if err != nil {
		return "", ErrorGetType
	}

	if m := strings.Split(string(t), ": ")[1]; m != "" {
		return m, nil
	}
	return "", ErrorUnknownType
}
