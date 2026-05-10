package main

import (
	"fmt"
	"io/fs"
	"net/http"
	"strconv"

	"nhatp.com/go/sugar"
	"nhatp.com/go/sugar/tools/lexeme-viewer/asset"
)

func main() {
	DefaultPort := strconv.Itoa(sugar.ToolLexemeViewerDefaultPort)

	resourceFS, err := fs.Sub(asset.Content, "resource")
	if err != nil {
		panic(err)
	}

	fileServer := http.FileServer(http.FS(resourceFS))

	http.Handle("/", fileServer)

	fmt.Println("Lexeme Viewer server starting at http://localhost:" + DefaultPort)

	err = http.ListenAndServe(":"+DefaultPort, nil)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
