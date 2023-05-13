package main

import (
	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"io"
	"io/ioutil"
	"os"
	"workermangrpc/php"
)

func main() {
	req, err := readRequest(os.Stdin)
	if err != nil {
		panic(err)
	}
	if err = writeResponse(os.Stdout, php.Generate(req)); err != nil {
		panic(err)
	}
}

func readRequest(in io.Reader) (*plugin.CodeGeneratorRequest, error) {
	data, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}
	req := new(plugin.CodeGeneratorRequest)
	if err = proto.Unmarshal(data, req); err != nil {
		return nil, err
	}
	return req, nil
}

func writeResponse(out io.Writer, resp *plugin.CodeGeneratorResponse) error {
	data, err := proto.Marshal(resp)
	if err != nil {
		return err
	}
	_, err = out.Write(data)
	return err
}
