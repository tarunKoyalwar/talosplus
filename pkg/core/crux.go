package core

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/tarunKoyalwar/talosplus/pkg/db"
	"github.com/tarunKoyalwar/talosplus/pkg/internal"
	"github.com/tarunKoyalwar/talosplus/pkg/ioutils"
	"github.com/tarunKoyalwar/talosplus/pkg/stringutils"

	"github.com/tarunKoyalwar/talosplus/pkg/scheduler"
	"github.com/tarunKoyalwar/talosplus/pkg/shared"
	"github.com/tarunKoyalwar/talosplus/pkg/shell"
	"github.com/tarunKoyalwar/talosplus/pkg/workshop"
)

// Engine : Template Processing Engine of Talosplus
type Engine struct {
	CMDs        []*shell.CMDWrap          // Command Instances
	IndexedCMDs map[string]*shell.CMDWrap // Indexed Commands
	ExecPyramid [][]*scheduler.Node       // Scheduled / Indexed Pyramid
	ShowOutput  bool                      // ShowOutput of every command
	BlackList   map[string]bool           // Blacklist specific variable (Beta)
}

// fillindex : Index all commands and store them in map
func (e *Engine) fillindex() {
	for _, v := range e.CMDs {
		if v.Alerts == nil {
			ioutils.Cout.PrintInfo("Alerts Instance of CMDWrap is nil")
		}
		e.IndexedCMDs[v.UID] = v
	}
}

// Compile : Compile template
func (e *Engine) Compile(shellscript string) {

	cmdarr := internal.ParseScript(shellscript)

	for _, v := range cmdarr {
		// create cmdwraps from basic commands
		wrap := shell.NewCMDWrap(v.Raw, v.Comment)
		e.CMDs = append(e.CMDs, &wrap)

	}
	e.fillindex()
}

// Schedule : Schedule programs by analyzing dependencies
func (e *Engine) Schedule() {

	t := scheduler.NewScheduler()

	// remove commands from scheduling tree whose value is available
	arr := e.removeAvailable()

	//these are needless commands that should not run at all
	needless := map[string]bool{}
	for _, v := range arr {
		addresses := shared.DefaultRegistry.VarAddress[v]
		for _, addr := range addresses {
			needless[addr] = true
		}
	}

	ioutils.Cout.Header("[+] Skipping Following commands\n")

	for k := range needless {
		val := e.IndexedCMDs[k].Comment
		if val == "" {
			val = e.IndexedCMDs[k].Raw
		}
		ioutils.Cout.PrintColor(ioutils.Azure, "%v : %v", k, val)
	}

	ioutils.Cout.Seperator(60)

	// autobalance dependencies of blacklisted nodes
	t.BlackListed = needless

	for _, v := range e.CMDs {
		t.AddNode(v.UID, v.Comment)
	}

	// Run Scheduler
	t.Run()

	e.ExecPyramid = t.ExecPyramid

}

// removeAvailable : [BETA]remove Commands from Schedule whose output is already available/assigned by user
func (e *Engine) removeAvailable() []string {
	updatedvars := []string{}
	if db.DB == nil {
		return updatedvars
	}

	//get all runtime vars
	z, err := db.DB.GetAllImplicit()
	if err != nil {
		return updatedvars
	}

	/*
		for every runtime variable
		remove all its dependents
		TODO Add Blacklist

	*/
	for k, v := range z {
		if shared.DefaultRegistry.Dependents[k] != nil {

			// If BlackListed Will run commands again
			if e.BlackList[k] {
				continue
			}

			//unmark these runtime variables as dependents and fill their values
			delete(shared.DefaultRegistry.Dependents, k)

			//fill value of this in the shared variable store
			err := shared.SharedVars.Set(k, v, false)
			if err != nil {
				ioutils.Cout.PrintWarning("failed to add variable from db %v", err.Error())
			} else {
				updatedvars = append(updatedvars, k)
			}
		}
	}

	return updatedvars

}

// Execute : Will Execute Template In orderly Fashion
func (e *Engine) Execute() {

	count := 0

	for _, v := range e.ExecPyramid {
		ioutils.Cout.Header("[^_^] Executing Level %v Commands\n", count)
		count += 1

		queue := []shell.CMDWrap{}

		//check if it for loop and dissolve for each level
		for _, ftest := range v {

			uid := ftest.UID
			c := e.IndexedCMDs[uid]

			if c == nil {
				ioutils.Cout.PrintWarning("This was not supposed to happen")
				ioutils.Cout.PrintWarning("%v with %v Not Found", ftest.Comment, uid)
			} else {

				c.Process()

				if !c.IsForLoop {
					queue = append(queue, *c)

				} else {
					dissolved, er := c.Disolve()

					if er != nil {
						ioutils.Cout.Printf("[-] %v Will Not be Executed because :%v\n", c.Comment, ioutils.Cout.ErrColor(er).Bold())
					} else {
						for _, tinstance := range dissolved {

							dx := tinstance

							dx.Process()

							e.IndexedCMDs[dx.UID] = &dx
							e.CMDs = append(e.CMDs, &dx)
							queue = append(queue, dx)
						}
					}

				}

			}
		}

		finalqueue := []*shell.CMDWrap{}

		for _, c := range queue {
			//All Checks Passed
			if !c.IsInvalid {

				ioutils.Cout.PrintInfo("(*) Scheduled... %v", strings.Join(c.CMD.Cmdsplit, " "))
				finalqueue = append(finalqueue, e.IndexedCMDs[c.UID])
				// fmt.Println(unsafe.Sizeof(c))

			} else {
				ioutils.Cout.Printf("[-] %v Will Not be Executed because :\n%v", c.Comment, ioutils.Cout.GetColor(ioutils.Azure, strings.Join(c.CauseofFailure, "\n")))
			}
		}

		workshop.ExecQueue(finalqueue, shared.DefaultSettings.Limit, e.ShowOutput)
		ioutils.Cout.Seperator(60)

	}
	//cleanup

	defer cleanup()

}

// Evaluate : Summarizes & Evaluate All Script Data
func (e *Engine) Evaluate() {

	tmp := []string{}

	if ioutils.Cout.VeryVerbose {
		//Only in very verbose Mode

		ioutils.Cout.Header("\n[*] Parsed Settings\n")
		ioutils.Cout.Value("%-16v : %v", "Purge Cache", shared.DefaultSettings.Purge)
		ioutils.Cout.Value("%-16v : %v", "Concurrency", shared.DefaultSettings.Limit)
		ioutils.Cout.Value("%-16v : %v", "ProjectName", shared.DefaultSettings.ProjectName)
		ioutils.Cout.Value("%-16v : %v", "CacheDir", shared.DefaultSettings.CacheDIR)
		ioutils.Cout.Value("%-16v : %v", "Verbose", ioutils.Cout.Verbose)

		ioutils.Cout.Seperator(60)

		gvars := shared.SharedVars.GetGlobalVars()

		ioutils.Cout.Header("\n[*] Used Explicit declared Variables\n")

		for k, v := range gvars {
			tarr := strings.Split(v, "\n")
			if len(tarr) == 1 {
				ioutils.Cout.Value("%-16v : %v", k, tarr[0])
			} else {
				ioutils.Cout.Value("%-16v : %v", k, tarr[0])
				for _, zx := range tarr[1:] {
					ioutils.Cout.Value("%-16v : %v", "", zx)
				}
			}

		}

		ioutils.Cout.Seperator(60)

		ioutils.Cout.Header("\n[*] Generated UIDs For Commands\n")
		for k, v := range e.IndexedCMDs {
			identifier := v.Comment
			if identifier == "" {
				identifier = v.Raw
			}
			ioutils.Cout.Value("%v : %v", k, identifier)
		}

		ioutils.Cout.Seperator(60)

		ioutils.Cout.Header("[*] Dependencies Found\n")

	}

	// Evaluate required , available and extra variables
	for k, v := range shared.DefaultRegistry.Dependents {
		if len(v) != 0 {

			vaddress := shared.DefaultRegistry.VarAddress[k]
			providers := []string{}

			for _, addr := range vaddress {
				if e.IndexedCMDs[addr] != nil {
					providers = append(providers, e.IndexedCMDs[addr].Comment)
				} else {
					providers = append(providers, k)
				}
			}

			requiredby := []string{}

			for _, uid := range v {
				if e.IndexedCMDs[uid] != nil {
					inst := e.IndexedCMDs[uid]
					requiredby = append(requiredby, inst.Comment)
				}
			}

			if ioutils.Cout.VeryVerbose {
				// Print Only in Verbose Mode

				zx := fmt.Sprintf("[ %v ] Will be Executed After [%v]\n", strings.Join(requiredby, " , "), strings.Join(providers, " , "))
				ioutils.Cout.Printf("%v", ioutils.Cout.GetColor(ioutils.LightGreen, "%v", zx))
			}

		} else {
			tmp = append(tmp, k)
		}
	}

	if ioutils.Cout.VeryVerbose {
		ioutils.Cout.Seperator(60)
	}

	if len(tmp) > 0 {

		ioutils.Cout.Header("[*] Following Values Were Never Used :")
		ioutils.Cout.PrintColor(ioutils.Azure, strings.Join(tmp, "\n"))
		ioutils.Cout.Seperator(60)

	}

	if ioutils.Cout.VeryVerbose {
		ioutils.Cout.Header("[*] Implicit Declarations\n")
	}

	notdeclared := []string{}
	for k, v := range shared.DefaultRegistry.FoundVars {
		if v {

			if ioutils.Cout.VeryVerbose {
				ioutils.Cout.Value("[+] Found Implicit Declaration of %v", k)
			}

		} else {
			notdeclared = append(notdeclared, k)
		}
	}

	if ioutils.Cout.VeryVerbose {
		ioutils.Cout.Seperator(60)
	}

	if len(notdeclared) > 0 {
		ioutils.Cout.Header("\n[*] Following Variables Where Not Found")
		fatal := ""

		for k, v := range shared.DefaultRegistry.FoundVars {
			if !v {
				fatal += fmt.Sprintf("[-] Missing Declaration for %v\n", k)
			}
		}

		ioutils.Cout.PrintColor(ioutils.Red, fatal)

		os.Exit(1)
	}

}

// PrintAllCMDs : Pretty Prints All Commands
func (e *Engine) PrintAllCMDs() {

	if !ioutils.Cout.VeryVerbose {
		return
	}

	ioutils.Cout.Header("[*] All Accepted Commands\n")

	for _, v := range e.IndexedCMDs {
		ioutils.Cout.PrintColor(ioutils.Azure, "\n[+] %v", v.Comment)
		ioutils.Cout.Printf("=> %v", ioutils.Cout.GetColor(ioutils.Green, v.Raw))
		if v.Alerts != nil {
			if v.Alerts.NotifyEnabled {

				ioutils.Cout.PrintColor(ioutils.Grey, "[&] %vResult", v.Alerts.NotifyMsg)

			}
		}

	}

	ioutils.Cout.Seperator(60)

}

// NewEngine
func NewEngine() *Engine {
	z := Engine{
		CMDs:        []*shell.CMDWrap{},
		IndexedCMDs: map[string]*shell.CMDWrap{},
		ExecPyramid: [][]*scheduler.Node{},
		BlackList:   map[string]bool{},
	}

	return &z
}

// NewVolatileEngine : Engine preconfigured with defaults for volatile use
func NewVolatileEngine() *Engine {
	ioutils.Cout.Verbose = true

	z := Engine{
		CMDs:        []*shell.CMDWrap{},
		IndexedCMDs: map[string]*shell.CMDWrap{},
		ExecPyramid: [][]*scheduler.Node{},
		BlackList:   map[string]bool{},
	}

	z.ShowOutput = true

	shared.DefaultSettings.Limit = 8
	shared.DefaultSettings.Purge = true

	db.UseBBoltDB(os.TempDir(), "talos"+stringutils.RandomString(3)+".db", "TalosDefault")

	return &z
}

// cleanup : clean uneeded items from fs i.e Exports or runtime files
func cleanup() {

	exportpath := path.Join(shared.DefaultSettings.CacheDIR, shared.DefaultSettings.ProjectExportName)

	_, err := os.Stat(exportpath)
	if err != nil {
		return
	}

	// exports are runtime files {file} created and are not persistent
	// and is not part of fs cache
	os.RemoveAll(exportpath)
}
