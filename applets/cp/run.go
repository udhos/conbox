package cp

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

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

	if errCp := cp(*recursive, list); errCp != nil {
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

func isDir(path string) bool {
	stat, errStat := os.Stat(path)
	return errStat == nil && stat.IsDir()
}

func cp(recursive bool, list []string) error {
	if len(list) < 2 {
		return fmt.Errorf("cp: missing operand")
	}

	last := len(list) - 1
	dst := list[last]
	dstDir := isDir(dst)

	srcMultiple := len(list) > 2
	if srcMultiple {
		if !dstDir {
			// create dir
			if errMkdir := os.Mkdir(dst, 0777); errMkdir != nil {
				return errMkdir
			}
			dstDir = true
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
		errCp := cpFile(src, dstFull)
		if errCp != nil {
			fmt.Printf("cp: %v\n", errCp)
			errLast = errCp
		}
	}
	return errLast
}

func cpFile(src, dst string) error {

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
