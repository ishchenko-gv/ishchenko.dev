package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

func main() {
	go listenLivereload()

	logMsg("Server is running on http://localhost:3000")
	err := http.ListenAndServe(":3000", http.HandlerFunc(handler))
	if err != nil {
		logMsg("%v", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	reqPath := r.URL.Path
	file, mimeType, err := loadFile(reqPath)
	if err != nil {
		logMsg("error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if strings.HasPrefix(mimeType, "text/html") {
		file, err = handleHtmlTemplate(file)
		if err != nil {
			logMsg("error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.Header().Add("Content-Type", mimeType)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(file)
	if err != nil {
		logMsg("error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func loadFile(filePath string) ([]byte, string, error) {
	filePath = path.Join("src", filePath)

	fi, err := os.Stat(filePath)
	if err != nil {
		return nil, "", err
	}

	switch mode := fi.Mode(); {
	case mode.IsDir():
		filePath = path.Join(filePath, "index.html")
	}

	f, err := os.ReadFile(filePath)
	if err != nil {
		return nil, "", err
	}

	ext := filepath.Ext(filePath)
	mimeType := mime.TypeByExtension(ext)

	logMsg("Loading file %s - extension %s - mime type %s", filePath, ext, mimeType)

	return f, mimeType, nil
}

func handleHtmlTemplate(file []byte) ([]byte, error) {
	tmpl, err := template.New("page").Parse(string(file))
	if err != nil {
		return nil, err
	}

	pageData, err := getPageData()
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	err = tmpl.Execute(buf, pageData)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type PageData struct {
	LivereloadScript template.HTML
}

func getPageData() (*PageData, error) {
	livereloadScriptSrc, err := os.ReadFile("dev/livereload.js")
	if err != nil {
		return nil, err
	}

	livereloadScriptContent := fmt.Sprintf("<script>%s</script>", string(livereloadScriptSrc))

	return &PageData{
		LivereloadScript: template.HTML(livereloadScriptContent),
	}, nil
}

func logMsg(message string, args ...any) {
	now := time.Now().Format(time.RFC822)
	msg := fmt.Sprintf(message+"\n", args...)
	fmt.Printf("[%s] %s", now, msg)
}

func listenLivereload() {
	fsChan := make(chan bool)

	go listenLivereloadFs(fsChan)
	go listenLivereloadSse(fsChan)
}

func listenLivereloadFs(fsChan chan<- bool) {
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

func listenLivereloadSse(fsChan <-chan bool) {
	var liverealoadSseHandler = func(w http.ResponseWriter, r *http.Request) {
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

	http.ListenAndServe(":3001", http.HandlerFunc(liverealoadSseHandler))
}
