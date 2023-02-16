// Package chmod implements an utility.
package chmod

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/udhos/conbox/common"
)

// Run executes the applet.
func Run(tab map[string]common.AppletFunc, args []string) int {

	flagSet := flag.NewFlagSet("chmod", flag.ContinueOnError)
	help := flagSet.Bool("h", false, "Show command-line help")
	recursive := flagSet.Bool("r", false, "Change directories recursively")

	if err := flagSet.Parse(args); err != nil {
		usage(flagSet)
		return 1 // exit status
	}

	if *help {
		usage(flagSet)
		return 2 // exit status
	}

	if flagSet.NArg() < 2 {
		fmt.Println("chmod: missing operand")
		usage(flagSet)
		return 3
	}

	var status int

	list := flagSet.Args()
	mode := list[0]

	m, errConv := strconv.ParseUint(mode, 8, 32)
	if errConv != nil {
		fmt.Printf("chmod: bad file mode: %s: %v", mode, errConv)
		return 4
	}

	fileMode := os.FileMode(uint32(m))

	for _, f := range list[1:] {
		if err := chmod(f, fileMode, *recursive); err != nil {
			fmt.Printf("chmod: %v\n", err)
			status = 5
		}
	}

	return status // exit status
}

func usage(flagSet *flag.FlagSet) {
	common.ShowVersion()
	fmt.Println("chmod [OPTION]... MODE FILE...")
	flagSet.PrintDefaults()
}

func chmod(path string, mode os.FileMode, recursive bool) error {

	stat, errStat := os.Stat(path)
	if errStat != nil {
		return errStat
	}

	if stat.IsDir() {
		if recursive {
			return chmodDir(path, mode)
		}
	}

	return os.Chmod(path, mode)
}

func chmodDir(path string, mode os.FileMode) error {

	list, errList := ioutil.ReadDir(path)
	if errList != nil {
		return errList
	}

	var errSave error

	// chmod children
	for _, f := range list {
		p := filepath.Join(path, f.Name())
		if errChmod := chmod(p, mode, true); errChmod != nil {
			errSave = errChmod
			fmt.Printf("chmod: %v\n", errChmod)
		}
	}

	// chmod path
	errChmod := os.Chmod(path, mode)
	if errChmod == nil {
		return errSave
	}
	return errChmod
}
