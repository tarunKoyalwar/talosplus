package subcmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/tarunKoyalwar/talosplus/pkg/alerts"
	"github.com/tarunKoyalwar/talosplus/pkg/core"
	"github.com/tarunKoyalwar/talosplus/pkg/ioutils"
	"github.com/tarunKoyalwar/talosplus/pkg/shared"
	"github.com/tarunKoyalwar/talosplus/pkg/shell"
)

// Default Settings
var (
	limit   = 4 // max number of concurrent programs
	Purge   = false
	Verbose = false

	CachedDir    = os.TempDir()
	DiscordId    = ""
	DiscordToken = ""
	ShowOutput   = false
	DryRun       = false
	blacklist    = []string{}
	Nocolor      = false
)

var RunScript cobra.Command = cobra.Command{
	Use:   "run",
	Short: "Run Given Script",
	Long: `Execute Bash Script With All Features of talos and save to db
Settings Like  Purge, CacheDIR , pname etc can also be set using "get|set"
	`,
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

		if scriptpath != "" {
			DefaultSettings.ActiveScript = scriptpath

		}

		if DefaultSettings.ActiveScript == "" {
			return fmt.Errorf("Bash Script Not Supplied")
		}

		var scriptdata string
		bin, err := ioutil.ReadFile(DefaultSettings.ActiveScript)
		if err != nil {
			return err
		}

		scriptdata = string(bin)

		// Set Given Settings
		shell.Settings.CacheDIR = CachedDir
		shell.Settings.Limit = limit
		shell.Settings.ProjectName = DefaultSettings.ActiveColl
		shell.Settings.ProjectExportName = DefaultSettings.ActiveColl + "Exports"
		shell.Settings.Purge = Purge

		// Set Verbosity
		ioutils.Cout.Verbose = Verbose
		ioutils.Cout.DisableColor = Nocolor

		// check env for cachedir
		cdirenv := os.Getenv("TALOS_CACHEDIR")
		if cdirenv != "" {
			shell.Settings.CacheDIR = cdirenv
		}

		// Configure Notifications if given
		SetupAlerts()

		// Connect to mongodb
		if DefaultSettings.ActiveDB == "" || DefaultSettings.ActiveColl == "" {
			return fmt.Errorf("database and program name must be given using `use` or cmdline args")
		} else {
			err := PrepareDB()
			if err != nil {
				fmt.Printf("\ntalosplus requires mongodb\n")
				return err
			}
		}

		// Create Script Engine
		s := core.NewScripter()
		s.ShowOutput = ShowOutput

		//load data from db
		shared.SharedVars.AddGlobalVarsFromDB()

		if len(blacklist) > 0 {
			for _, v := range blacklist {
				s.BlackList[v] = true
			}
		}

		// Compile Given Bash Script
		s.Compile(scriptdata)

		for _, v := range s.IndexedCMDs {
			if v.Alerts == nil {
				fmt.Println("nil after compile")
			}
		}

		// Summarize all details of bash script
		s.Summarize()

		// Pretty print identified commands
		s.PrintAllCMDs()

		// Schedule commands [Must]
		s.Schedule()

		if DryRun {
			if alerts.Alert == nil {
				fmt.Printf("[warn] Discord Tokens Not Found\n")
			}
			os.Exit(1)
		}

		// Execute All Commands
		s.Execute()

		return nil
	},
}

func init() {
	RunScript.Flags().SortFlags = false
	RunScript.Flags().StringVar(&dbname, "db", "", "Database Name to use")
	RunScript.Flags().StringVarP(&collname, "program", "p", "", "Program/Collection Name")
	RunScript.Flags().StringVarP(&scriptpath, "script", "s", "", "Script Path")

	// Flags Related to talos
	RunScript.Flags().StringVar(&DiscordId, "id", "", "Discord Webhook ID(ENV: DISCORD_WID)")
	RunScript.Flags().StringVar(&DiscordToken, "token", "", "Discord Webhook Token(ENV: DISCORD_WTOKEN")
	RunScript.Flags().StringVar(&CachedDir, "cdir", os.TempDir(), "Cache Directory(ENV: TALOS_CACHEDIR)")
	RunScript.Flags().IntVarP(&limit, "limit", "t", 4, "Max Number of Concurrent Programs")
	RunScript.Flags().BoolVar(&ShowOutput, "show", false, "Show Output of Each Command")
	RunScript.Flags().BoolVar(&Purge, "purge", false, "Purge Cached Data")
	RunScript.Flags().BoolVarP(&Verbose, "verbose", "v", false, "Verbose Output (Includes Warn,Info)")
	RunScript.Flags().BoolVar(&Nocolor, "no-color", false, "Disable Colored Output")
	RunScript.Flags().BoolVar(&DryRun, "dryrun", false, "Do Everything Except Running the commands")
	RunScript.Flags().StringSliceVarP(&blacklist, "blacklist", "b", []string{}, "These variables will not be used from db and cmds will be rerun")

}

// SetupAlerts : Configure Discord using valid tokens
func SetupAlerts() {
	id := ""
	token := ""
	if DiscordId != "" && DiscordToken != "" {
		id = DiscordId
		token = DiscordToken
	} else {
		id = os.Getenv("DISCORD_WID")
		token = os.Getenv("DISCORD_WTOKEN")
	}
	if id != "" && token != "" {

		alerts.Alert = alerts.NewDiscordHook(id, token)
	}

	if alerts.Alert != nil {
		alerts.Alert.Title = DefaultSettings.ActiveColl
	}
}
