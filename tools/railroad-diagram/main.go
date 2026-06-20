package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strconv"

	"nhatp.com/go/sugar"
)

//go:embed data
var embeddedAssets embed.FS

// OverlayFS tries the disk FS first, then falls back to the embedded FS.
type OverlayFS struct {
	primary   fs.FS
	secondary fs.FS
}

func (o OverlayFS) Open(name string) (fs.File, error) {
	f, err := o.primary.Open(name)
	if err == nil {
		return f, nil
	}
	return o.secondary.Open(name)
}

func main() {
	DefaultPort := strconv.Itoa(sugar.ToolRailroadDiagramPort)

	resourceFS, err := fs.Sub(embeddedAssets, "data/tabatkins")
	if err != nil {
		panic(err)
	}

	overlay := OverlayFS{
		primary:   os.DirFS("./overrides"), // put override files here
		secondary: resourceFS,
	}

	http.Handle("/", http.FileServer(http.FS(overlay)))

	fmt.Println("Railroad Diagram server starting at http://localhost:" + DefaultPort)

	err = http.ListenAndServe(":"+DefaultPort, nil)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
