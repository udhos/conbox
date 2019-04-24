package ls

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/udhos/conbox/common"
)

// Run executes the applet.
func Run(tab map[string]common.AppletFunc, args []string) int {

	flagSet := flag.NewFlagSet("ls", flag.ContinueOnError)
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
		err := ls(".")
		if err != nil {
			fmt.Printf("ls: %v\n", err)
			return 3
		}
		return 0
	}

	var status int

	for _, f := range flagSet.Args() {
		if err := ls(f); err != nil {
			fmt.Printf("ls: %v\n", err)
			status = 4
		}
	}

	return status // exit status
}

func usage(flagSet *flag.FlagSet) {
	common.ShowVersion()
	fmt.Println("ls [OPTION]... [FILE]...")
	flagSet.PrintDefaults()
}

func ls(path string) error {

	stat, errStat := os.Stat(path)
	if errStat != nil {
		return errStat
	}

	if stat.IsDir() {
		return lsDir(path)
	}

	fmt.Println(path)

	return nil
}

func lsDir(path string) error {

	list, errList := ioutil.ReadDir(path)
	if errList != nil {
		return errList
	}

	for _, f := range list {
		//p := path + "/" + f.Name()
		fmt.Println(f.Name())
	}

	return nil
}
