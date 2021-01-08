package actiondoc

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"text/template"

	"gopkg.in/yaml.v2"
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

type ioVal struct {
	Key   string
	Props map[string]interface{}
}

type ioVals []ioVal

func (v *ioVals) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var slice yaml.MapSlice
	err := unmarshal(&slice)
	if err != nil {
		return err
	}
	var mp map[string]map[string]interface{}
	err = unmarshal(&mp)
	if err != nil {
		return err
	}
	*v = make(ioVals, len(slice))
	for i := 0; i < len(slice); i++ {
		key, ok := slice[i].Key.(string)
		if !ok {
			return fmt.Errorf("expected key to be a string but got a %T", slice[i].Key)
		}
		(*v)[i] = ioVal{
			Key:   key,
			Props: mp[key],
		}
	}
	return nil
}

type actionProperties struct {
	Name        string
	Author      string
	Description string
	Inputs      ioVals
	Outputs     ioVals
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

{{end}}{{$.Config.HeaderPrefix}}# Inputs

{{ range .Inputs -}}{{$.Config.HeaderPrefix}}## {{ .Key }}
{{if .Props.deprecationMessage }}
__Deprecated__ - {{.Props.deprecationMessage }}
{{end -}}
{{if .Props.required }}
__Required__
{{end -}}
{{if .Props.default }}
default: bktk{{ .Props.default }}bktk
{{end -}}
{{if .Props.description }}
{{.Props.description}}

{{end -}}
{{end -}}
{{- if .Outputs -}}
{{$.Config.HeaderPrefix}}# Outputs
{{range .Outputs}}
{{$.Config.HeaderPrefix}}## {{ .Key }}

{{ .Props.description }}
{{end -}}
{{end -}}{{- end -}}
`))
