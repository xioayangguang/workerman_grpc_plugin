package php

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

const phpBody = `<?php
# source: {{ .File.Name }}
{{ $ns := .Namespace -}}
{{if $ns.Namespace}}
namespace {{ $ns.Namespace }};
{{end}}
use Mix\Grpc;
use Mix\Grpc\Context;
{{- range $n := $ns.Import}}
use {{ $n }};
{{- end}}

interface {{ .Service.Name | interface }} extends Grpc\ServiceInterface
{
    public const NAME = "{{ .File.Package }}.{{ .Service.Name }}";{{ "\n" }}
{{- range $m := .Service.Method}}
    /**
    * @param Context $context
    * @param {{ name $ns $m.InputType }} $request
    * @return {{ name $ns $m.OutputType }}
    */
    public function {{ $m.Name }}(Context $context, {{ name $ns $m.InputType }} $request): {{ name $ns $m.OutputType }};
{{end -}}
}
`

type imports []string

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("phpBody").Funcs(template.FuncMap{
		"interface": func(name *string) string {
			return identifier(*name, "interface")
		},
		"name": func(ns *ns, name *string) string {
			return ns.resolve(name)
		},
	}).Parse(phpBody))
}

// generate php filename
func filename(file *descriptor.FileDescriptorProto, name *string) string {
	ns := namespace(file.Package, "/")
	if file.Options != nil && file.Options.PhpNamespace != nil {
		ns = strings.Replace(*file.Options.PhpNamespace, `\`, `/`, -1)
	}
	return fmt.Sprintf("%s/%s.php", ns, identifier(*name, "interface"))
}

// generate php file body
func body(
	req *plugin.CodeGeneratorRequest,
	file *descriptor.FileDescriptorProto,
	service *descriptor.ServiceDescriptorProto,
) string {
	out := bytes.NewBuffer(nil)

	data := struct {
		Namespace *ns
		File      *descriptor.FileDescriptorProto
		Service   *descriptor.ServiceDescriptorProto
	}{
		Namespace: newNamespace(req, file, service),
		File:      file,
		Service:   service,
	}

	err := tpl.Execute(out, data)
	if err != nil {
		panic(err)
	}

	return out.String()
}
