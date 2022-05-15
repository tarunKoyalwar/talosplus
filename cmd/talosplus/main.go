package main

import (
	"github.com/tarunKoyalwar/talosplus/cmd/talosplus/subcmd"
)

func main() {

	err := subcmd.RootCMD.Execute()
	if err != nil {
		panic(err)
	}

}
