package main

import (
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/kreuzwerker/envplate"
	"github.com/spf13/cobra"
	"github.com/yawn/doubledash"
)

var (
	build   string
	version string
)

func main() {

	var ( // flags
		backup  *bool
		dryRun  *bool
		strict  *bool
		verbose *bool
	)

	os.Args = doubledash.Args

	// commands
	root := &cobra.Command{
		Use:   "ep",
		Short: "envplate provides trivial templating for configuration files using environment keys",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			envplate.Logger.Verbose = *verbose
		},
	}

	parse := &cobra.Command{
		Use:   "parse",
		Short: "Parse globs and exec after doubledash",
		Run: func(cmd *cobra.Command, args []string) {

			var h = envplate.Handler{
				Backup: *backup,
				DryRun: *dryRun,
				Strict: *strict,
			}

			if err := h.Apply(args); err != nil {
				os.Exit(1)
			}

			if h.DryRun {
				os.Exit(0)
			}

			if len(doubledash.Xtra) > 0 {

				if err := syscall.Exec(doubledash.Xtra[0], doubledash.Xtra, os.Environ()); err != nil {
					log.Fatalf("Cannot exec '%v': %v", doubledash.Xtra, err)
				}

			}

		},
	}

	version := &cobra.Command{
		Use:   "version",
		Short: "Print the version information of ep",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Envplate %s (%s)\n", version, build)
		},
	}

	root.AddCommand(parse)
	root.AddCommand(version)

	// flag parsing
	backup = parse.Flags().BoolP("backup", "b", false, "Create a backup file when using inline mode")
	dryRun = parse.Flags().BoolP("dry-run", "d", false, "Dry-run - output templates to stdout instead of inline replacement")
	strict = parse.Flags().BoolP("strict", "s", false, "Strict-mode - fail when falling back on defaults")
	verbose = parse.Flags().BoolP("verbose", "v", false, "Verbose logging")

	if err := root.Execute(); err != nil {
		log.Fatalf("Failed to start the application: %v", err)
	}

}
