package utils

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func HttpServeDir(path string, url string, done chan bool) {
	fmt.Printf("Serving directory %s at url: %s\n", path, url)

	srv := &http.Server{Addr: url}
	srv.Handler = &MyHandler{dir: path}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()

	<-done

	_ = srv.Shutdown(context.Background())

	fmt.Println("Stopped serving directory")
}

type MyHandler struct {
	dir string
}

func (m *MyHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var path = request.URL.Path
	if path == "" || path == "/" {
		m.serveIndex(writer, request)
		return
	}

	requestedFilePath := filepath.Join(m.dir, path)
	_, err := os.Open(requestedFilePath)
	if err != nil {
		//could not open the requested file, try index.html
		m.serveIndex(writer, request)
		return
	}

	http.ServeFile(writer, request, requestedFilePath)
}

func (m *MyHandler) serveIndex(writer http.ResponseWriter, request *http.Request) {
	indexPath := filepath.Join(m.dir, "index.html")
	_, err := os.Open(indexPath)
	if err != nil {
		http.NotFound(writer, request)
		return
	}

	http.ServeFile(writer, request, indexPath)
}