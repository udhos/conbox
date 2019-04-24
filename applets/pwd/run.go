package pwd

import (
	"fmt"
	"os"

	"github.com/udhos/conbox/common"
)

// Run executes the applet.
func Run(tab map[string]common.AppletFunc, args []string) int {

	dir, err := os.Getwd()

	if err != nil {
		fmt.Printf("pwd: %v\n", err)
		return 1
	}

	fmt.Println(dir)

	return 0 // exit status
}
