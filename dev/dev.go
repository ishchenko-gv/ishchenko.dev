package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"regexp"
	"time"
)

func main() {
	log("Server is running on http://localhost:3000")
	http.ListenAndServe(":3000", http.HandlerFunc(handler))
}

func handler(w http.ResponseWriter, r *http.Request) {
	reqPath := r.URL.Path
	file, mime, err := loadFile(reqPath)
	if err != nil {
		log("error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", string(mime))
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(file)
	if err != nil {
		log("error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type MimeType string

var (
	MimeTypeHtml  MimeType = "text/html"
	MimeTypeCss   MimeType = "text/css"
	MimeTypeJs    MimeType = "text/javascript"
	MimeTypeJson  MimeType = "application/json"
	MimeTypeXIcon MimeType = "image/x-icon"
	MimeTypeSvg   MimeType = "image/svg+xml"
	MimeTypePng   MimeType = "image/png"
)

func loadFile(filePath string) ([]byte, MimeType, error) {
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

	fileExt := extFromPath(filePath)
	mime := mimeFromExt(fileExt)

	log("Loading file %s - extension %s - mime type %s", filePath, fileExt, mime)

	return f, mime, nil
}

func extFromPath(filePath string) string {
	var fileExt = regexp.MustCompile(`\.[A-Za-z0-9]{3,4}`)
	return fileExt.FindString(filePath)
}

func mimeFromExt(fileExt string) MimeType {
	switch fileExt {
	case ".html":
		return MimeTypeHtml
	case ".css":
		return MimeTypeCss
	case ".js":
		return MimeTypeJs
	case ".json":
		return MimeTypeJson
	case ".ico":
		return MimeTypeXIcon
	case ".svg":
		return MimeTypeSvg
	case ".png":
		return MimeTypePng
	default:
		panic(fmt.Errorf("unsupported file extension: %s", fileExt))
	}
}

func log(message string, args ...any) {
	now := time.Now().Format(time.RFC822)
	msg := fmt.Sprintf(message+"\n", args...)
	fmt.Printf("[%s] %s", now, msg)
}
