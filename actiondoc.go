package actiondoc

import (
	"bytes"
	"io"
	"io/ioutil"
	"text/template"

	"github.com/ghodss/yaml"
)

type markdownConfig struct {
	SkipActionName        bool
	SkipActionDescription bool
	SkipActionAuthor      bool
	HeaderPrefix          string
}

// MarkdownOption is an option for configuring markdown output
type MarkdownOption func(config *markdownConfig)

// SkipActionName skip outputting the name of the action
func SkipActionName(val bool) MarkdownOption {
	return func(config *markdownConfig) {
		config.SkipActionName = val
	}
}

// SkipActionDescription skip outputting the description of the action
func SkipActionDescription(val bool) MarkdownOption {
	return func(config *markdownConfig) {
		config.SkipActionDescription = val
	}
}

// SkipActionAuthor skip outputting the action author
func SkipActionAuthor(val bool) MarkdownOption {
	return func(config *markdownConfig) {
		config.SkipActionAuthor = val
	}
}

// HeaderPrefix some extra #s to put in front of each markdown header
func HeaderPrefix(val string) MarkdownOption {
	return func(config *markdownConfig) {
		config.HeaderPrefix = val
	}
}

// ActionMarkdown reads an action.yml file and returns some markdown suitable for a quality README.md
func ActionMarkdown(r io.Reader, option ...MarkdownOption) ([]byte, error) {
	var cfg markdownConfig
	for _, opt := range option {
		opt(&cfg)
	}
	if !cfg.SkipActionName {
		cfg.HeaderPrefix += "#"
	}
	var props actionProperties
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(b, &props)
	if err != nil {
		return nil, err
	}
	data := tmplData{
		Config:     cfg,
		Properties: props,
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, &data)
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

type tmplData struct {
	Config     markdownConfig
	Properties actionProperties
}

var tmpl = template.Must(template.New("").Parse(`
{{- with .Properties}}{{ if not $.Config.SkipActionName }}{{$.Config.HeaderPrefix}} {{ .Name }}
{{ end }}
{{ if not $.Config.SkipActionAuthor }}{{if .Author }}Author: {{.Author}}

{{end}}{{end}}
{{- if not $.Config.SkipActionDescription }}{{ .Description }}

{{end}}{{ if .Inputs }}{{$.Config.HeaderPrefix}}# Inputs

{{ range $name, $props := .Inputs -}}{{$.Config.HeaderPrefix}}## {{ $name }}
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
{{$.Config.HeaderPrefix}}# Outputs
{{range $name, $props := .Outputs}}
{{$.Config.HeaderPrefix}}## {{ $name }}

{{ $props.Description }}
{{end -}}
{{end -}}{{- end -}}
`))
