package main

import (
	"log"
	"net/http"

	"github.com/fsnotify/fsnotify"
)

func startLivereload() {
	fsChan := make(chan bool)

	go watchFsForLivereload(fsChan)
	go serveLivereloadSse(fsChan)
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
				log.Println("event:", event)
				if event.Has(fsnotify.Write) {
					log.Println("modified file:", event.Name)
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

	// Block main goroutine forever.
	<-make(chan struct{})
}

func serveLivereloadSse(fsChan <-chan bool) {
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
				return
			}
		}
	}

	http.ListenAndServe(":3001", http.HandlerFunc(handler))
}
