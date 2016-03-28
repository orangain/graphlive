package main

import (
	"errors"
	"flag"
	"fmt"
)

// See: https://lawlessguy.wordpress.com/2013/07/23/filling-a-slice-using-command-line-flags-in-go-golang/
type stringslice []string

func (s *stringslice) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *stringslice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

type Opts struct {
	commands stringslice
	port     int
}

func ParseOpts() (result Opts, err error) {
	var opts Opts
	var help bool

	flag.Var(&opts.commands, "c", "")
	flag.BoolVar(&help, "h", false, "")
	flag.BoolVar(&help, "help", false, "")
	flag.IntVar(&opts.port, "p", 9999, "")
	flag.IntVar(&opts.port, "port", 9999, "")
	flag.Parse()

	if help || len(opts.commands) == 0 {
		return opts, errors.New(
			`Usage:
    graphlive [options] -c COMMAND [-c COMMAND]...

Options:
    -p PORT, --port PORT
        Port to listen on (default 9999)
    -h, --help
        Show this help message`)
	} else {
		return opts, nil
	}
}
