package utils

import (
	"fmt"
	"log"
	"net/http"
)

func HttpServeDir(path string, url string, done chan bool) {
	fmt.Printf("Serving directory %s at url: %s\n", path, url)

	err := http.ListenAndServe(url, http.FileServer(http.Dir(path)))
	if err != nil {
		log.Fatal(err)
	}

	<-done

	fmt.Println("Stopped serving directory")
}
