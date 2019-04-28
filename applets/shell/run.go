package shell

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

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

		// read input from stdin

		stdin := os.Stdin
		if !*interactive {
			// check if stdin is terminal (interactive mode)
			info, errStat := stdin.Stat()
			if errStat != nil {
				fmt.Printf("shell: stat stdin: %v\n", errStat)
			} else {
				if (info.Mode() & os.ModeCharDevice) != 0 {
					*interactive = true
				}
			}
		}
		return loop(tab, stdin, *interactive)
	}

	// read input from file

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

	builtins := loadBuiltins()

	if interactive {
		common.ShowVersion()
		fmt.Print(`
welcome to conbox shell.
this tiny shell is very limited in features.
however you can run external programs normally.
some hints:
       - use 'conbox' to see all applets available as shell commands.
       - use 'help' to list shell built-in commands.
       - 'exit' terminates the shell.

`)

		builtinHelp(builtins, nil)

		fmt.Println()
	}

	var prompt string
	euid := os.Geteuid()
	if euid == 0 {
		prompt = "#"
	} else {
		prompt = "$"
	}

	input := bufio.NewReader(r)

LOOP:
	for {
		if interactive {
			fmt.Printf("conbox shell%s ", prompt)
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

		parameters := common.Tokenize(line)

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

		exit, status := execute(tab, builtins, parameters)
		if exit {
			return status
		}
	}

	return 0
}

func expand(params []string) []string {

	// 1. expand env vars

	//fmt.Printf("== %d %v\n", len(params), params)

	var exp1 []string
	for _, p := range params {
		exp1 = append(exp1, os.ExpandEnv(p))
	}

	//fmt.Printf("$$ %d %v\n", len(exp1), exp1)

	// 2. expand file globs

	var exp2 []string
	for _, p := range exp1 {
		if p == "" {
			exp2 = append(exp2, p)
			continue
		}
		g, errGlob := filepath.Glob(p)
		if errGlob != nil {
			fmt.Printf("shell: glob: %v", errGlob)
			exp2 = append(exp2, p)
			continue
		}
		if len(g) < 1 {
			exp2 = append(exp2, p)
			continue
		}
		exp2 = append(exp2, g...)
	}

	//fmt.Printf("** %d %v\n", len(exp2), exp2)

	return exp2
}

func execute(tab map[string]common.AppletFunc, builtins map[string]builtinFunc, params []string) (bool, int) {

	exp := expand(params)

	prog := exp[0]

	// 1. lookup shell built-in

	if b, found := builtins[prog]; found {
		return b(builtins, exp[1:])
	}

	// 2. lookup conbox applet

	if applet, found := tab[prog]; found {
		return false, applet(tab, exp[1:])
	}

	// 3. call external program

	return false, external(exp)
}

func external(params []string) int {

	c := exec.Command(params[0], params[1:]...)

	c.Stdout = os.Stdout
	c.Stdin = os.Stdin
	c.Stderr = os.Stderr

	if errStart := c.Start(); errStart != nil {
		fmt.Printf("shell: %v\n", errStart)
		return 30
	}

	if errWait := c.Wait(); errWait != nil {
		if err, isExit := errWait.(*exec.ExitError); isExit {
			return err.ExitCode()
		}
		fmt.Printf("shell: uexpected exit error: %v\n", errWait)
		return 31
	}

	return 0
}
