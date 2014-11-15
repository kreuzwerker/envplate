package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"

	env "github.com/yawn/envplate"
)

var build string

func split() (flags []string, execArgs []string) {

	split := len(os.Args)

	for idx, e := range os.Args {

		if e == "--" {
			split = idx
			break
		}

	}

	flags = os.Args[1:split]

	if split < len(os.Args) {
		execArgs = os.Args[split+1 : len(os.Args)]
	}

	return flags, execArgs

}

func main() {

	flags, execArgs := split()

	os.Args = flags

	help := flag.Bool("h", false, "display usage")

	e := env.Envplate{
		Backup:  flag.Bool("b", false, "create a backup file when using inline mode"),
		Debug:   flag.Bool("d", false, "output templates to stdout instead of inline replacement"),
		Verbose: flag.Bool("v", false, "verbose logging"),
	}

	flag.Parse()

	if *help {
		fmt.Fprintf(os.Stderr, "Usage of %s (%s):\n", os.Args[0], build)
		flag.PrintDefaults()
	} else {
		e.Apply(flag.Args())
	}

	if len(execArgs) > 0 {

		if err := syscall.Exec(execArgs[0], execArgs, os.Environ()); err != nil {
			e.Log(env.ERROR, "Cannot exec '%v': %v", execArgs, err)
		}

	}

}
