package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
)

var (
	commit  string
	version string
)

func main() {
	cfg, err := ParseCommandLine(os.Args[0], os.Args[1:])
	if err != nil {
		log.Fatalf("Could not parse the command line: %v", err)
	}

	if cfg.Version {
		if err = writeVersion(os.Stdout); err != nil {
			log.Fatalf("Could not print the version information: %v", err)
		}

		return
	}

	sfs := make(map[string]bool, len(cfg.Paths))

	h := http.FileServer(subFS(sfs))

	if err = http.ListenAndServe(cfg.Address, h); err != nil {
		log.Fatalf("Error while running the server: %v", err)
	}
}

func writeVersion(w io.Writer) error {
	if _, err := fmt.Fprint(w, "Version: ", version); err != nil {
		return err
	}

	_, err := fmt.Fprint(w, "Commit: ", commit)

	return err
}

//type Handler struct {
//	logger *log.Logger
//	sfs    subFS
//}
//
//func NewHandler(logger *log.Logger, paths []string) *Handler {
//	h := Handler{
//		logger: logger,
//		sfs:    make(map[string]bool, len(paths)),
//	}
//
//	for _, p := range paths {
//		h.sfs[p] = true
//	}
//
//	return &h
//}
//
//func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	http.FileServer()
//}

type subFS map[string]bool

func (sfs subFS) Open(name string) (http.File, error) {

	fs.Sub()
	log.Printf("Looking for %q", name)

	if !sfs[name] {
		return nil, fs.ErrNotExist
	}

	return os.Open(name)
}
