package main

import (
	"flag"
	"fmt"
	"os"

	env "github.com/yawn/envplate"
)

var build string

func main() {

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

}
