// Package common holds some utilities.
package common

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

const version = "0.1.0"

// AppletFunc type is the signature for the applet Run() function.
type AppletFunc func(tab map[string]AppletFunc, args []string) int

// ShowVersion prints conbox version.
func ShowVersion() {
	fmt.Printf("conbox: version %s runtime %s GOMAXPROC=%d OS=%s ARCH=%s\n", version, runtime.Version(), runtime.GOMAXPROCS(0), runtime.GOOS, runtime.GOARCH)
}

// Tokenize parses an input line into tokens.
func Tokenize(line string) []string {
	return strings.Fields(line) // FIXME WRITEME
}

// IsDir reports if path exists as directory.
func IsDir(path string) bool {
	stat, errStat := os.Stat(path)
	return errStat == nil && stat.IsDir()
}
