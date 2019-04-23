package common

import (
	"fmt"
	"runtime"
)

const version = "0.0"

// ShowVersion prints conbox version.
func ShowVersion() {
	fmt.Printf("conbox: version %s runtime %s GOMAXPROC=%d OS=%s ARCH=%s\n", version, runtime.Version(), runtime.GOMAXPROCS(0), runtime.GOOS, runtime.GOARCH)
}
