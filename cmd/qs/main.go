package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/qbarrand/quickserve/pkg/middlewares"
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

	filesystem, err := fsFromPaths(cfg.Paths...)
	if err != nil {
		log.Fatalf("Could not create a filesystem to serve: %v", err)
	}

	h := http.FileServer(http.FS(filesystem))

	if !cfg.AllowDotFiles {
		h = middlewares.HideDotFiles(http.StatusForbidden, h)
	}

	if err = http.ListenAndServe(cfg.Address, h); err != nil {
		log.Fatalf("Error while running the server: %v", err)
	}
}

func writeVersion(w io.Writer) error {
	if _, err := fmt.Fprintln(w, "Version: ", version); err != nil {
		return err
	}

	_, err := fmt.Fprintln(w, "Commit: ", commit)

	return err
}

type aliasDirEntry struct {
	de      fs.DirEntry
	sysPath string
}

type dirFileInfo struct{}

func (dfi *dirFileInfo) Name() string       { return "/" }
func (dfi *dirFileInfo) Size() int64        { return 0 }
func (dfi *dirFileInfo) Mode() os.FileMode  { return os.ModeDir }
func (dfi *dirFileInfo) ModTime() time.Time { return time.Now() }
func (dfi *dirFileInfo) IsDir() bool        { return true }
func (dfi *dirFileInfo) Sys() interface{}   { return nil }

// rootFile implements both fs.FS and fs.ReadDirFile.
type rootFile map[string]*aliasDirEntry

func (rf rootFile) Close() error { return nil }

func (rf rootFile) Read(b []byte) (int, error) {
	return 0, errors.New("not implemented")
}

func (rf rootFile) Stat() (fs.FileInfo, error) {
	return &dirFileInfo{}, nil
}

func (rf rootFile) Open(name string) (fs.File, error) {
	log.Printf("open(%q)", name)

	if name == "/" || name == "." {
		return rf, nil
	}

	if ade, ok := rf[name]; ok {
		return os.Open(ade.sysPath)
	}

	elems := strings.SplitN(name, "/", 2)

	// Check the first segment in name. Is it a directory name owned by sf?
	if len(elems) == 2 {
		dir := elems[0]

		dirEntry, ok := rf[dir]
		if !ok {
			return nil, os.ErrNotExist
		}

		return os.DirFS(dirEntry.sysPath).Open(elems[1])
	}

	return nil, os.ErrNotExist
}

func (rf rootFile) ReadDir(n int) ([]fs.DirEntry, error) {
	if n <= 0 || len(rf) < n {
		n = len(rf)
	}

	entries := make([]fs.DirEntry, 0, n)

	for _, ade := range rf {
		// append up to n values
		if len(entries) >= n {
			break
		}

		entries = append(entries, ade.de)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	return entries, nil
}

func (rf rootFile) addMapping(name, fullPath string) error {
	if ade := rf[name]; ade != nil {
		return fmt.Errorf("%s already maps to %s", name, ade.sysPath)
	}

	fi, err := os.Stat(fullPath)
	if err != nil {
		return fmt.Errorf("could not stat(%q): %v", fullPath, err)
	}

	rf[name] = &aliasDirEntry{
		de:      fs.FileInfoToDirEntry(fi),
		sysPath: fullPath,
	}

	return nil
}

func fsFromPaths(paths ...string) (fs.FS, error) {
	if len(paths) <= 0 {
		return nil, errors.New("at least one path is required")
	}

	sf := make(rootFile)

	for _, path := range paths {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return nil, fmt.Errorf("could not get the absolute path for %q: %v", path, err)
		}

		fi, err := os.Stat(absPath)
		if err != nil {
			return nil, fmt.Errorf("could not stat(%q): %v", absPath, err)
		}

		name := fi.Name()

		log.Printf("Adding %s=%s", name, absPath)

		if err := sf.addMapping(name, absPath); err != nil {
			return nil, fmt.Errorf("could not add mapping %s=%s: %v", name, absPath, err)
		}
	}

	return sf, nil
}
