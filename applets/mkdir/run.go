// Package mkdir implements an utility.
package mkdir

import (
	"flag"
	"fmt"
	"os"

	"github.com/udhos/conbox/common"
)

// Run executes the applet.
func Run(tab map[string]common.AppletFunc, args []string) int {

	flagSet := flag.NewFlagSet("mkdir", flag.ContinueOnError)
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
		fmt.Println("mkdir: missing operand")
		usage(flagSet)
		return 3
	}

	var status int

	for _, f := range flagSet.Args() {
		if err := os.Mkdir(f, 0777); err != nil {
			fmt.Printf("mkdir: %v\n", err)
			status = 4
		}
	}

	return status // exit status
}

func usage(flagSet *flag.FlagSet) {
	common.ShowVersion()
	fmt.Println("mkdir [OPTION]... FILE...")
	flagSet.PrintDefaults()
}
