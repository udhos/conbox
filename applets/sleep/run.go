package sleep

import (
	"flag"
	"fmt"
	"strconv"
	"time"

	"github.com/udhos/conbox/common"
)

// Run executes the applet.
func Run(tab map[string]common.AppletFunc, args []string) int {

	flagSet := flag.NewFlagSet("sleep", flag.ContinueOnError)
	help := flagSet.Bool("h", false, "Show command-line help")

	if err := flagSet.Parse(args); err != nil {
		usage(flagSet)
		return 1 // exit status
	}

	if *help {
		usage(flagSet)
		return 2 // exit status
	}

	if flagSet.NArg() != 1 {
		fmt.Println("sleep: single operand required")
		usage(flagSet)
		return 3
	}

	sec := flagSet.Args()[0]

	f, errParse := strconv.ParseFloat(sec, 64)
	if errParse != nil {
		fmt.Printf("sleep: invalid operand: %s: %v\n", sec, errParse)
		return 4
	}

	if f < 0 {
		fmt.Printf("sleep: negative operand: %f\n", f)
		return 5
	}

	dur := time.Duration(f * float64(time.Second))

	time.Sleep(dur)

	return 0 // exit status
}

func usage(flagSet *flag.FlagSet) {
	common.ShowVersion()
	fmt.Println("sleep SECONDS")
	flagSet.PrintDefaults()
}
