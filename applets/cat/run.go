// Package cat implements an utility.
package cat

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/udhos/conbox/common"
)

// Run executes the applet.
func Run(tab map[string]common.AppletFunc, args []string) int {

	flagSet := flag.NewFlagSet("cat", flag.ContinueOnError)
	helpFlag := flagSet.Bool("h", false, "Show command-line help")

	if err := flagSet.Parse(args); err != nil {
		usage(flagSet)
		return 1 // exit status
	}

	if *helpFlag {
		usage(flagSet)
		return 2 // exit status
	}

	if flagSet.NArg() < 1 {
		if err := catReader(os.Stdin); err != nil {
			fmt.Printf("cat: %v\n", err)
			return 4
		}
		return 0
	}

	var status int

	for _, f := range flagSet.Args() {
		if err := catFile(f); err != nil {
			fmt.Printf("cat: %v\n", err)
			status = 5
		}
	}

	return status // exit status
}

func usage(flagSet *flag.FlagSet) {
	common.ShowVersion()
	fmt.Println("cat [OPTION]... [FILE]...")
	flagSet.PrintDefaults()
}

func catFile(path string) error {
	f, errOpen := os.Open(path)
	if errOpen != nil {
		return errOpen
	}

	err := catReader(f)

	f.Close()

	return err
}

func catReader(r io.Reader) error {
	_, err := io.Copy(os.Stdout, r)
	return err
}
