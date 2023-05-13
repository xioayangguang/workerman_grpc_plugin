package php

import (
	"bytes"
	"fmt"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"strings"
	"text/template"
)

const phpBody = `<?php
declare(strict_types=1);
# source: {{ .File.Name }}  
{{ $ns := .Namespace -}}
{{if $ns.Namespace}}
namespace {{ $ns.Namespace }};
{{end}}
use parse\Http2Stream;
use parse\Response;

{{- range $n := $ns.Import}}
use {{ $n }};
{{- end}}


class {{ .Service.Name | service }}
{
	 public static $Streaming = [
	{{- range $k,$v := .Ttype}}
		"{{$k}}"=>[{{- range $kk,$vv := $v }}"{{ $vv }}",{{- end}}],
	{{- end}}
	];
	
	{{ $s := .Service -}}
	{{ $f := .File -}}
	 public static $Route  = [
	{{- range $m := .Service.Method}}
		"/{{ $f.Package }}.{{ $s.Name }}/{{ $m.Name }}" => [{{ $s.Name }}Service::class, "{{ $m.Name }}"],
	{{- end}}
	];

     public static $Parameter  = [
	{{- range $m := .Service.Method}}
		"/{{ $f.Package }}.{{ $s.Name }}/{{ $m.Name }}" => {{ name $ns $m.InputType }}::class,
	{{- end}}
	];

    public const NAME = "{{ .File.Package }}.{{ .Service.Name }}";

{{- range $m := .Service.Method}}
	{{if $m.ClientStreaming }}
		{{if $m.ServerStreaming }}
	/**
	*此处实现自己的业务逻辑
	* DoubleStreaming
	* @param {{ name $ns $m.InputType }} $request
	*/
	public static function {{ $m.Name }}({{ name $ns $m.InputType }} $request) 
    {
	}
		{{else}}
	/**
	*此处实现自己的业务逻辑
	* ClientStreaming
	* @param {{ name $ns $m.InputType }} $request
	*/
	public static function {{ $m.Name }}({{ name $ns $m.InputType }} $request) : {{ name $ns $m.OutputType }}
    {
	}
		{{end}}
	{{else}}
		{{if $m.ServerStreaming }}
	/**
	*此处实现自己的业务逻辑
	* ServerStreaming
	* @param {{ name $ns $m.InputType }} $request
	*/
	public static function {{ $m.Name }}({{ name $ns $m.InputType }} $request) 
    {
	}
		{{else}}
	/**
	*此处实现自己的业务逻辑
	* Simple
	* @param {{ name $ns $m.InputType }} $request
	* @return {{ name $ns $m.OutputType }}
	*/
	public static function {{ $m.Name }}({{ name $ns $m.InputType }} $request): {{ name $ns $m.OutputType }} 
    {
	}
		{{end}}
{{end}}
{{end -}}
}
`

type imports []string

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("phpBody").Funcs(template.FuncMap{
		"service": func(name *string) string {
			return identifier(*name, "Service")
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
	return fmt.Sprintf("%s/%s.php", ns, identifier(*name, "Service"))
}

// generate php file body
func body(
	req *plugin.CodeGeneratorRequest,
	file *descriptor.FileDescriptorProto,
	service *descriptor.ServiceDescriptorProto,
) string {
	out := bytes.NewBuffer(nil)
	ttype := make(map[string][]string)
	for _, m := range service.Method {
		path := fmt.Sprintf("/%s.%s/%s", *file.Package, *service.Name, *m.Name)
		if m.ServerStreaming != nil && m.ClientStreaming != nil {
			if *m.ServerStreaming && *m.ClientStreaming {
				ttype["double_streaming"] = append(ttype["double_streaming"], path)
				continue
			}
		}
		if m.ServerStreaming != nil && *m.ServerStreaming {
			ttype["server_streaming"] = append(ttype["server_streaming"], path)
			continue
		}
		if m.ClientStreaming != nil && *m.ClientStreaming {
			ttype["client_streaming"] = append(ttype["client_streaming"], path)
			continue
		}
		ttype["simple"] = append(ttype["simple"], path)
	}
	data := struct {
		Namespace *ns
		File      *descriptor.FileDescriptorProto
		Service   *descriptor.ServiceDescriptorProto
		Ttype     map[string][]string
	}{
		Namespace: newNamespace(req, file, service),
		File:      file,
		Service:   service,
		Ttype:     ttype,
	}
	err := tpl.Execute(out, data)
	if err != nil {
		panic(err)
	}

	return out.String()
}
