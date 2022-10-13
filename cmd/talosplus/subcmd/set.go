package subcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tarunKoyalwar/talosplus/pkg/db"
)

var varname string

var set cobra.Command = cobra.Command{
	Use:   "set",
	Short: "Set/Update Variable/Document Value",
	Long:  "Update Values of Variables Stored in Database",
	RunE: func(cmd *cobra.Command, args []string) error {
		LoadSettings()

		if !DefaultSettings.Available() {
			if dbname == "" || collname == "" {
				return fmt.Errorf("either use parameters `db,coll` or `use` command to set required values")
			} else {
				DefaultSettings.ActiveDB = dbname
				DefaultSettings.ActiveColl = collname
			}
		}

		if varname == "" {
			return fmt.Errorf("enter the variable name for which data needs to be set")
		}

		err := PrepareDB()
		if err != nil {
			return err
		}
		var data string
		manageinput(&data, args)

		errres := db.DB.Put(varname, data, true)
		if errres != nil {
			return fmt.Errorf("failed to save to db %v", errres)
		} else {
			fmt.Printf("Saved in buffer/variable %v\n", varname)
		}

		return nil
	},
}

func init() {
	set.Flags().StringVar(&dbname, "db", "", "Databse Name to use")
	set.Flags().StringVarP(&collname, "program", "p", "", "Program/Collection Name")
	set.Flags().StringVar(&varname, "var", "", "Variable/Buffer Name to set Value for")
}
