// Package ls implements an utility.
package ls

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"

	"github.com/udhos/conbox/common"
)

// Run executes the applet.
func Run(tab map[string]common.AppletFunc, args []string) int {

	var help, long bool

	flagSet := flag.NewFlagSet("ls", flag.ContinueOnError)
	flagSet.BoolVar(&help, "h", false, "Show command-line help")
	flagSet.BoolVar(&long, "l", false, "Long list mode")

	if err := flagSet.Parse(args); err != nil {
		usage(flagSet)
		return 1 // exit status
	}

	if help {
		usage(flagSet)
		return 2 // exit status
	}

	if flagSet.NArg() < 1 {
		err := ls(".", long)
		if err != nil {
			fmt.Printf("ls: %v\n", err)
			return 3
		}
		return 0
	}

	var status int

	for _, f := range flagSet.Args() {
		if err := ls(f, long); err != nil {
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

func ls(path string, long bool) error {

	stat, errStat := os.Stat(path)
	if errStat != nil {
		return errStat
	}

	if stat.IsDir() {
		return lsDir(path, long)
	}

	showInfo(path, long, stat)

	return nil
}

func lsDir(path string, long bool) error {

	list, errList := ioutil.ReadDir(path)
	if errList != nil {
		return errList
	}

	for _, f := range list {
		show(filepath.Join(path, f.Name()), long)
	}

	return nil
}

func show(path string, long bool) {

	stat, errStat := os.Stat(path)
	if errStat != nil {
		fmt.Printf("ls: %v\n", errStat)
		return
	}

	showInfo(path, long, stat)
}

func showInfo(path string, long bool, info os.FileInfo) {
	if long {
		showFileMode(info)
	}

	fmt.Println(path)
}

func showFileMode(info os.FileInfo) {
	showBool(info.IsDir(), "d")
	mode := info.Mode()
	showBool(mode&0400 != 0, "r")
	showBool(mode&0200 != 0, "w")
	showBool(mode&0100 != 0, "x")
	showBool(mode&0040 != 0, "r")
	showBool(mode&0020 != 0, "w")
	showBool(mode&0010 != 0, "x")
	showBool(mode&0004 != 0, "r")
	showBool(mode&0002 != 0, "w")
	showBool(mode&0001 != 0, "x")

	uid := -1
	gid := -1

	sys := info.Sys()
	if s, ok := sys.(*syscall.Stat_t); ok {
		uid = int(s.Uid)
		gid = int(s.Gid)
	}

	fmt.Printf(" %4d %4d ", uid, gid)
}

func showBool(value bool, ifTrue string) {
	if value {
		fmt.Print(ifTrue)
	} else {
		fmt.Print("-")
	}
}
