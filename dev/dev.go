package main

import (
	"bytes"
	"fmt"
	"html/template"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	log("Server is running on http://localhost:3000")
	err := http.ListenAndServe(":3000", http.HandlerFunc(handler))
	if err != nil {
		log("%v", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path
	file, mimeType, err := loadFile(filePath)
	if err != nil {
		log("error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if strings.HasPrefix(mimeType, "text/html") {
		file, err = handleHtmlTemplate(file)
		if err != nil {
			log("error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.Header().Add("Content-Type", mimeType)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(file)
	if err != nil {
		log("error: %v", err)
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

	log("Loading file %s - extension %s - mime type %s", filePath, ext, mimeType)

	return f, mimeType, nil
}

type PageData struct {
	LivereloadScript template.HTML
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

func log(message string, args ...any) {
	now := time.Now().Format(time.RFC822)
	msg := fmt.Sprintf(message+"\n", args...)
	fmt.Printf("[%s] %s", now, msg)
}
