package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"regexp"
)

func main() {
	log("Server is running on http://localhost:3000\n")
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
	_, err = w.Write(file)
	if err != nil {
		log("error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
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
	var fileExt = regexp.MustCompile(`\.[A-Za-z0-9]{3,4}`)
	ext := fileExt.FindString(filePath)

	log("filePath %s, fileExt %s", filePath, ext)

	filePath = path.Join("src", filePath)
	if ext == "" {
		filePath = path.Join(filePath, "index.html")
	}

	f, err := os.ReadFile(filePath)
	if err != nil {
		return nil, "", err
	}

	mime := mimeFromExt(ext)

	return f, mime, nil
}

func mimeFromExt(fileExt string) MimeType {
	switch fileExt {
	case "":
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
	fmt.Printf(message+"\n", args...)
}
