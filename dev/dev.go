package main

import (
	"fmt"
	"html/template"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/ishchenko-gv/ishchenko-dev/tmpl"
)

func main() {
	go startLivereload()

	log.Println("Server is running on http://localhost:3000")
	err := http.ListenAndServe(":3000", http.HandlerFunc(handler))
	if err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	reqPath := r.URL.Path
	file, mimeType, err := loadFile(reqPath)
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if strings.HasPrefix(mimeType, "text/html") {
		pageData, err := getPageData()
		if err != nil {
			log.Printf("%v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		file, err = tmpl.HandleHtmlTemplate(file, pageData)
		if err != nil {
			log.Printf("%v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.Header().Add("Content-Type", mimeType)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(file)
	if err != nil {
		log.Printf("%v", err)
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

	log.Printf("Loading file %s - extension %s - mime type %s", filePath, ext, mimeType)

	return f, mimeType, nil
}
func getPageData() (*tmpl.PageData, error) {
	livereloadScriptSrc, err := os.ReadFile("dev/livereload.js")
	if err != nil {
		return nil, err
	}

	livereloadScriptContent := fmt.Sprintf("<script>%s</script>", string(livereloadScriptSrc))

	return &tmpl.PageData{
		LivereloadScript: template.HTML(livereloadScriptContent),
	}, nil
}
