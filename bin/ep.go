package main

import (
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/kreuzwerker/envplate"
	_ "github.com/paulrosania/go-charset/data"
	"github.com/spf13/cobra"
	"github.com/yawn/doubledash"
)

var (
	build   string
	version string
)

func init() {
	os.Args = doubledash.Args
}

func main() {

	var ( // flags
		backup  *bool
		dryRun  *bool
		strict  *bool
		verbose *bool
		charset *string
	)

	root := &cobra.Command{
		Use:   "ep",
		Short: fmt.Sprintf("envplate %s (%s) provides trivial templating for configuration files using environment keys", version, build),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			envplate.Logger.Verbose = *verbose
		},
		Run: func(cmd *cobra.Command, args []string) {

			var h = envplate.Handler{
				Backup:  *backup,
				DryRun:  *dryRun,
				Strict:  *strict,
				Charset: *charset,
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

	// flag parsing
	backup = root.Flags().BoolP("backup", "b", false, "Create a backup file when using inline mode")
	dryRun = root.Flags().BoolP("dry-run", "d", false, "Dry-run - output templates to stdout instead of inline replacement")
	strict = root.Flags().BoolP("strict", "s", false, "Strict-mode - fail when falling back on defaults")
	verbose = root.Flags().BoolP("verbose", "v", false, "Verbose logging")
	charset = root.Flags().StringP("charset", "c", "", "Output charset")

	if err := root.Execute(); err != nil {
		log.Fatalf("Failed to start the application: %v", err)
	}

}
