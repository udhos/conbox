package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/udhos/conbox/common"
)

func main() {

	appletTable := loadApplets()

	// 1. try basename
	appletName := filepath.Base(os.Args[0])
	if applet, found := appletTable[appletName]; found {
		run(appletTable, applet, os.Args[1:])
		return
	}

	if appletName != "conbox" {
		common.ShowVersion()
		fmt.Printf("conbox: basename: applet '%s' not found\n", appletName)
		usage(appletTable)
		os.Exit(1)
	}

	// 2. try arg 1
	if len(os.Args) > 1 {
		arg := os.Args[1]
		switch arg {
		case "-h":
			common.ShowVersion()
			usage(appletTable)
			return
		case "-l":
			listApplets(appletTable, "\n")
			return
		}
		appletName = arg
		if applet, found := appletTable[appletName]; found {
			run(appletTable, applet, os.Args[2:])
			return
		}
		common.ShowVersion()
		fmt.Printf("conbox: arg 1: applet '%s' not found\n", appletName)
		usage(appletTable)
		os.Exit(2)
	}

	common.ShowVersion()
	usage(appletTable)
	os.Exit(3)
}

func usage(tab map[string]common.AppletFunc) {
	fmt.Println()
	fmt.Println("usage: conbox APPLET [ARG]... : run APPLET")
	fmt.Println("       conbox -h              : show command-line help")
	fmt.Println("       conbox -l              : list applets")
	fmt.Println()
	fmt.Println("conbox: registered applets:")
	listApplets(tab, " ")
	fmt.Println()
}

func listApplets(tab map[string]common.AppletFunc, sep string) {
	var list []string
	for n := range tab {
		list = append(list, n)
	}
	sort.Strings(list)
	for _, n := range list {
		fmt.Printf("%s%s", n, sep)
	}
}

func run(tab map[string]common.AppletFunc, applet common.AppletFunc, args []string) {
	exitCode := applet(tab, args)
	if exitCode != 0 {
		os.Exit(exitCode)
	}
}
