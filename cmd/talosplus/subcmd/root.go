package subcmd

import (
	"os"

	"github.com/spf13/cobra"
)

// MongoDB URL
var MongoURL string

// ClipBoardIn
var ClipBoardIn bool

var RootCMD cobra.Command = cobra.Command{
	Use:   "talosplus",
	Short: "Bash Scripting Middleware",
	Long:  `Create and Run Robust Bash Automation Scripts at Lightspeed with Notification support`,
}

func init() {
	RootCMD.PersistentFlags().StringVarP(&MongoURL, "url", "u", os.Getenv("MONGO_URL"), "MongoDB Connection String")
	RootCMD.PersistentFlags().BoolVarP(&ClipBoardIn, "clipin", "i", false, "Use this to get data from clipboard")
	RootCMD.AddCommand(&use, &get, &set, &RunScript)
}
