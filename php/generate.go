package php

import (
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

// Generate generates needed service classes
func Generate(req *plugin.CodeGeneratorRequest) *plugin.CodeGeneratorResponse {
	resp := &plugin.CodeGeneratorResponse{}
	for _, file := range req.ProtoFile {
		for _, service := range file.Service {
			resp.File = append(resp.File, generate(req, file, service))
		}
	}
	for _, file := range req.ProtoFile {
		for _, service := range file.Service {
			resp.File = append(resp.File, generate1(req, file, service))
		}
	}
	return resp
}

func generate(
	req *plugin.CodeGeneratorRequest,
	file *descriptor.FileDescriptorProto,
	service *descriptor.ServiceDescriptorProto,
) *plugin.CodeGeneratorResponse_File {
	return &plugin.CodeGeneratorResponse_File{
		Name:    str(filename(file, service.Name)),
		Content: str(body(req, file, service)),
	}
}

func generate1(
	req *plugin.CodeGeneratorRequest,
	file *descriptor.FileDescriptorProto,
	service *descriptor.ServiceDescriptorProto,
) *plugin.CodeGeneratorResponse_File {
	return &plugin.CodeGeneratorResponse_File{
		Name:    str(filename1(file, service.Name)),
		Content: str(body1(req, file, service)),
	}
}

// helper to convert string into string pointer
func str(str string) *string {
	return &str
}
