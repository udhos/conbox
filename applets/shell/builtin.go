package shell

import (
	"fmt"
	"os"
	"sort"
)

type builtinFunc func(builtins map[string]builtinFunc, params []string) (bool, int)

func loadBuiltins() map[string]builtinFunc {
	tab := map[string]builtinFunc{
		"builtins": builtinBuiltins,
		"cd":       builtinCd,
		"exit":     builtinExit,
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

func builtinBuiltins(builtins map[string]builtinFunc, params []string) (bool, int) {
	listBuiltins(builtins, " ")
	fmt.Println()
	return false, 0
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
