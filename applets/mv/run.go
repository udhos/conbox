// Package mv implements an utility.
package mv

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/udhos/conbox/common"
)

// Run executes the applet.
func Run(tab map[string]common.AppletFunc, args []string) int {

	flagSet := flag.NewFlagSet("mv", flag.ContinueOnError)
	help := flagSet.Bool("h", false, "Show command-line help")

	if err := flagSet.Parse(args); err != nil {
		usage(flagSet)
		return 1 // exit status
	}

	if *help {
		usage(flagSet)
		return 2 // exit status
	}

	if flagSet.NArg() < 2 {
		fmt.Println("mv: missing operand")
		usage(flagSet)
		return 3
	}

	list := flagSet.Args()

	if errMv := mv(list); errMv != nil {
		fmt.Printf("mv: %v\n", errMv)
		return 4
	}

	return 0
}

func usage(flagSet *flag.FlagSet) {
	common.ShowVersion()
	fmt.Println("mv [OPTION]... SOURCE... DESTINATION")
	flagSet.PrintDefaults()
}

func mv(list []string) error {
	if len(list) < 2 {
		return fmt.Errorf("mv: missing operand")
	}

	last := len(list) - 1
	dst := list[last]
	dstDir := common.IsDir(dst)

	srcMultiple := len(list) > 2
	if srcMultiple {
		if !dstDir {
			return fmt.Errorf("%s: destination is not a directory", dst)
		}
	}

	var errLast error
	for i := 0; i < last; i++ {
		src := list[i]
		var dstFull string
		if dstDir {
			dstFull = filepath.Join(dst, filepath.Base(src))
		} else {
			dstFull = dst
		}
		errMv := mvFile(src, dstFull)
		if errMv != nil {
			fmt.Printf("mv: %v\n", errMv)
			errLast = errMv
		}
	}
	return errLast
}

func mvFile(src, dst string) error {
	return os.Rename(src, dst)
}
