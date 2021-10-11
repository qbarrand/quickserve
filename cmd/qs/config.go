package main

import "flag"

type CommandLine struct {
	Address string
	Paths   []string
	Version bool
}

func ParseCommandLine(programName string, args []string) (*CommandLine, error) {
	cl := CommandLine{}

	fs := flag.NewFlagSet(programName, flag.ContinueOnError)

	fs.StringVar(&cl.Address, "address", ":8080", "the address to bind to")
	fs.BoolVar(&cl.Version, "version", false, "displays the version and exits")

	err := fs.Parse(args)

	cl.Paths = fs.Args()

	return &cl, err
}
