package subcmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tarunKoyalwar/talosplus/pkg/shared"
)

var show bool
var list bool
var local bool

var get cobra.Command = cobra.Command{
	Use:   "get",
	Short: "Get Specific Document/Variable Value",
	Long:  "Retrieve Given variable value from DB",
	RunE: func(cmd *cobra.Command, args []string) error {
		LoadSettings()

		if local {
			fmt.Printf("[+] Local Settings\n")
			fmt.Printf("%-25v: %v\n", "Active DB URL", DefaultSettings.ActiveURL)
			fmt.Printf("%-25v: %v\n", "Active Database", DefaultSettings.ActiveDB)
			fmt.Printf("%-25v: %v\n", "Active Program", DefaultSettings.ActiveColl)
			fmt.Printf("%-25v: %v\n", "Active Bash Script", DefaultSettings.ActiveScript)
			return nil
		}

		if !DefaultSettings.Available() {
			if dbname == "" || collname == "" {
				return fmt.Errorf("either use parameters `db,coll` or `use` command to set required values")
			} else {
				DefaultSettings.ActiveDB = dbname
				DefaultSettings.ActiveColl = collname
			}
		}

		err := PrepareDB()
		if err != nil {
			return err
		}

		if show {
			zx, errx := shared.LoadAllExplicitVars()
			if errx != nil {
				return errx
			} else {
				fmt.Printf("[*] Explicitly Set Variables in DB\n\n")
				for k, v := range zx {
					fmt.Printf("%-16v : %v\n", k, v)
				}
				return nil
			}
		} else if list {
			zx, errx := shared.GetAllVarNames()
			if errx != nil {
				return errx
			} else {
				fmt.Printf("[*] Explicitly Set Variables in DB\n\n")
				for k, v := range zx {
					if v {
						fmt.Println(k)
					}
				}

				fmt.Printf("\n\n[*] Executed Commands Set Variables in DB\n\n")
				for k, v := range zx {
					if !v {
						fmt.Println(k)
					}
				}

				return nil
			}
		}

		if len(args) < 1 {
			return fmt.Errorf("no args given required positional arg")
		}

		key := strings.TrimSpace(args[0])

		val, err := shared.GetFromDB(key)

		if err != nil {
			return err
		} else {
			fmt.Printf("%v", val)
		}

		return nil
	},
}

func init() {
	get.Flags().StringVar(&dbname, "db", "", "Database Name to use")
	get.Flags().StringVarP(&collname, "program", "p", "", "Program/Collection Name")
	get.Flags().BoolVar(&show, "show", false, "Show All Enviornment Variables present in DB")
	get.Flags().BoolVar(&list, "list", false, "List All Variable Names Present in DB")
	get.Flags().BoolVar(&local, "local", false, "Show Local Settings [set using `use`]")
}
