package main

import "flag"

type CommandLine struct {
	Address       string
	AllowDotFiles bool
	Paths         []string
	Version       bool
}

func ParseCommandLine(programName string, args []string) (*CommandLine, error) {
	cl := CommandLine{}

	fs := flag.NewFlagSet(programName, flag.ContinueOnError)

	fs.StringVar(&cl.Address, "address", ":8080", "the address to bind to")
	fs.BoolVar(&cl.AllowDotFiles, "allow-dotfiles", false, "allow access to resources with a dot-leading name (e.g. .htpasswd)")
	fs.BoolVar(&cl.Version, "version", false, "displays the version and exits")

	err := fs.Parse(args)

	cl.Paths = fs.Args()

	return &cl, err
}
