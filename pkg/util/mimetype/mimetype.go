package mimetype

import (
	"errors"
	"strings"
)

var (
	ErrorGetType     = errors.New("Unable to get the type")
	ErrorUnknownType = errors.New("Unknown type")

	defaultChecker = newExecChecker()
)

type Checker interface {
	GetMIMEFromFilepath(filepath string) (string, error)
}

func FromFilepath(filepath string) (string, error) {
	return defaultChecker.GetMIMEFromFilepath(filepath)
}

func Base(mime string) string {
	return strings.Split(mime, "/")[0]
}
