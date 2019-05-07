package sh

import (
	//"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	//"os/exec"
	//"path/filepath"

	"github.com/udhos/conbox/common"
	"golang.org/x/crypto/ssh/terminal"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

type shell struct {
	parser           *syntax.Parser
	mainRunner       *interp.Runner
	command          string
	forceInteractive bool
	flagSet          *flag.FlagSet
}

// Run executes the applet.
func Run(tab map[string]common.AppletFunc, args []string) int {

	var help, interactive bool
	var cmd string

	flagSet := flag.NewFlagSet("sh", flag.ContinueOnError)
	flagSet.BoolVar(&help, "h", false, "Show command-line help")
	flagSet.BoolVar(&interactive, "i", false, "Force interactive mode")
	flagSet.StringVar(&cmd, "c", "", "Command to be executed")

	if err := flagSet.Parse(args); err != nil {
		usage(flagSet)
		return 1 // exit status
	}

	if help {
		usage(flagSet)
		return 2 // exit status
	}

	// create shell
	runner, errInterp := interp.New(interp.StdIO(os.Stdin, os.Stdout, os.Stderr))
	if errInterp != nil {
		fmt.Printf("sh: new interp: %v\n", errInterp)
		return 3
	}
	s := shell{
		parser:           syntax.NewParser(),
		mainRunner:       runner,
		command:          cmd,
		forceInteractive: interactive,
		flagSet:          flagSet,
	}

	switch err := runAll(s).(type) {
	case nil:
	case interp.ShellExitStatus:
		return int(err)
	case interp.ExitStatus:
		return int(err)
	default:
		fmt.Fprintln(os.Stderr, err)
		return 4
	}

	return 0
}

func usage(flagSet *flag.FlagSet) {
	common.ShowVersion()
	fmt.Println("sh [OPTION]... [FILE]")
	flagSet.PrintDefaults()
}

func runAll(s shell) error {
	if s.command != "" {
		r := strings.NewReader(s.command)
		if s.forceInteractive {
			return interactive(s, r)
		}
		return run(s, r, "")
	}
	if s.flagSet.NArg() == 0 {
		if s.forceInteractive {
			return interactive(s, s.mainRunner.Stdin)
		}
		if terminal.IsTerminal(int(os.Stdin.Fd())) {
			return interactive(s, s.mainRunner.Stdin)
		}
		return run(s, os.Stdin, "")
	}
	for _, path := range s.flagSet.Args() {
		if err := runPath(s, path); err != nil {
			return err
		}
	}
	return nil
}

func runPath(s shell, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	if s.forceInteractive {
		return interactive(s, f)
	}
	return run(s, f, path)
}

func run(s shell, reader io.Reader, name string) error {
	prog, err := s.parser.Parse(reader, name)
	if err != nil {
		return err
	}
	s.mainRunner.Reset()
	ctx := context.Background()
	return s.mainRunner.Run(ctx, prog)
}

func interactive(s shell, reader io.Reader) error {
	fmt.Fprintf(s.mainRunner.Stdout, "$ ")

	var errStmt error

	fn := func(stmts []*syntax.Stmt) bool {
		if s.parser.Incomplete() {
			fmt.Fprintf(s.mainRunner.Stdout, "> ")
			return true
		}
		ctx := context.Background()
		for _, stmt := range stmts {
			switch err := s.mainRunner.Run(ctx, stmt).(type) {
			case nil:
			case interp.ShellExitStatus:
				//os.Exit(int(err))
				errStmt = err
				return false
			case interp.ExitStatus:
			default:
				fmt.Fprintln(s.mainRunner.Stderr, err)
				//os.Exit(1)
				errStmt = err
				return false
			}
		}
		fmt.Fprintf(s.mainRunner.Stdout, "$ ")
		return true
	}

	errInter := s.parser.Interactive(reader, fn)
	if errInter == nil {
		return errStmt
	}
	return errInter
}
