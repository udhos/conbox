package cp

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/udhos/conbox/common"
)

// Run executes the applet.
func Run(tab map[string]common.AppletFunc, args []string) int {

	flagSet := flag.NewFlagSet("cp", flag.ContinueOnError)
	help := flagSet.Bool("h", false, "Show command-line help")
	recursive := flagSet.Bool("r", false, "Copy directories recursively")

	if err := flagSet.Parse(args); err != nil {
		usage(flagSet)
		return 1 // exit status
	}

	if *help {
		usage(flagSet)
		return 2 // exit status
	}

	if flagSet.NArg() < 2 {
		fmt.Println("cp: missing operand")
		usage(flagSet)
		return 3
	}

	list := flagSet.Args()

	if errCp := cp(*recursive, list[0], list[1]); errCp != nil {
		fmt.Printf("cp: %v\n", errCp)
		return 4
	}

	return 0
}

func usage(flagSet *flag.FlagSet) {
	common.ShowVersion()
	fmt.Println("cp [OPTION]... SOURCE DESTINATION")
	flagSet.PrintDefaults()
}

func cp(recursive bool, src, dst string) error {

	statSrc, errStatSrc := os.Stat(src)
	if errStatSrc != nil {
		return errStatSrc
	}
	if statSrc.IsDir() {
		return fmt.Errorf("cannot copy from directory")
	}

	statDst, errStatDst := os.Stat(dst)
	if errStatDst == nil && statDst.IsDir() {
		return fmt.Errorf("cannot copy into directory")
	}

	s, errOpen := os.Open(src)
	if errOpen != nil {
		return errOpen
	}
	defer s.Close()

	d, errCreate := os.Create(dst)
	if errCreate != nil {
		return errCreate
	}
	defer d.Close()

	_, errCopy := io.Copy(d, s)

	return errCopy
}
