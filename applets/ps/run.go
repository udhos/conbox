package ps

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/udhos/conbox/common"
)

// Run executes the applet.
func Run(tab map[string]common.AppletFunc, args []string) int {

	flagSet := flag.NewFlagSet("ps", flag.ContinueOnError)
	help := flagSet.Bool("h", false, "Show command-line help")

	if err := flagSet.Parse(args); err != nil {
		usage(flagSet)
		return 1 // exit status
	}

	if *help {
		usage(flagSet)
		return 2 // exit status
	}

	return ps()
}

func usage(flagSet *flag.FlagSet) {
	common.ShowVersion()
	fmt.Println("ps [OPTION]...")
	flagSet.PrintDefaults()
}

func ps() int {

	pids, errPids := listPids()
	if errPids != nil {
		fmt.Printf("ps: %v\n", errPids)
		return 10
	}

	var status int

	fmt.Printf("%-5s %-6s\n", "PID", "NAME")

	for _, pid := range pids {

		info, errProc := procStatus(pid)
		if errProc != nil {
			fmt.Printf("ps: %v\n", errProc)
			status = 11
		}

		fmt.Printf("%5d %-6s\n", pid, info["Name"])
	}

	return status
}

func listPids() ([]int, error) {
	var pids []int

	f, errOpen := os.Open("/proc")
	if errOpen != nil {
		return pids, errOpen
	}
	defer f.Close()

	names, errRead := f.Readdirnames(0)
	if errRead != nil {
		return pids, errRead
	}

	for _, n := range names {
		pid, errConv := strconv.Atoi(n)
		if errConv == nil {
			pids = append(pids, pid)
		}
	}

	return pids, nil
}

func procStatus(pid int) (map[string]string, error) {
	path := fmt.Sprintf("/proc/%d/status", pid)
	f, errOpen := os.Open(path)
	if errOpen != nil {
		return nil, errOpen
	}
	defer f.Close()

	info := map[string]string{}

	input := bufio.NewReader(f)

LOOP:
	for {
		line, errInput := input.ReadString('\n')
		switch errInput {
		case nil:
		case io.EOF:
			break LOOP
		default:
			fmt.Printf("ps: %v\n", errInput)
			return info, errInput
		}

		pair := strings.SplitN(line, ":", 2)
		if len(pair) < 2 {
			fmt.Printf("ps: bad status line: %s", line)
			continue
		}

		key := strings.TrimSpace(pair[0])
		value := strings.TrimSpace(pair[1])

		info[key] = value
	}

	return info, nil
}
