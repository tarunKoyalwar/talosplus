package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/tarunKoyalwar/talosplus/pkg/gopool"
	"github.com/tarunKoyalwar/talosplus/pkg/internal"
	"github.com/tarunKoyalwar/talosplus/pkg/ioutils"
	"github.com/tarunKoyalwar/talosplus/pkg/mongodb"
	"github.com/tarunKoyalwar/talosplus/pkg/scheduler"
	"github.com/tarunKoyalwar/talosplus/pkg/shared"
	"github.com/tarunKoyalwar/talosplus/pkg/shell"
)

// CRUX : High Level Struct Which takes shell script and executes
type Scripter struct {
	CMDs        []*shell.CMDWrap
	IndexedCMDs map[string]*shell.CMDWrap
	ExecPyramid [][]*scheduler.Node
	Backup      []ioutils.CSave
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

	// ioutils.Cout.PrintInfo("Resume Kicked In ")

	//these are needless commands that should not run at all
	needless := map[string]bool{}
	for _, v := range arr {
		addresses := shared.DefaultRegistry.VarAddress[v]
		for _, addr := range addresses {
			needless[addr] = true
		}
	}

	ioutils.Cout.Printf("[+] Skipping Following commands\n")

	for k := range needless {
		val := s.IndexedCMDs[k].Comment
		if val == "" {
			val = s.IndexedCMDs[k].Raw
		}
		ioutils.Cout.Printf("%v : %v", k, val)
	}

	ioutils.Cout.DrawLine(30)

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

	p := gopool.NewPool(shared.DefaultSettings.Limit)
	//this is a must
	defer p.Release()

	p.HandleError = func(er error) {
		if er != nil {
			ioutils.Cout.PrintWarning(er.Error())
		}
	}

	//Execute This after Every SuccessFul COmmand Completion
	p.OnCompletion = func(resp gopool.JobResponse) {
		// ioutils.Cout.PrintInfo("%v was executed \n", resp.Uid)
		//check if job was successful
		if resp.Err == nil {

			b := ioutils.CSave{
				UID: resp.Uid,
			}

			instance := s.IndexedCMDs[resp.Uid]
			if instance != nil {
				ioutils.Cout.Printf("[$] %v Executed Successfully", instance.Comment)

				// WIll show output of commands
				if s.ShowOutput {
					if !instance.Ignore {
						ioutils.Cout.Printf("[+] %v\n", instance.Raw)
						if instance.ExportFromFile == "" {
							ioutils.Cout.Printf("%v", instance.CMD.COutStream.String())
						} else {
							dat, _ := ioutil.ReadFile(instance.ExportFromFile)
							ioutils.Cout.Printf("%v", string(dat))
						}
					}
				}

				b.Comment = instance.Comment

				b.CacheKey = instance.CacheKey

				varname := instance.ExportAs

				//Just Extract Correct Value
				re := regexp.MustCompile("{.*}")

				matched := re.FindStringSubmatchIndex(varname)
				if len(matched) < 2 {
					b.Output, _ = shared.SharedVars.Get(varname)
				} else {
					varname = varname[:matched[0]]
					b.Output, _ = shared.SharedVars.Get(varname)
				}

			}

			s.Backup = append(s.Backup, b)

		} else {
			instance := s.IndexedCMDs[resp.Uid]
			if instance != nil {
				ioutils.Cout.Printf("[Failed] %v\n", strings.Join(instance.CMD.Cmdsplit, " "))
			}

			ioutils.Cout.Printf("[-] %v responded with error %v\n", resp.Uid, resp.Err)
		}
	}

	count := 0

	for _, v := range s.ExecPyramid {
		ioutils.Cout.Printf("[^_^] Executing Level %v Commands\n", count)
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

				// fmt.Printf("before process")
				// process the command
				c.Process()
				// fmt.Printf("after process")

				if !c.IsForLoop {
					queue = append(queue, *c)
					// fmt.Printf("added to queue")
				} else {
					dissolved, er := c.Disolve()

					if er != nil {
						ioutils.Cout.Printf("[-] %v Will Not be Executed because :%v\n", c.Comment, er)
					} else {
						for _, tinstance := range dissolved {

							dx := tinstance

							dx.Process()

							// fmt.Printf("%v : %v\n", dx.UID, dx.Raw)
							// also add them to indexdb
							s.IndexedCMDs[dx.UID] = &dx
							s.CMDs = append(s.CMDs, &dx)
							queue = append(queue, dx)
						}
					}

				}

			}
		}

		for _, c := range queue {
			//All Checks Passed
			if !c.IsInvalid {

				ioutils.Cout.PrintInfo("(*) Scheduled... %v", c.Raw)
				p.AddJobWithId(s.IndexedCMDs[c.UID], c.UID)
				// fmt.Println(unsafe.Sizeof(c))

			} else {
				ioutils.Cout.Printf("[-] %v Will Not be Executed because :\n%v", c.Comment, strings.Join(c.CauseofFailure, "\n"))
			}
		}

		p.Wait()
		ioutils.Cout.DrawLine(32)

	}
	//All Jobs Assigned
	p.Done()

	//cleanup

	defer cleanup()

}

// Summarize : Summarizes All Script Data
func (s *Scripter) Summarize() {

	ioutils.Cout.Printf("\n[*] Used Explicit declared Variables\n")
	gvars := shared.SharedVars.GetGlobalVars()

	for k, v := range gvars {
		ioutils.Cout.Printf("%-16v : %v", k, v)
	}

	ioutils.Cout.DrawLine(30)

	tmp := []string{}

	ioutils.Cout.Printf("\n[*] Generated UIDs For Commands\n")
	for k, v := range s.IndexedCMDs {
		identifier := v.Comment
		if identifier == "" {
			identifier = v.Raw
		}
		ioutils.Cout.Printf("%v : %v", k, identifier)
	}

	ioutils.Cout.DrawLine(30)

	ioutils.Cout.Printf("[*] Dependencies Found\n")

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

			ioutils.Cout.Printf("[ %v ] Will be Executed After [%v]\n", strings.Join(requiredby, " , "), strings.Join(providers, " , "))

		} else {
			tmp = append(tmp, k)
		}
	}

	ioutils.Cout.DrawLine(30)

	if len(tmp) > 0 {

		ioutils.Cout.Printf("[*] Following Values Were Never Used : \n%v \n", strings.Join(tmp, "\n"))
		ioutils.Cout.DrawLine(30)

	}

	ioutils.Cout.Printf("[*] Implicit Declarations\n")
	notdeclared := []string{}
	for k, v := range shared.DefaultRegistry.FoundVars {
		if v {
			ioutils.Cout.Printf("[+] Found Implicit Declaration of %v", k)
		} else {
			notdeclared = append(notdeclared, k)
		}
	}

	ioutils.Cout.DrawLine(30)

	if len(notdeclared) > 0 {
		ioutils.Cout.Printf("\n[*] Following Variables Where Not Found")
		fatal := ""

		for k, v := range shared.DefaultRegistry.FoundVars {
			if !v {
				fatal += fmt.Sprintf("[-] Missing Declaration for %v\n", k)
			}
		}

		ioutils.Cout.Printf(fatal)

		os.Exit(1)

		ioutils.Cout.DrawLine(30)
	}

}

// PrintAllCMDs : Pretty Prints All Commands
func (e *Scripter) PrintAllCMDs() {

	ioutils.Cout.Printf("[*] All Accepted Commands\n")

	for _, v := range e.IndexedCMDs {
		ioutils.Cout.Printf("\n[+] %v", v.Comment)
		ioutils.Cout.Printf("=> %v", v.Raw)
		if v.Alerts != nil {
			if v.Alerts.NotifyEnabled {

				ioutils.Cout.Printf("[&] %vResult", v.Alerts.NotifyMsg)

			}
		} else {
			// fmt.Printf("nil for %v\n", v.Raw)
		}

	}

	ioutils.Cout.DrawLine(30)

}

func (s *Scripter) Export(filename string) {
	a := ioutils.AllSave{
		Commands: s.Backup,
		Exports:  shared.SharedVars.GetGlobalVars(),
	}

	_ = a
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
