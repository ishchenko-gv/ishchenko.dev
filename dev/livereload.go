package main

import (
	"io/fs"
	"log"
	"net/http"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func startLivereload() {
	fsChan := make(chan bool)

	go watchFsForLivereload(fsChan)
	go serveLivereloadEventStream(fsChan)
}

func watchFsForLivereload(fsChan chan<- bool) {
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
				log.Println("file system event:", event)
				if event.Has(fsnotify.Write) {
					fsChan <- true
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("src")
	if err != nil {
		log.Fatal(err)
	}

	err = filepath.WalkDir("src", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			watcher.Add(path)
		}
		return nil
	})

	// Block main goroutine forever.
	<-make(chan struct{})
}

func serveLivereloadEventStream(fsChan <-chan bool) {
	var handler = func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/livereload" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		flusher := w.(http.Flusher)
		for {
			select {
			case <-fsChan:
				w.Write([]byte("data: fsChange\n\n"))
				flusher.Flush()
			case <-r.Context().Done():
				log.Println("Event stream closed")
				return
			}
		}
	}

	http.ListenAndServe(":3001", http.HandlerFunc(handler))
}
