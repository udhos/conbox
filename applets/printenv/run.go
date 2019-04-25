package printenv

import (
	"flag"
	"fmt"
	"os"

	"github.com/udhos/conbox/common"
)

// Run executes the applet.
func Run(tab map[string]common.AppletFunc, args []string) int {

	flagSet := flag.NewFlagSet("printenv", flag.ContinueOnError)
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
		for _, e := range os.Environ() {
			fmt.Println(e)
		}
		return 0 // exit status
	}

	var status int

	for _, v := range flagSet.Args() {
		vv := os.Getenv(v)
		if vv == "" {
			status = 3
			continue
		}
		fmt.Println(vv)
	}

	return status // exit status
}

func usage(flagSet *flag.FlagSet) {
	common.ShowVersion()
	fmt.Println("printenv [OPTION]... [VARIABLE]...")
	flagSet.PrintDefaults()
}
