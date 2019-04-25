package shell

import (
	"fmt"
	"os"
	"sort"
)

type builtinFunc func(builtins map[string]builtinFunc, params []string) (bool, int)

func loadBuiltins() map[string]builtinFunc {
	tab := map[string]builtinFunc{
		"builtin": builtinBuiltin,
		"cd":      builtinCd,
		"exit":    builtinExit,
		"help":    builtinHelp,
	}
	return tab
}

func listBuiltins(builtins map[string]builtinFunc, sep string) {
	var list []string
	for n := range builtins {
		list = append(list, n)
	}
	sort.Strings(list)
	for _, n := range list {
		fmt.Printf("%s%s", n, sep)
	}
}

func builtinHelp(builtins map[string]builtinFunc, params []string) (bool, int) {
	listBuiltins(builtins, " ")
	fmt.Println()
	return false, 0
}

func builtinBuiltin(builtins map[string]builtinFunc, params []string) (bool, int) {
	if len(params) < 1 {
		return builtinHelp(builtins, params)
	}

	prog := params[0]

	if b, found := builtins[prog]; found {
		return b(builtins, params[1:])
	}

	fmt.Printf("cd: %s: not a shell builtin\n", prog)

	return false, 1
}

func builtinExit(builtins map[string]builtinFunc, params []string) (bool, int) {
	return true, 0
}

func builtinCd(builtins map[string]builtinFunc, params []string) (bool, int) {
	if len(params) < 1 {
		fmt.Println("cd: missing directory")
		return false, 1
	}
	if err := os.Chdir(params[0]); err != nil {
		fmt.Printf("cd: %v\n", err)
		return false, 2
	}
	return false, 0
}
