package host

import (
	"flag"
	"fmt"
	"net"

	"github.com/udhos/conbox/common"
)

// Run executes the applet.
func Run(tab map[string]common.AppletFunc, args []string) int {

	flagSet := flag.NewFlagSet("host", flag.ContinueOnError)
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
		fmt.Println("host: missing operand")
		usage(flagSet)
		return 3
	}

	var status int

	for _, h := range flagSet.Args() {
		if err := host(h); err != nil {
			fmt.Printf("host: %v\n", err)
			status = 4
		}
	}

	return status // exit status
}

func usage(flagSet *flag.FlagSet) {
	common.ShowVersion()
	fmt.Println("host [OPTION]... hostname")
	flagSet.PrintDefaults()
}

func host(hostname string) error {
	addrs, err := net.LookupHost(hostname)
	if err != nil {
		return err
	}
	for _, a := range addrs {
		fmt.Printf("%s has address %s\n", hostname, a)
	}
	return nil
}
