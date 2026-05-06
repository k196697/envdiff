// Package config handles parsing and validation of CLI flags and options
// for the envdiff tool.
package config

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

// Format represents the output format for diff results.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Options holds all configuration parsed from CLI flags.
type Options struct {
	Files       []string
	Dir         string
	Prefix      string
	ExcludeKeys []string
	Format      Format
	Quiet       bool
}

// Parse reads arguments and returns a populated Options struct.
func Parse(args []string) (*Options, error) {
	fs := flag.NewFlagSet("envdiff", flag.ContinueOnError)

	dir := fs.String("dir", "", "directory containing .env files to compare")
	prefix := fs.String("prefix", "", "only include keys with this prefix")
	exclude := fs.String("exclude", "", "comma-separated list of keys to exclude")
	format := fs.String("format", "text", "output format: text or json")
	quiet := fs.Bool("quiet", false, "suppress output, only set exit code")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	opts := &Options{
		Files:  fs.Args(),
		Dir:    *dir,
		Prefix: *prefix,
		Quiet:  *quiet,
	}

	switch Format(*format) {
	case FormatText, FormatJSON:
		opts.Format = Format(*format)
	default:
		return nil, fmt.Errorf("unknown format %q: must be \"text\" or \"json\"", *format)
	}

	if *exclude != "" {
		for _, k := range strings.Split(*exclude, ",") {
			if k = strings.TrimSpace(k); k != "" {
				opts.ExcludeKeys = append(opts.ExcludeKeys, k)
			}
		}
	}

	if err := opts.validate(); err != nil {
		return nil, err
	}

	return opts, nil
}

func (o *Options) validate() error {
	if o.Dir == "" && len(o.Files) < 2 {
		return errors.New("provide at least two .env files or use --dir")
	}
	if o.Dir != "" && len(o.Files) > 0 {
		return errors.New("--dir and positional file arguments are mutually exclusive")
	}
	return nil
}
