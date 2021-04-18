package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func devServe() {
	var buildPath, publicPath string

	if len(os.Args) != 4 {
		_, _ = fmt.Scanln(&buildPath)
		_, _ = fmt.Scanln(&publicPath)
	} else {
		buildPath = os.Args[2]
		publicPath = os.Args[3]
	}

	log.Println("Starting server")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op & fsnotify.Write == fsnotify.Write || event.Op & fsnotify.Create == fsnotify.Create {
					log.Println("modified file:", event.Name)
					handleBuildFolderModified(buildPath, publicPath)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(buildPath)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Watcher created")

	go func() {
		log.Fatal(http.ListenAndServe(`:8080`, http.FileServer(http.Dir(publicPath))))
	}()

	log.Println("Listening...")

	fmt.Scanln()
}

func handleBuildFolderModified(buildPath, publicPath string) {
	log.Println("Build directory updated")
	log.Println("Copying to the public directory")
	files, err := listFiles(buildPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		p1 := buildPath + "/" + f
		p2 := publicPath + "/" + f
		_, _ = copyFile(p1, p2)
	}
	log.Println("Copied")
}

func listFiles(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(files))
	for _, f := range files {
		if !f.IsDir() {
			paths = append(paths, f.Name())
		}
	}

	return paths, nil
}

func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
