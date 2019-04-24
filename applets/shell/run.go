package shell

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/udhos/conbox/common"
)

// Run executes the applet.
func Run(tab map[string]common.AppletFunc, args []string) int {

	flagSet := flag.NewFlagSet("shell", flag.ContinueOnError)
	help := flagSet.Bool("h", false, "Show command-line help")
	interactive := flagSet.Bool("i", false, "Force interactive mode")

	if err := flagSet.Parse(args); err != nil {
		usage(flagSet)
		return 1 // exit status
	}

	if *help {
		usage(flagSet)
		return 2 // exit status
	}

	if flagSet.NArg() < 1 {
		return loop(tab, os.Stdin, true)
	}

	list := flagSet.Args()
	input := list[0]
	f, errOpen := os.Open(input)
	if errOpen != nil {
		fmt.Printf("shell: %v\n", errOpen)
		return 3 // exit status
	}

	defer f.Close()

	info, errStat := f.Stat()
	if errStat != nil {
		fmt.Printf("shell: %v\n", errStat)
		return 4 // exit status
	}

	if info.IsDir() {
		fmt.Printf("shell: %s: is a directory\n", input)
		return 5 // exit status
	}

	return loop(tab, f, *interactive)
}

func usage(flagSet *flag.FlagSet) {
	common.ShowVersion()
	fmt.Println("shell [OPTION]... [FILE]")
	flagSet.PrintDefaults()
}

func loop(tab map[string]common.AppletFunc, r io.Reader, interactive bool) int {

	if interactive {
		common.ShowVersion()
		fmt.Println("welcome to conbox shell")
		fmt.Println()
	}

	input := bufio.NewReader(r)

LOOP:
	for {
		if interactive {
			fmt.Print("conbox shell$ ")
		}

		line, errInput := input.ReadString('\n')
		switch errInput {
		case nil:
		case io.EOF:
			break LOOP
		default:
			fmt.Printf("shell: %v\n", errInput)
			return 10
		}

		parameters := strings.Fields(line)

		if len(parameters) < 1 {
			continue // empty line
		}

		p0 := parameters[0]

		if p0 == "" {
			continue // blank line
		}

		if p0[0] == '#' {
			continue // comment
		}

		exit, status := execute(tab, parameters)
		if exit {
			return status
		}
	}

	return 0
}

func execute(tab map[string]common.AppletFunc, params []string) (bool, int) {

	prog := params[0]

	// 1. lookup shell built-in

	switch prog {
	case "exit":
		return true, 0
	}

	// 2. lookup conbox applet

	if applet, found := tab[prog]; found {
		return false, applet(tab, params[1:])
	}

	// 3. lookup PATH

	fmt.Printf("shell: not found: %s\n", prog)

	return false, 0
}
