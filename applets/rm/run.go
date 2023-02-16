// Package rm implements an utility.
package rm

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/udhos/conbox/common"
)

// Run executes the applet.
func Run(tab map[string]common.AppletFunc, args []string) int {

	flagSet := flag.NewFlagSet("rm", flag.ContinueOnError)
	help := flagSet.Bool("h", false, "Show command-line help")
	recursive := flagSet.Bool("r", false, "Remove recursively into directories")

	if err := flagSet.Parse(args); err != nil {
		usage(flagSet)
		return 1 // exit status
	}

	if *help {
		usage(flagSet)
		return 2 // exit status
	}

	if flagSet.NArg() < 1 {
		fmt.Println("rm: missing operand")
		usage(flagSet)
		return 3
	}

	var status int

	for _, f := range flagSet.Args() {
		if err := rm(f, *recursive); err != nil {
			fmt.Printf("rm: %v\n", err)
			status = 4
		}
	}

	return status // exit status
}

func usage(flagSet *flag.FlagSet) {
	common.ShowVersion()
	fmt.Println("rm [OPTION]... FILE...")
	flagSet.PrintDefaults()
}

func rm(path string, recursive bool) error {

	stat, errStat := os.Stat(path)
	if errStat != nil {
		return errStat
	}

	if stat.IsDir() {
		if recursive {
			return rmDir(path)
		}
		return fmt.Errorf("can't remove a directory: %s", path)
	}

	return os.Remove(path)
}

func rmDir(path string) error {

	list, errList := ioutil.ReadDir(path)
	if errList != nil {
		return errList
	}

	for _, f := range list {
		p := filepath.Join(path, f.Name())
		if errRm := rm(p, true); errRm != nil {
			fmt.Printf("rm: %v\n", errRm)
		}
	}

	return os.Remove(path)
}
