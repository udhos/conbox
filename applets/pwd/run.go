package pwd

import (
	"fmt"
	"os"
)

// Run executes the applet.
func Run(args []string) int {

	dir, err := os.Getwd()

	if err != nil {
		fmt.Printf("pwd: %v\n", err)
		return 1
	}

	fmt.Println(dir)

	return 0 // exit status
}
