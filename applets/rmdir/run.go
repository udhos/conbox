// Package rmdir implements an utility.
package rmdir

import (
	"flag"
	"fmt"
	"os"

	"github.com/udhos/conbox/common"
)

// Run executes the applet.
func Run(tab map[string]common.AppletFunc, args []string) int {

	flagSet := flag.NewFlagSet("rmdir", flag.ContinueOnError)
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
		fmt.Println("rmdir: missing operand")
		usage(flagSet)
		return 3
	}

	var status int

	for _, f := range flagSet.Args() {
		if err := rmdir(f); err != nil {
			fmt.Printf("rmdir: %v\n", err)
			status = 4
		}
	}

	return status // exit status
}

func usage(flagSet *flag.FlagSet) {
	common.ShowVersion()
	fmt.Println("rmdir [OPTION]... FILE...")
	flagSet.PrintDefaults()
}

func rmdir(path string) error {

	stat, errStat := os.Stat(path)
	if errStat != nil {
		return errStat
	}

	if !stat.IsDir() {
		return fmt.Errorf("not a directory: %s", path)
	}

	return os.Remove(path)
}
