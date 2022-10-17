package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/projectdiscovery/goflags"
)

type DB_CLIENT_OPTS struct {
	READ_VAR  string
	WRITE_VAR string
	LIST_ALL  bool
	FROM_FILE string
}

type Templates_OPTS struct {
	DryRun      bool
	TemplateDir string
	Template    string
}

type Config_OPTS struct {
	Purge            bool
	SkipNotification bool
	Concurrency      int
	CacheDIR         string
	DiscordWID       string
	DiscordWTOKEN    string
	BlacklistVars    goflags.StringSlice
}

type DB_OPTS struct {
	UseMongoDB  bool
	DatabaseURI string
	DBName      string
	ContextName string
}

type Options struct {
	Templates_OPTS
	DB_CLIENT_OPTS
	Config_OPTS
	DB_OPTS

	Silent      bool // No Banner
	NoColor     bool //Disable Color Output
	ShowOutput  bool // Show Ouput of all commands
	Verbose     bool // Verbose Mode
	VeryVerbose bool
}

func parseInput() Options {
	opts := Options{}

	flagset := goflags.NewFlagSet()

	description := "Talosplus is a template based Automation Framework that harnesses power of GoLang"

	flagset.SetDescription(fmt.Sprintf("%v\n%v", fmt.Sprintf(banner, Version), description))

	flagset.CreateGroup("templates", "Templates",
		flagset.BoolVarP(&opts.DryRun, "dry-run", "n", false, "Dry/Test Run the template"),
		flagset.StringVarEnv(&opts.TemplateDir, "templates", "td", "", "TALOS_TEMPLATE_DIR", "Run all templates from directory [Env:TALOS_TEMPLATE_DIR]"),
		flagset.StringVarP(&opts.Template, "template", "t", "", "Run a single template"),
	)

	flagset.CreateGroup("output", "Output",
		flagset.BoolVar(&opts.Silent, "silent", false, "Don't Print Banner"),
		flagset.BoolVarP(&opts.NoColor, "no-color", "nc", false, "Disable Color Output"),
		flagset.BoolVarP(&opts.ShowOutput, "show", "s", false, "Show Output of All Commands"),
		flagset.BoolVarP(&opts.Verbose, "verbose", "v", false, "Verbose Mode (Show Scheduled Tasks & Warnings)"),
		flagset.BoolVarP(&opts.VeryVerbose, "very-verbose", "vv", false, "Max Verbosity"),
	)

	flagset.CreateGroup("configs", "Configurations",
		flagset.IntVarP(&opts.Concurrency, "limit", "c", 8, "Max Number of Concurrent Programs"),
		flagset.StringSliceVarP(&opts.BlacklistVars, "blacklist-vars", "b", []string{}, "Blacklist ", goflags.CommaSeparatedStringSliceOptions),
		flagset.BoolVarP(&opts.Purge, "purge", "p", false, "Purge Cache"),
		flagset.BoolVar(&opts.SkipNotification, "skip-notify", false, "Skip Sending Notification to Discord"),
		flagset.StringVarEnv(&opts.CacheDIR, "cache-dir", "cdir", os.TempDir(), "TALOS_CACHE_DIR", "Cache Directory [Env:TALOS_CACHE_DIR](All command outputs are saved here)"),
		flagset.StringVarEnv(&opts.DiscordWID, "discord-wid", "wid", "", "DISCORD_WID", "Discord Webhook ID [Env:DISCORD_WID]"),
		flagset.StringVarEnv(&opts.DiscordWTOKEN, "discord-wtoken", "wtoken", "", "DISCORD_WTOKEN", "Discord Webhook Token [Env:DISCORD_WTOKEN]"),
	)

	flagset.CreateGroup("database", "Database",
		flagset.StringVarEnv(&opts.DatabaseURI, "uri", "u", "", "TALOS_URI", "URI [Env: TALOS_URI] (MongoDB URL/Directory)"),
		flagset.StringVarEnv(&opts.DBName, "database-name", "db", "talosplus", "TALOS_DBNAME", "Database Name [Env:TALOS_DBNAME]"),
		flagset.StringVarEnv(&opts.ContextName, "context-name", "cn", "automation", "TALOS_CN", "Similar to Table Name in SQL can be anything version,subdomain etc [Env:TALOS_CN]"),
		flagset.BoolVar(&opts.UseMongoDB, "mongodb", false, "Use MongoDB (default : BBolt DB)"),
	)

	flagset.CreateGroup("dbclient", "DB Client(Similar to BBRF)",
		flagset.StringVarP(&opts.READ_VAR, "read-var", "get", "", "Read Variable Value from Database"),
		flagset.StringVarP(&opts.WRITE_VAR, "write-var", "put", "", "Save Data to Variable"),
		flagset.StringVarP(&opts.FROM_FILE, "file", "f", "", "Read From File"),
		flagset.BoolVarP(&opts.LIST_ALL, "list", "l", false, "List All Variables"),
	)

	if err := flagset.Parse(); err != nil {
		log.Fatalf("Could not parse flags: %s\n", err)
	}

	return opts
}

// HasStdin : Check if Stdin is present
func HasStdin() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	mode := stat.Mode()

	isPipedFromChrDev := (mode & os.ModeCharDevice) == 0
	isPipedFromFIFO := (mode & os.ModeNamedPipe) != 0

	return isPipedFromChrDev || isPipedFromFIFO
}

// GetStdin : Get all Data present on stdin
func GetStdin() string {
	bin, _ := ioutil.ReadAll(os.Stdin)
	return string(bin)
}
