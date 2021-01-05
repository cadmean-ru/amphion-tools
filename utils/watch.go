package utils

import (
	"github.com/fsnotify/fsnotify"
	"log"
)

func WatchDirs(dirs []string, done chan bool, onUpdate func()) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	for _, dir := range dirs {
		err = watcher.Add(dir)
		if err != nil {
			log.Fatal(err)
		}
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			log.Println("event:", event)
			if event.Op & fsnotify.Write == fsnotify.Write || event.Op & fsnotify.Create == fsnotify.Create {
				log.Println("modified file:", event.Name)
				onUpdate()
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		case <-done:
			return
		}
	}
}
