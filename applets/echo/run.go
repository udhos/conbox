package echo

import (
	"fmt"

	"github.com/udhos/conbox/common"
)

// Run executes the applet.
func Run(tab map[string]common.AppletFunc, args []string) int {

	var suppressNewline bool

	if len(args) > 0 {
		suppressNewline = args[0] == "-n"
		args = args[1:]
	}

	if len(args) > 0 {
		fmt.Print(args[0])
	}

	for i := 1; i < len(args); i++ {
		fmt.Print(" ")
		fmt.Print(args[i])
	}

	if !suppressNewline {
		fmt.Println()
	}

	return 0 // exit status
}
