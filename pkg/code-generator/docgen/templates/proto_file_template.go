package templates

import (
	"text/template"

	"github.com/solo-io/solo-kit/pkg/code-generator/codegen/templates"
)

var ProtoFileTemplate = template.Must(template.New("resource").Funcs(templates.Funcs).Parse(`
## Package:
{{ .Package }}

## Source File:
{{ .Name }} 

## Description:
{{- range .SyntaxComments.Detached }}  
{{ remove_magic_comments (printf "%v" .) }}
{{- end }}  

## Contents:
- Messages:
{{- range .Messages }}  
	- [{{ printfptr "%v" .Name }}](#{{.Name}})
{{- end }}

{{- if gt (len .Enums) 0 }} 
- Enums:
{{- range .Enums}}
	- [{{ printfptr "%v" .Name }}](#{{.Name}})
{{- end}}
{{- end}}

---
{{range .Messages }}  
### <a name="{{ printfptr "%v" .Name }}">{{ printfptr "%v" .Name }}</a>

Description: {{ remove_magic_comments .Comments.Leading }}

` + "```" + `yaml
{{range .Fields -}}
"{{ printfptr "%v" .Name}}": {{ fieldType . }}
{{end}}
` + "```" + `

| Field | Type | Description | Default |
| ----- | ---- | ----------- |----------- | 
{{range .Fields -}}
| {{ printfptr "%v" .Name }} | {{linkForType . }} | {{ remove_magic_comments (nobr .Comments.Leading) }} | {{if .DefaultValue}} Default: {{.DefaultValue}}{{end}} |
{{end}}

{{- end }}

<!-- Start of HubSpot Embed Code -->
<script type="text/javascript" id="hs-script-loader" async defer src="//js.hs-scripts.com/5130874.js"></script>
<!-- End of HubSpot Embed Code -->
`))
