package main

import (
	"embed"
	"io/fs"

	"github.com/afq984/ilmsmock/internal/cmdmain"
	"github.com/afq984/ilmsmock/internal/proxy"
	"github.com/spf13/pflag"
)

//go:embed cache/*
var cache embed.FS

func main() {
	fs, err := fs.Sub(cache, "cache")
	if err != nil {
		panic(err)
	}

	pflag.Parse()

	cmdmain.Main(proxy.NewFileSystemUpstream(fs))
}
