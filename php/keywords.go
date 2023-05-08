package php

import (
	"bytes"
	"github.com/c9s/inflect"
	"strings"
)

// @see https://github.com/protocolbuffers/protobuf/blob/master/php/ext/google/protobuf/protobuf.c#L168
var reservedKeywords = []string{
	"abstract", "and", "array", "as", "break",
	"callable", "case", "catch", "class", "clone",
	"const", "continue", "declare", "default", "die",
	"do", "echo", "else", "elseif", "empty",
	"enddeclare", "endfor", "endforeach", "endif", "endswitch",
	"endwhile", "eval", "exit", "extends", "final",
	"for", "foreach", "function", "global", "goto",
	"if", "implements", "include", "include_once", "instanceof",
	"insteadof", "interface", "isset", "list", "namespace",
	"new", "or", "print", "private", "protected",
	"public", "require", "require_once", "return", "static",
	"switch", "throw", "trait", "try", "unset",
	"use", "var", "while", "xor", "int",
	"float", "bool", "string", "true", "false",
	"null", "void", "iterable",
}

// Check if given name/keyword is reserved by php.
func isReserved(name string) bool {
	name = strings.ToLower(name)
	for _, k := range reservedKeywords {
		if name == k {
			return true
		}
	}

	return false
}

// generate php namespace or path
func namespace(pkg *string, sep string) string {
	if pkg == nil {
		return ""
	}

	result := bytes.NewBuffer(nil)
	for _, p := range strings.Split(*pkg, ".") {
		result.WriteString(identifier(p, ""))
		result.WriteString(sep)
	}

	return strings.Trim(result.String(), sep)
}

// create php identifier for class or message
func identifier(name string, suffix string) string {
	name = inflect.Camelize(name)
	if suffix != "" {
		return name + inflect.Camelize(suffix)
	}

	return name
}

func resolveReserved(identifier string, pkg string) string {
	if isReserved(strings.ToLower(identifier)) {
		if pkg == ".google.protobuf" {
			return "GPB" + identifier
		}
		return "PB" + identifier
	}
	return identifier
}
