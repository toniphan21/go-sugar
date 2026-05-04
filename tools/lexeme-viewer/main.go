package main

import (
	"fmt"
	"io/fs"
	"net/http"

	"nhatp.com/go/sugar/tools/lexeme-viewer/asset"
)

const DefaultPort = "39800"

func main() {
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
