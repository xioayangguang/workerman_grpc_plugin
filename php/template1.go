package php

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

const phpBody1 = `<?php
{{ $ns := .Namespace -}}
{{if $ns.Namespace}}
namespace {{ $ns.Namespace }};
{{end}}
use Mix\Grpc;
use Mix\Grpc\Context;
{{- range $n := $ns.Import}}
use {{ $n }};
{{- end}}

class {{ .Service.Name | client }} extends XXXX
{
{{- range $m := .Service.Method}}
    /**
    * @param {{ name $ns $m.InputType }} $request
    * @param array $options
    * @return {{ name $ns $m.OutputType }}
    */
    public function {{ $m.Name }}({{ name $ns $m.InputType }} $request): {{ name $ns $m.OutputType }}
    {
        return $this->_simpleRequest('/{{ $.File.Package }}.{{ $.Service.Name }}/{{ $m.Name }}', $context, $request, new {{ name $ns $m.OutputType }}());
    }
{{end -}}
}
`

var tpl1 *template.Template

func init() {
	tpl1 = template.Must(template.New("phpBody1").Funcs(template.FuncMap{
		"client": func(name *string) string {
			return identifier(*name, "Client")
		},
		"name": func(ns *ns, name *string) string {
			return ns.resolve(name)
		},
	}).Parse(phpBody1))
}

// generate php filename
func filename1(file *descriptor.FileDescriptorProto, name *string) string {
	ns := namespace(file.Package, "/")
	if file.Options != nil && file.Options.PhpNamespace != nil {
		ns = strings.Replace(*file.Options.PhpNamespace, `\`, `/`, -1)
	}

	return fmt.Sprintf("%s/%s.php", ns, identifier(*name, "Client"))
}

// generate php file body
func body1(
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

	err := tpl1.Execute(out, data)
	if err != nil {
		panic(err)
	}

	return out.String()
}
