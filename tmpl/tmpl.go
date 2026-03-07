package tmpl

import (
	"bytes"
	"html/template"
)

func HandleHtmlTemplate(file []byte, pageData *PageData) ([]byte, error) {
	tmpl, err := template.New("page").Parse(string(file))
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
