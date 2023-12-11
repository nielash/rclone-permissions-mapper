SHELL = bash
VERSION := 1.0

compile_all:
	go run bin/cross-compile.go -compile-only $(VERSION)
