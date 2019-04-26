package touch

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/udhos/conbox/common"
)

// Run executes the applet.
func Run(tab map[string]common.AppletFunc, args []string) int {

	flagSet := flag.NewFlagSet("touch", flag.ContinueOnError)
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
		fmt.Println("touch: missing operand")
		usage(flagSet)
		return 3
	}

	var status int

	for _, f := range flagSet.Args() {
		if err := touch(f); err != nil {
			fmt.Printf("touch: %v\n", err)
			status = 4
		}
	}

	return status // exit status
}

func usage(flagSet *flag.FlagSet) {
	common.ShowVersion()
	fmt.Println("touch [OPTION]... FILE...")
	flagSet.PrintDefaults()
}

func touch(path string) error {

	_, errStat := os.Stat(path)
	if errStat != nil {
		if os.IsNotExist(errStat) {
			f, errCreate := os.Create(path)
			if errCreate != nil {
				return errCreate
			}
			return f.Close()
		}

		return errStat
	}

	now := time.Now()

	return os.Chtimes(path, now, now)
}
