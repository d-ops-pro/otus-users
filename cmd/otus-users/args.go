package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
)

func validateFS(fs *flag.FlagSet, httpAddr, dbURI string) {
	if httpAddr == "" {
		fs.Usage()
		panic("http-addr is not set")
	}

	if dbURI == "" {
		fs.Usage()
		panic("db-uri is not set")
	}
}

func usageFor(fs *flag.FlagSet, short string) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "USAGE\n")
		fmt.Fprintf(os.Stderr, "  %s\n", short)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "FLAGS\n")
		w := tabwriter.NewWriter(os.Stderr, 0, 2, 2, ' ', 0)
		fs.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(w, "\t-%s %s\t%s\n", f.Name, f.DefValue, f.Usage)
		})
		w.Flush()
		fmt.Fprintf(os.Stderr, "\n")
	}
}
