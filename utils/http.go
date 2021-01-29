package utils

import (
	"context"
	"fmt"
	"net/http"
)

func HttpServeDir(path string, url string, done chan bool) {
	fmt.Printf("Serving directory %s at url: %s\n", path, url)

	srv := &http.Server{Addr: url}
	srv.Handler = http.FileServer(http.Dir(path))

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()

	<-done

	_ = srv.Shutdown(context.Background())

	fmt.Println("Stopped serving directory")
}
