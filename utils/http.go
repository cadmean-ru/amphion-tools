package utils

import (
	"log"
	"net/http"
)

func HttpServeDir(path string, url string, done chan bool) {
	err := http.ListenAndServe(url, http.FileServer(http.Dir(path)))
	if err != nil {
		log.Fatal(err)
	}

	<-done
}
