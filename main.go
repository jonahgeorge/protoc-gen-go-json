package main

import (
	"text/template"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

func main() {
	pgs.Init(
		pgs.DebugEnv("DEBUG"),
	).RegisterModule(
		JSONify(),
	).RegisterPostProcessor(
		pgsgo.GoFmt(),
	).Render()
}

// JSONifyPlugin adds encoding/json Marshaler and Unmarshaler methods on PB
// messages that utilizes the more correct jsonpb package.
// See: https://godoc.org/github.com/golang/protobuf/jsonpb
type JSONifyModule struct {
	*pgs.ModuleBase
	ctx pgsgo.Context
	tpl *template.Template
}

// JSONify returns an initialized JSONifyPlugin
func JSONify() *JSONifyModule { return &JSONifyModule{ModuleBase: &pgs.ModuleBase{}} }

func (p *JSONifyModule) InitContext(c pgs.BuildContext) {
	p.ModuleBase.InitContext(c)
	p.ctx = pgsgo.InitContext(c.Parameters())

	tpl := template.New("jsonify").Funcs(map[string]interface{}{
		"package":         p.ctx.PackageName,
		"name":            p.ctx.Name,
		"marshaler":       p.marshaler,
		"unmarshaler":     p.unmarshaler,
		"discard_unknown": p.discardUnknown,
	})

	p.tpl = template.Must(tpl.Parse(jsonifyTpl))
}

// Name satisfies the generator.Plugin interface.
func (p *JSONifyModule) Name() string { return "jsonify" }

func (p *JSONifyModule) Execute(targets map[string]pgs.File, pkgs map[string]pgs.Package) []pgs.Artifact {
	for _, t := range targets {
		p.generate(t)
	}

	return p.Artifacts()
}

func (p *JSONifyModule) generate(f pgs.File) {
	if len(f.Messages()) == 0 {
		return
	}

	name := p.ctx.OutputPath(f).SetExt(".json.go")
	p.AddGeneratorTemplateFile(name.String(), p.tpl, f)
}

func (p *JSONifyModule) discardUnknown(m pgs.Message) bool {
	b, err := p.ctx.Params().Bool("discard_unknown")
	if err != nil {
		return false
	}
	return b
}

func (p *JSONifyModule) marshaler(m pgs.Message) pgs.Name {
	return p.ctx.Name(m) + "JSONMarshaler"
}

func (p *JSONifyModule) unmarshaler(m pgs.Message) pgs.Name {
	return p.ctx.Name(m) + "JSONUnmarshaler"
}

const jsonifyTpl = `package {{ package . }}

import (
	"encoding/json"

  "google.golang.org/protobuf/encoding/protojson"
)

{{ range .AllMessages }}
// {{ marshaler . }} describes the default jsonpb.Marshaler used by all
// instances of {{ name . }}. This struct is safe to replace or modify but
// should not be done so concurrently.
var {{ marshaler . }} = protojson.MarshalOptions{}

// MarshalJSON satisfies the encoding/json Marshaler interface. This method
// uses the more correct jsonpb package to correctly marshal the message.
func (m *{{ name . }}) MarshalJSON() ([]byte, error) {
	if m == nil {
		return json.Marshal(nil)
	}
	return {{ marshaler . }}.Marshal(m)
}

var _ json.Marshaler = (*{{ name . }})(nil)

// {{ unmarshaler . }} describes the default jsonpb.Unmarshaler used by all
// instances of {{ name . }}. This struct is safe to replace or modify but
// should not be done so concurrently.
var {{ unmarshaler . }} = protojson.UnmarshalOptions{
  DiscardUnknown: {{ discard_unknown . }},
}

// UnmarshalJSON satisfies the encoding/json Unmarshaler interface. This method
// uses the more correct jsonpb package to correctly unmarshal the message.
func (m *{{ name . }}) UnmarshalJSON(b []byte) error {
	return {{ unmarshaler . }}.Unmarshal(b, m)
}

var _ json.Unmarshaler = (*{{ name . }})(nil)
{{ end }}
`
