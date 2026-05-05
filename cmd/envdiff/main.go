package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/reporter"
)

const usage = `envdiff - Compare .env files across environments

Usage:
  envdiff [flags] <file1> <file2>

Flags:
`

func main() {
	var (
		quiet   = flag.Bool("quiet", false, "suppress output, exit code only")
		noColor = flag.Bool("no-color", false, "disable colored output")
	)

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
		flag.PrintDefaults()
	}

	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		os.Exit(2)
	}

	fileA, fileB := args[0], args[1]

	envA, err := parser.ParseFile(fileA)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", fileA, err)
		os.Exit(1)
	}

	envB, err := parser.ParseFile(fileB)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", fileB, err)
		os.Exit(1)
	}

	results := diff.Compare(envA, envB)
	summary := diff.Summary(results)

	if !*quiet {
		reporter.Report(os.Stdout, results, fileA, fileB, reporter.Options{
			NoColor: *noColor,
		})
		fmt.Fprintf(os.Stdout, "\nSummary: %d matched, %d mismatched, %d missing\n",
			summary.Matched, summary.Mismatched, summary.Missing)
	}

	if summary.Mismatched > 0 || summary.Missing > 0 {
		os.Exit(1)
	}
}
