package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"regexp"
)

func main() {
	fmt.Printf("Server is running on http://localhost:3000\n")
	http.ListenAndServe(":3000", http.HandlerFunc(handler))
}

func handler(w http.ResponseWriter, r *http.Request) {
	reqPath := r.URL.Path
	file, mime, err := loadFile(reqPath)

	if err != nil {
		fmt.Printf("error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", string(mime)+"; charset=utf-8")
	_, err = w.Write(file)
	if err != nil {
		fmt.Printf("error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type MimeType string

var (
	htmlMimeType MimeType = "text/html"
	cssMimeType  MimeType = "text/css"
	jsMimeType   MimeType = "text/javascript"
	jsonMimeType MimeType = "application/json"
)

func loadFile(filePath string) ([]byte, MimeType, error) {
	var fileExt = regexp.MustCompile(`\.[A-Za-z0-9]{3,4}`)
	ext := fileExt.FindString(filePath)

	fmt.Printf("filePath %s, fileExt %s\n", filePath, ext)

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
		return htmlMimeType
	case ".css":
		return cssMimeType
	case ".js":
		return jsMimeType
	case ".json":
		return jsonMimeType
	default:
		panic(fmt.Errorf("unsupported file extension: %s", fileExt))
	}
}
