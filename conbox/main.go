package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/udhos/conbox/applets/cat"
)

const conboxVersion = "0.0"

func main() {

	fmt.Printf("conbox: version %s runtime %s GOMAXPROC=%d OS=%s ARCH=%s\n", conboxVersion, runtime.Version(), runtime.GOMAXPROCS(0), runtime.GOOS, runtime.GOARCH)

	appletTable := loadApplets()

	// 1. try basename
	appletName := filepath.Base(os.Args[0])
	if applet, found := appletTable[appletName]; found {
		run(applet, os.Args[1:])
		return
	}

	// 2. try arg 1
	if len(os.Args) > 1 {
		appletName = os.Args[1]
		if applet, found := appletTable[appletName]; found {
			run(applet, os.Args[2:])
			return
		}
		fmt.Printf("conbox: arg 1: applet '%s' not found\n", appletName)
	} else {
		fmt.Printf("conbox: basename: applet '%s' not found\n", appletName)
	}

	fmt.Println("registered applets:")
	for n := range appletTable {
		fmt.Printf("%s ", n)
	}
	fmt.Println()
}

func run(applet appletFunc, args []string) {
	exitCode := applet(args)
	if exitCode != 0 {
		os.Exit(exitCode)
	}
}

type appletFunc func(args []string) int

func loadApplets() map[string]appletFunc {
	tab := map[string]appletFunc{}
	tab["cat"] = cat.Run
	return tab
}
