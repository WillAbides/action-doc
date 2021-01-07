package actiondoc

import (
	"bytes"
	"io"
	"io/ioutil"
	"text/template"

	"github.com/ghodss/yaml"
)

// ActionMarkdown reads an action.yml file and returns some markdown suitable for a quality README.md
func ActionMarkdown(r io.Reader) ([]byte, error) {
	var props actionProperties
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(b, &props)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, &props)
	if err != nil {
		panic(err)
	}
	result := bytes.ReplaceAll(buf.Bytes(), []byte("bktk"), []byte("`"))
	return result, nil
}

type actionProperties struct {
	Name        string                            `json:"name"`
	Author      string                            `json:"author"`
	Description string                            `json:"description"`
	Inputs      map[string]actionInputProperties  `json:"inputs"`
	Outputs     map[string]actionOutputProperties `json:"outputs"`
}

type actionInputProperties struct {
	Description        string `json:"description"`
	DeprecationMessage string `json:"deprecationMessage"`
	Required           bool   `json:"required"`
	Default            string `json:"default"`
}

type actionOutputProperties struct {
	Description string `json:"description"`
}

var tmpl = template.Must(template.New("").Parse(`# {{ .Name }}

{{if .Author }}Author: {{.Author}}
    
{{end}}
{{- .Description }}

{{ if .Inputs }}## Inputs

{{ range $name, $props := .Inputs -}}### {{ $name }}
{{if $props.DeprecationMessage }}
__Deprecated__ - {{$props.DeprecationMessage }}
{{end -}}
{{if $props.Required }}
__Required__
{{end -}}
{{if $props.Default }}
default: bktk{{ $props.Default }}bktk
{{end -}}
{{if $props.Description }}
{{.Description}}

{{end -}}
{{end -}}
{{end -}}
{{- if .Outputs -}}
## Outputs
{{range $name, $props := .Outputs}}
### {{ $name }}

{{ $props.Description }}
{{end -}}
{{end -}}
`))
