package core

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/tarunKoyalwar/talosplus/pkg/internal"
	"github.com/tarunKoyalwar/talosplus/pkg/ioutils"
	"github.com/tarunKoyalwar/talosplus/pkg/mongodb"
	"github.com/tarunKoyalwar/talosplus/pkg/scheduler"
	"github.com/tarunKoyalwar/talosplus/pkg/shared"
	"github.com/tarunKoyalwar/talosplus/pkg/shell"
	"github.com/tarunKoyalwar/talosplus/pkg/workshop"
)

// CRUX : High Level Struct Which takes shell script and executes
type Scripter struct {
	CMDs        []*shell.CMDWrap
	IndexedCMDs map[string]*shell.CMDWrap
	ExecPyramid [][]*scheduler.Node
	ShowOutput  bool
	BlackList   map[string]bool
}

func (s *Scripter) fillindex() {
	for _, v := range s.CMDs {
		if v.Alerts == nil {
			fmt.Println("nil  in array")
		}
		s.IndexedCMDs[v.UID] = v
	}
}

// Compile : Self Explainatory
func (s *Scripter) Compile(shellscript string) {

	cmdarr := internal.ParseScript(shellscript)

	for _, v := range cmdarr {
		// create cmdwraps from basic commands
		wrap := shell.NewCMDWrap(v.Raw, v.Comment)
		s.CMDs = append(s.CMDs, &wrap)

	}

	s.fillindex()

}

// Schedule : Schedule programs by analyzing dependencies
func (s *Scripter) Schedule() {

	t := scheduler.NewScheduler()

	// remove dependency if variable is already present in data
	arr := s.filterCompleted()

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
		val := s.IndexedCMDs[k].Comment
		if val == "" {
			val = s.IndexedCMDs[k].Raw
		}
		ioutils.Cout.PrintColor(ioutils.Azure, "%v : %v", k, val)
	}

	ioutils.Cout.Seperator(60)

	// autobalance dependcies of blacklisted nodes
	t.BlackListed = needless

	for _, v := range s.CMDs {
		t.AddNode(v.UID, v.Comment)
	}

	t.Run()

	s.ExecPyramid = t.ExecPyramid

}

// filtercompleted : Filter Commands that are already completed
func (s *Scripter) filterCompleted() []string {
	updatedvars := []string{}
	if mongodb.MDB == nil {
		return updatedvars
	}

	//get all runtime vars
	z, err := shared.LoadAllRuntimeVars()
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
			if s.BlackList[k] {
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

// Execute : Will Execute Script In orderly Fashion
func (s *Scripter) Execute() {

	count := 0

	for _, v := range s.ExecPyramid {
		ioutils.Cout.Header("[^_^] Executing Level %v Commands\n", count)
		count += 1

		queue := []shell.CMDWrap{}

		//check if it for loop and dissolve for each level
		for _, ftest := range v {

			uid := ftest.UID
			c := s.IndexedCMDs[uid]

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

							s.IndexedCMDs[dx.UID] = &dx
							s.CMDs = append(s.CMDs, &dx)
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
				finalqueue = append(finalqueue, s.IndexedCMDs[c.UID])
				// fmt.Println(unsafe.Sizeof(c))

			} else {
				ioutils.Cout.Printf("[-] %v Will Not be Executed because :\n%v", c.Comment, ioutils.Cout.GetColor(ioutils.Azure, strings.Join(c.CauseofFailure, "\n")))
			}
		}

		workshop.ExecQueue(finalqueue, shared.DefaultSettings.Limit, s.ShowOutput)
		ioutils.Cout.Seperator(60)

	}
	//cleanup

	defer cleanup()

}

// Summarize : Summarizes All Script Data
func (s *Scripter) Summarize() {

	if ioutils.Cout.Verbose {
		//Only in verbose Mode

		ioutils.Cout.Header("\n[*] Parsed Settings\n")
		ioutils.Cout.Value("%-16v : %v", "Purge Cache", shared.DefaultSettings.Purge)
		ioutils.Cout.Value("%-16v : %v", "Concurrency", shared.DefaultSettings.Limit)
		ioutils.Cout.Value("%-16v : %v", "ProjectName", shared.DefaultSettings.ProjectName)
		ioutils.Cout.Value("%-16v : %v", "CacheDir", shared.DefaultSettings.CacheDIR)
		ioutils.Cout.Value("%-16v : %v", "Verbose", ioutils.Cout.Verbose)

		ioutils.Cout.Seperator(60)

	}

	ioutils.Cout.Header("\n[*] Used Explicit declared Variables\n")
	gvars := shared.SharedVars.GetGlobalVars()

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

	tmp := []string{}

	if ioutils.Cout.Verbose {
		// Only in Verbose Mode

		ioutils.Cout.Header("\n[*] Generated UIDs For Commands\n")
		for k, v := range s.IndexedCMDs {
			identifier := v.Comment
			if identifier == "" {
				identifier = v.Raw
			}
			ioutils.Cout.Value("%v : %v", k, identifier)
		}

		ioutils.Cout.Seperator(60)
	}

	if ioutils.Cout.Verbose {
		ioutils.Cout.Header("[*] Dependencies Found\n")
	}

	for k, v := range shared.DefaultRegistry.Dependents {
		if len(v) != 0 {

			vaddress := shared.DefaultRegistry.VarAddress[k]
			providers := []string{}

			for _, addr := range vaddress {
				if s.IndexedCMDs[addr] != nil {
					providers = append(providers, s.IndexedCMDs[addr].Comment)
				} else {
					providers = append(providers, k)
				}
			}

			requiredby := []string{}

			for _, uid := range v {
				if s.IndexedCMDs[uid] != nil {
					inst := s.IndexedCMDs[uid]
					requiredby = append(requiredby, inst.Comment)
				}
			}

			if ioutils.Cout.Verbose {
				// Print Only in Verbose Mode

				zx := fmt.Sprintf("[ %v ] Will be Executed After [%v]\n", strings.Join(requiredby, " , "), strings.Join(providers, " , "))
				ioutils.Cout.Printf("%v", ioutils.Cout.GetColor(ioutils.LightGreen, "%v", zx))
			}

		} else {
			tmp = append(tmp, k)
		}
	}

	if ioutils.Cout.Verbose {
		ioutils.Cout.Seperator(60)
	}

	if len(tmp) > 0 {

		ioutils.Cout.Header("[*] Following Values Were Never Used :")
		ioutils.Cout.PrintColor(ioutils.Azure, strings.Join(tmp, "\n"))
		ioutils.Cout.Seperator(60)

	}

	if ioutils.Cout.Verbose {
		ioutils.Cout.Header("[*] Implicit Declarations\n")
	}

	notdeclared := []string{}
	for k, v := range shared.DefaultRegistry.FoundVars {
		if v {

			if ioutils.Cout.Verbose {
				ioutils.Cout.Value("[+] Found Implicit Declaration of %v", k)
			}

		} else {
			notdeclared = append(notdeclared, k)
		}
	}

	if ioutils.Cout.Verbose {
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

		ioutils.Cout.Seperator(60)
	}

}

// PrintAllCMDs : Pretty Prints All Commands
func (e *Scripter) PrintAllCMDs() {

	ioutils.Cout.Header("[*] All Accepted Commands\n")

	for _, v := range e.IndexedCMDs {
		ioutils.Cout.PrintColor(ioutils.Azure, "\n[+] %v", v.Comment)
		ioutils.Cout.Printf("=> %v", ioutils.Cout.GetColor(ioutils.Green, v.Raw))
		if v.Alerts != nil {
			if v.Alerts.NotifyEnabled {

				ioutils.Cout.PrintColor(ioutils.Grey, "[&] %vResult", v.Alerts.NotifyMsg)

			}
		} else {
			// fmt.Printf("nil for %v\n", v.Raw)
		}

	}

	ioutils.Cout.Seperator(60)

}

func NewScripter() *Scripter {
	z := Scripter{
		CMDs:        []*shell.CMDWrap{},
		IndexedCMDs: map[string]*shell.CMDWrap{},
		ExecPyramid: [][]*scheduler.Node{},
		BlackList:   map[string]bool{},
	}

	return &z
}

// cleanup : clean uneeded items from fs
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
