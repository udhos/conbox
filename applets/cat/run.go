package cat

import (
	"flag"
	"fmt"
	"io"
	"os"
)

// Run executes the applet.
func Run(args []string) int {

	flagSet := flag.NewFlagSet("cat", flag.ExitOnError)
	helpFlag := flagSet.Bool("h", false, "Show command-line help")

	if err := flagSet.Parse(args); err != nil {
		return 1 // exit status
	}

	if *helpFlag {
		fmt.Println("cat [OPTIONS]... [FILES]...")
		flagSet.PrintDefaults()
		return 2 // exit status
	}

	if flagSet.NArg() < 1 {
		if err := catReader(os.Stdin); err != nil {
			fmt.Printf("cat: %v\n", err)
			return 4
		}
		return 0
	}

	var status int

	for _, f := range flagSet.Args() {
		if err := catFile(f); err != nil {
			fmt.Printf("cat: %v\n", err)
			status = 5
		}
	}

	return status // exit status
}

func catFile(path string) error {
	f, errOpen := os.Open(path)
	if errOpen != nil {
		return errOpen
	}

	err := catReader(f)

	f.Close()

	return err
}

func catReader(r io.Reader) error {
	_, err := io.Copy(os.Stdout, r)
	return err
}
