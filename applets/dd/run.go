package dd

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/udhos/conbox/common"
)

// Run executes the applet.
func Run(tab map[string]common.AppletFunc, args []string) int {

	flagSet := flag.NewFlagSet("dd", flag.ContinueOnError)
	help := flagSet.Bool("h", false, "Show command-line help")

	if err := flagSet.Parse(args); err != nil {
		usage(flagSet)
		return 1 // exit status
	}

	if *help {
		usage(flagSet)
		return 2 // exit status
	}

	var src, dst string

	list := flagSet.Args()
	for _, arg := range list {
		value := strings.TrimPrefix(arg, "if=")
		if value != arg {
			src = value
			continue
		}
		value = strings.TrimPrefix(arg, "of=")
		if value != arg {
			dst = value
			continue
		}
	}

	if errDd := dd(src, dst); errDd != nil {
		fmt.Printf("dd: %v\n", errDd)
		return 3
	}

	return 0
}

func usage(flagSet *flag.FlagSet) {
	common.ShowVersion()
	fmt.Println("dd [OPERAND]...")
	flagSet.PrintDefaults()
}

func dd(src, dst string) error {
	var s io.Reader
	var d io.Writer

	if src == "" || src == "-" {
		s = os.Stdin
	} else {
		statSrc, errStatSrc := os.Stat(src)
		if errStatSrc != nil {
			return errStatSrc
		}
		if statSrc.IsDir() {
			return fmt.Errorf("cannot copy from directory")
		}
		in, errOpen := os.Open(src)
		if errOpen != nil {
			return errOpen
		}
		defer in.Close()
		s = in
	}

	if dst == "" || dst == "-" {
		d = os.Stdout
	} else {
		statDst, errStatDst := os.Stat(dst)
		if errStatDst == nil && statDst.IsDir() {
			return fmt.Errorf("cannot copy into directory")
		}
		out, errCreate := os.Create(dst)
		if errCreate != nil {
			return errCreate
		}
		defer out.Close()
		d = out
	}

	return ddStream(s, d)
}

func ddStream(src io.Reader, dst io.Writer) error {
	_, errCopy := io.Copy(dst, src)
	return errCopy
}
