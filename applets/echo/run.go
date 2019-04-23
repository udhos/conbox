package echo

import (
	"flag"
	"fmt"
)

// Run executes the applet.
func Run(args []string) int {

	flagSet := flag.NewFlagSet("echo", flag.ExitOnError)
	suppressNewline := flagSet.Bool("n", false, "Do not output the trailing newline")

	if err := flagSet.Parse(args); err != nil {
		return 1 // exit status
	}

	argList := flagSet.Args()
	if len(argList) > 0 {
		fmt.Print(argList[0])
	}

	for i := 1; i < len(argList); i++ {
		fmt.Print(" ")
		fmt.Print(argList[i])
	}

	if !*suppressNewline {
		fmt.Println()
	}

	return 0 // exit status
}
