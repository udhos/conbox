package head

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/udhos/conbox/common"
)

// Run executes the applet.
func Run(tab map[string]common.AppletFunc, args []string) int {

	var help bool
	var files []string
	var lines int

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--":
			files = append(files, args[i+1:]...)
		case arg == "-h":
			help = true
		case arg == "-n":
			if i >= len(args)-1 {
				usage()
				return 1 // exit status
			}
			i++ // consume next arg
			a := args[i]
			v, errConv := strconv.Atoi(a)
			if errConv != nil {
				usage()
				return 1
			}
			lines = v
		case strings.HasPrefix(arg, "-n"):
			count := strings.TrimPrefix(arg, "-n")
			v, errConv := strconv.Atoi(count)
			if errConv != nil {
				usage()
				return 1
			}
			lines = v
		case strings.HasPrefix(arg, "-"):
			count := strings.TrimPrefix(arg, "-")
			v, errConv := strconv.Atoi(count)
			if errConv != nil {
				usage()
				return 1
			}
			lines = v
		default:
			files = append(files, arg)
		}
	}

	if help {
		usage()
		return 2 // exit status
	}

	if lines < 1 {
		lines = 10
	}

	if len(files) < 1 {
		if err := headReader(lines, os.Stdin); err != nil {
			fmt.Printf("head: %v\n", err)
			return 4
		}
		return 0
	}

	var status int

	for _, f := range files {
		if err := headFile(lines, f); err != nil {
			fmt.Printf("head: %v\n", err)
			status = 5
		}
	}

	return status // exit status
}

func usage() {
	common.ShowVersion()
	fmt.Println("head [OPTION]... [FILE]...")
	fmt.Println("  -h   Show command-line help")
}

func headFile(lines int, path string) error {
	f, errOpen := os.Open(path)
	if errOpen != nil {
		return errOpen
	}

	err := headReader(lines, f)

	f.Close()

	return err
}

func headReader(lines int, r io.Reader) error {
	input := bufio.NewReader(r)
	var line int
	for line < lines {
		str, errRead := input.ReadString('\n')
		if str != "" {
			fmt.Print(str)
			line++
		}
		if errRead != nil {
			return errRead
		}
	}
	return nil
}
