package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"

	env "github.com/kreuzwerker/envplate"
	"github.com/yawn/doubledash"
)

var build string

func main() {

	flagArgs, execArgs := doubledash.Args, doubledash.Xtra
	os.Args = flagArgs

	var (
		help    = flag.Bool("h", false, "display usage")
		backup  = flag.Bool("b", false, "create a backup file when using inline mode")
		strict  = flag.Bool("s", false, "strict-mode - fail when falling back on defaults")
		dryRun  = flag.Bool("d", false, "dry-run - output templates to stdout instead of inline replacement")
		verbose = flag.Bool("v", false, "verbose logging")
	)

	flag.Parse()

	env.Config.Backup = *backup
	env.Config.DryRun = *dryRun
	env.Config.Strict = *strict
	env.Config.Verbose = *verbose

	if *help {
		fmt.Fprintf(os.Stderr, "Usage of %s (%s):\n", os.Args[0], build)
		flag.PrintDefaults()
	} else {
		env.Apply(flag.Args())
	}

	if len(execArgs) > 0 {

		if err := syscall.Exec(execArgs[0], execArgs, os.Environ()); err != nil {
			env.Log(env.ERROR, "Cannot exec '%v': %v", execArgs, err)
		}

	}

}
