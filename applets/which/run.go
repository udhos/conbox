// Package which implements an utility.
package which

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/udhos/conbox/common"
)

// Run executes the applet.
func Run(tab map[string]common.AppletFunc, args []string) int {

	flagSet := flag.NewFlagSet("which", flag.ContinueOnError)
	help := flagSet.Bool("h", false, "Show command-line help")

	if err := flagSet.Parse(args); err != nil {
		usage(flagSet)
		return 1 // exit status
	}

	if *help {
		usage(flagSet)
		return 2 // exit status
	}

	if flagSet.NArg() < 1 {
		return 3 // exit status
	}

	var status int

	paths := strings.Split(os.Getenv("PATH"), ":")

	for _, f := range flagSet.Args() {
		if s := which(f, paths); s != 0 {
			status = s
		}
	}

	return status // exit status
}

func usage(flagSet *flag.FlagSet) {
	common.ShowVersion()
	fmt.Println("which [OPTION]... [FILE]...")
	flagSet.PrintDefaults()
}

func which(cmd string, paths []string) int {
	for _, p := range paths {
		full := filepath.Join(p, cmd)
		info, errStat := os.Stat(full)
		if errStat != nil {
			continue
		}
		if info.IsDir() {
			continue
		}
		if (info.Mode() & 0111) == 0 {
			continue
		}
		fmt.Println(full)
		return 0
	}
	return 1
}
