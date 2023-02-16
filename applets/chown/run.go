// Package chown implements an utility.
package chown

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/udhos/conbox/common"
)

// Run executes the applet.
func Run(tab map[string]common.AppletFunc, args []string) int {

	var help, recursive bool
	flagSet := flag.NewFlagSet("chown", flag.ContinueOnError)
	flagSet.BoolVar(&help, "h", false, "Show command-line help")
	flagSet.BoolVar(&recursive, "r", false, "Change directories recursively")

	if err := flagSet.Parse(args); err != nil {
		usage(flagSet)
		return 1 // exit status
	}

	if help {
		usage(flagSet)
		return 2 // exit status
	}

	if flagSet.NArg() < 2 {
		fmt.Println("chown: missing operand")
		usage(flagSet)
		return 3
	}

	var status int

	list := flagSet.Args()
	owner := list[0]

	ug := strings.SplitN(owner, ":", 2)
	if len(ug) != 2 {
		ug = strings.SplitN(owner, ".", 2)
	}

	var strUser, strGroup string

	if len(ug) > 0 {
		strUser = ug[0]
	}
	if len(ug) > 1 {
		strGroup = ug[1]
	}

	uid, errUser := solveUser(strUser)
	if errUser != nil {
		fmt.Printf("chown: %v\n", errUser)
		return 4
	}

	gid, errGroup := solveGroup(strGroup)
	if errGroup != nil {
		fmt.Printf("chown: %v\n", errGroup)
		return 5
	}

	for _, f := range list[1:] {
		if err := chown(f, uid, gid, recursive); err != nil {
			fmt.Printf("chown: %v\n", err)
			status = 6
		}
	}

	return status // exit status
}

func solveUser(name string) (int, error) {
	if name == "" {
		return -1, nil
	}
	id, errConv := strconv.Atoi(name)
	if errConv == nil {
		return id, nil
	}
	u, errUser := user.Lookup(name)
	if errUser != nil {
		return -1, errUser
	}
	return strconv.Atoi(u.Uid)
}

func solveGroup(name string) (int, error) {
	if name == "" {
		return -1, nil
	}
	id, errConv := strconv.Atoi(name)
	if errConv == nil {
		return id, nil
	}
	g, errGroup := user.LookupGroup(name)
	if errGroup != nil {
		return -1, errGroup
	}
	return strconv.Atoi(g.Gid)
}

func usage(flagSet *flag.FlagSet) {
	common.ShowVersion()
	fmt.Println("chown [OPTION]... OWNER:GROUP FILE...")
	flagSet.PrintDefaults()
}

func chown(path string, uid, gid int, recursive bool) error {

	stat, errStat := os.Stat(path)
	if errStat != nil {
		return errStat
	}

	if stat.IsDir() {
		if recursive {
			return chownDir(path, uid, gid)
		}
	}

	return os.Chown(path, uid, gid)
}

func chownDir(path string, uid, gid int) error {

	list, errList := ioutil.ReadDir(path)
	if errList != nil {
		return errList
	}

	var errSave error

	// chown children
	for _, f := range list {
		p := filepath.Join(path, f.Name())
		if errChown := chown(p, uid, gid, true); errChown != nil {
			errSave = errChown
			fmt.Printf("chown: %v\n", errChown)
		}
	}

	// chown path
	errChown := os.Chown(path, uid, gid)
	if errChown == nil {
		return errSave
	}
	return errChown
}
