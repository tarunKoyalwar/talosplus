package main

import (
	_ "embed"

	"github.com/tarunKoyalwar/talosplus/pkg/ioutils"
)

//go:embed banner.txt
var banner string

var Version string = "v1.1.0"

func showBanner() {
	ioutils.Cout.Printf(banner, Version)
}
