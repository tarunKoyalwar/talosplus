package main

import (
	"bytes"
	"fmt"
	"os"
	"path"

	"github.com/tarunKoyalwar/talosplus/pkg/alerts"
	"github.com/tarunKoyalwar/talosplus/pkg/core"
	"github.com/tarunKoyalwar/talosplus/pkg/db"
	"github.com/tarunKoyalwar/talosplus/pkg/ioutils"
	"github.com/tarunKoyalwar/talosplus/pkg/shared"
	"github.com/tarunKoyalwar/talosplus/pkg/shell"
)

func main() {

	// Get Input
	opts := parseInput()

	showBanner()

	ioutils.Cout.DisableColor = opts.NoColor

	// Configure DB Settings

	if opts.UseMongoDB {
		// using mongodb
		err := db.UseMongoDB(opts.DatabaseURI, opts.DBName, opts.ContextName)
		if err != nil {
			ioutils.Cout.Fatalf(err, "Connection to MongoDB Failed %v", opts.DatabaseURI)
		}
		ioutils.Cout.PrintInfo("Connected to MongoDB")
	} else {
		opts.DBName += ".db"
		err := db.UseBBoltDB(opts.DatabaseURI, opts.DBName, opts.ContextName)
		if err != nil {
			ioutils.Cout.Fatalf(err, "Failed to open Database at %v %v", opts.DatabaseURI, opts.DBName)
		}
	}

	// Tasks related to Database Client

	if opts.READ_VAR != "" {
		val, err := db.DB.Get(opts.READ_VAR)
		if err != nil {
			ioutils.Cout.Fatalf(err, "Failed to retrieve %v from db", opts.READ_VAR)
		} else {
			fmt.Println(val)
		}
		os.Exit(0)
	} else if opts.WRITE_VAR != "" {
		var data string
		if HasStdin() {
			data = GetStdin()
		} else if opts.FROM_FILE != "" {
			bin, er := os.ReadFile(opts.FROM_FILE)
			if er != nil {
				ioutils.Cout.Fatalf(er, "Failed to read file %v", opts.FROM_FILE)
			} else {
				data = string(bin)
			}
		} else {
			ioutils.Cout.ErrExit("Input Missing . Exiting!!")
		}

		err := db.DB.Put(opts.WRITE_VAR, data, true)
		if err != nil {
			ioutils.Cout.Fatalf(err, "Failed to Write in %v", opts.WRITE_VAR)
		} else {
			fmt.Println("Saved to DB")
		}
		os.Exit(0)
	} else if opts.LIST_ALL {
		val, err := db.DB.GetAllVarNames()
		if err != nil {
			ioutils.Cout.Fatalf(err, "Failed to get list of variables from db")
		} else {
			for k := range val {
				fmt.Println(k)
			}
		}
		os.Exit(0)
	}

	//  Configure Discord Webhook

	if opts.DiscordWID != "" && opts.DiscordWTOKEN != "" {
		alerts.Alert = alerts.NewDiscordHook(opts.DiscordWID, opts.DiscordWTOKEN)
	}

	if opts.SkipNotification {
		alerts.Alert.Disabled = true
	}

	// Configure script args
	shell.Settings.CacheDIR = opts.CacheDIR
	shell.Settings.Limit = opts.Concurrency
	shell.Settings.ProjectName = opts.ContextName
	shell.Settings.ProjectExportName = opts.ContextName + "Exports"
	shell.Settings.Purge = opts.Purge

	// Template File Buffer
	var templateBuff bytes.Buffer

	// Configure Templates
	if opts.TemplateDir != "" {
		files, err := os.ReadDir(opts.TemplateDir)
		if err != nil {
			ioutils.Cout.Fatalf(err, "Failed to read templates from dir %v", opts.TemplateDir)
		}

		count := 0

		for _, v := range files {
			if !v.IsDir() {
				fbin, err := os.ReadFile(path.Join(opts.TemplateDir, v.Name()))
				if err != nil {
					ioutils.Cout.Printf("failed to read template %v got %v", opts.Template, err.Error())
				} else {
					templateBuff.WriteString("\n\n")
					templateBuff.Write(fbin)
					count += 1
				}
			}
		}

		ioutils.Cout.PrintInfo("Successfully Loaded %v Templates", count)

	} else if opts.Template != "" {
		// Load Templates
		fbin, err := os.ReadFile(opts.Template)
		if err != nil {
			ioutils.Cout.Fatalf(err, "failed to read template %v", opts.Template)
		}
		templateBuff.Write(fbin)
	} else {
		ioutils.Cout.ErrExit("No Templates Found. Exiting!!")
	}

	t := core.NewScripter()
	t.ShowOutput = opts.ShowOutput

	// Load Existing Variable Values
	shared.SharedVars.AddGlobalVarsFromDB()

	// Add Blacklisted Variables
	if len(opts.BlacklistVars) != 0 {
		for _, v := range opts.BlacklistVars {
			t.BlackList[v] = true
		}
	}

	t.Compile(templateBuff.String())

	t.Summarize()

	t.PrintAllCMDs()

	t.Schedule()

	if opts.DryRun {
		os.Exit(0)
	}

	t.Execute()

}
