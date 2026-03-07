package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/ishchenko-gv/ishchenko-dev/tmpl"
)

func main() {
	err := os.RemoveAll("dist")
	if err != nil {
		log.Fatal(err)
	}

	err = filepath.WalkDir("src", handleWalk)
	if err != nil {
		log.Fatal(err)
	}
}

func handleWalk(path string, d fs.DirEntry, err error) error {
	if d.IsDir() {
		return os.Mkdir(distRoot(path), 0755)
	}

	f, err := os.ReadFile(path)
	ext := filepath.Ext(path)
	if ext == ".html" {
		f, err = tmpl.HandleHtmlTemplate(f, &tmpl.PageData{})
		if err != nil {
			return err
		}
	}

	log.Println("filepath:", path)

	err = os.WriteFile(distRoot(path), f, 0755)
	if err != nil {
		return err
	}

	return nil
}

func distRoot(path string) string {
	p, err := filepath.Rel("src", path)
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join("dist", p)
}
