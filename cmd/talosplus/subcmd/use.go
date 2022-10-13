package subcmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tarunKoyalwar/talosplus/pkg/db"
	"github.com/tarunKoyalwar/talosplus/pkg/db/mongox"
)

var dbname string

var collname string

var scriptpath string

var use cobra.Command = cobra.Command{
	Use:   "use",
	Short: "Use/Set Parameters",
	Long:  "Use/Set Required Parameters like dbname,dburl,collname,activescript",
	RunE: func(cmd *cobra.Command, args []string) error {
		LoadSettings()

		if dbname != "" {
			DefaultSettings.ActiveDB = dbname
		}

		if collname != "" {
			DefaultSettings.ActiveColl = collname
		}

		if scriptpath != "" {
			_, err := os.Stat(scriptpath)
			if err != nil {
				return fmt.Errorf("script file does not exist")
			} else {
				DefaultSettings.ActiveScript = scriptpath
			}
		}

		if MongoURL != "" {
			DefaultSettings.ActiveURL = MongoURL
		}

		DefaultSettings.Save()

		// Create New Collection if it does not exist
		if DefaultSettings.ActiveColl != "" {
			err := PrepareDB()
			if err != nil {
				return err
			}
			return CreateProgramIfNotExists(DefaultSettings.ActiveColl)
		}

		return nil

	},
}

func init() {
	use.Flags().StringVar(&dbname, "db", "", "Databse Name to use")
	use.Flags().StringVarP(&collname, "program", "p", "", "Collection / Program Name")
	use.Flags().StringVarP(&scriptpath, "script", "s", "", "Script Path")
	use.MarkFlagFilename("script")

}

func CreateProgramIfNotExists(collname string) error {

	//check if collection exists
	arr, er := db.DB.(*mongox.Provider).ListDBCollections()
	if er != nil {
		return er
	}

	found := false
	for _, v := range arr {
		if v == collname {
			found = true
		}
	}

	// create new collection /program
	if !found {
		err := db.DB.(*mongox.Provider).CreateCollection(collname)
		return err
	}

	return nil
}
