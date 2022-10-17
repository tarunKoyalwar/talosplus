package shell

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/tarunKoyalwar/talosplus/pkg/internal"
	"github.com/tarunKoyalwar/talosplus/pkg/shared"
	"github.com/tarunKoyalwar/talosplus/pkg/stringutils"
)

// Settings :Pointer to default settings
var Settings *shared.Settings = shared.DefaultSettings

// Buffers : Pointer to shared values
var Buffers *shared.Shared = shared.SharedVars

// CMDWrap : A Wrapper Around SimpleCMD provides
// features like caching , exports etc
type CMDWrap struct {
	CMD *SimpleCMD // access to basic cmd

	Comment  string // Human Readable Comment
	Raw      string //Raw Command to be passed
	CacheKey string // Unique Hash of Command to Cache
	UID      string // A Temporary UID of Raw Command

	ExportAs       string // Export as variable name
	ExportFromFile string // IF @outfile is used then export data from here

	IsForLoop bool //Does the command have for loop
	IsInvalid bool //Will Be Valid Only if Dependency is Not Empty
	Ignore    bool // Ignore Flag Does not print output

	Alerts *Notifications // Access to Notifications if enabled

	CauseofFailure []string // Cause of Failure to execute command
}

// Process : Process command and fill values
func (c *CMDWrap) Process() {
	rawsplit := stringutils.SplitAtSpace(c.Raw)

	if strings.Contains(c.Raw, "#for:") {
		c.IsForLoop = true
		return
	}

	// missed source [temp fix] [Patched]
	// if c.CMD == nil {
	// 	c.CMD = &SimpleCMD{}
	// }

	// cmd after replacing and processing directives
	filtered := []string{}

	for _, v := range rawsplit {
		// fmt.Println(v)
		if strings.HasPrefix(v, "#as:") {
			val := strings.TrimLeft(v, "#as:")
			c.ExportAs = val

		} else if strings.HasPrefix(v, "#from:") {
			key := strings.TrimLeft(v, "#from:")
			val, er1 := Buffers.Get(key)
			if er1 != nil {
				c.addReason("\tfailed to fetch env value for %v\n no reason to execute", val)
			}
			c.CMD.UseStdin(val)

		} else if strings.HasPrefix(v, "#dir") {
			val := strings.TrimLeft(v, "#dir:")
			c.CMD.DIR = val
		} else if v == "#ignore" || strings.Contains(v, "#ignore") {
			// do no print output of the command
			c.Ignore = true
		} else if strings.HasPrefix(v, "@env:") {
			val := strings.TrimLeft(v, "@env:")
			envalue := os.Getenv(val)
			if envalue == "" {
				c.addReason("\tfailed to fetch env value for %v\n no reason to execute", val)
			}
			filtered = append(filtered)

		} else if strings.HasPrefix(v, "@") {

			// Sanitize v
			v = stringutils.ExtractVar(v)

			if v == "@outfile" {

				addr, er1 := shared.DefaultSettings.CreateDirectoryIfNotExist(shared.DefaultSettings.ProjectExportName)
				if er1 != nil {
					c.addReason("\tfailed to create directory no reason to execute %v", er1)
					continue
				}
				c.ExportFromFile = path.Join(addr, filtered[0]+"-texport-"+stringutils.RandomString(8))
				filtered = append(filtered, c.ExportFromFile)
			} else if v == "@tempfile" {
				addr, er1 := shared.DefaultSettings.CreateDirectoryIfNotExist(shared.DefaultSettings.ProjectExportName)
				if er1 != nil {
					c.addReason("\tfailed to create directory no reason to execute %v", er1)
					continue
				}
				tmploc := path.Join(addr, filtered[0]+"-texport-"+stringutils.RandomString(8))
				filtered = append(filtered, tmploc)

			} else {
				resp, er := internal.ProcessVariable(v)
				if er != nil {
					c.addReason("\t%v", er.Error())
					continue
				}

				filtered = append(filtered, resp)
			}

		} else {
			filtered = append(filtered, v)
		}

	}

	c.CMD.Cmdsplit = filtered

	c.genCacheKey()

}

// Export : Setup Export
func (c *CMDWrap) Export() {
	// Check Export Type
	// Use a tempfile as output

	if c.ExportFromFile != "" {
		bin, err := ioutil.ReadFile(c.ExportFromFile)
		if err != nil {
			c.addReason("Temp File was not created by command %v\n%v", c.Raw, err)
			return
		}

		fileval := string(bin)

		internal.CompleteOperation(c.ExportAs, fileval)

		c.Alerts.Notify(fileval)

	} else {
		//use stdout as output
		// ioutils.Cout.PrintInfo("got cout as %v for %v", c.CMD.COutStream.String(), c.UID)

		if c.ExportAs != "" {
			internal.CompleteOperation(c.ExportAs, c.CMD.COutStream.String())
		}

		// fmt.Printf("calling notify")
		c.Alerts.Notify(c.CMD.COutStream.String())

	}

}

// Execute : Runs Command If data is Not Present in Cache
func (c *CMDWrap) Execute() error {

	//don't run if prechecks failed

	if len(c.CauseofFailure) > 0 {
		er := fmt.Errorf("%v", strings.Join(c.CauseofFailure, "\n"))
		return er
	}

	var runerror error

	if !shared.DefaultSettings.Purge {
		err := c.cacheIn()
		if err != nil {
			// ioutils.Cout.PrintWarning("Catched Data Was Not Found %v\n", err.Error())
		} else {
			return nil
		}
	}

	// ioutils.Cout.PrintWarning("running cmd %v with cmd %v\n", c.UID, c.Raw)
	// ioutils.Cout.PrintWarning("running %v from cmdwrap", c.UID)
	runerror = c.CMD.Run()

	if runerror == nil {
		c.Export()

		//always cache
		c.cacheOut()
	}

	return runerror
}

// cacheIn : Check Cache And Import Output If Exists
func (c *CMDWrap) cacheIn() error {
	wdir := path.Join(shared.DefaultSettings.CacheDIR, shared.DefaultSettings.ProjectName)
	cpath := path.Join(wdir, c.CacheKey)

	// if _, err := os.StartProcess()
	_, err := os.Stat(cpath)

	if err == nil {
		bin, _ := ioutil.ReadFile(cpath)
		_, err = c.CMD.COutStream.Write(bin)
		internal.CompleteOperation(c.ExportAs, string(bin))
		if err != nil {
			return err
		}
	}

	return err
}

// cacheOut : Cache Output stream to be used later
func (c *CMDWrap) cacheOut() error {

	wdir := path.Join(shared.DefaultSettings.CacheDIR, shared.DefaultSettings.ProjectName)

	//check if cache dir exists and dir name
	_, err := os.Stat(wdir)

	if err != nil {

		//Create New DIrectory
		err := os.Mkdir(wdir, 0755)
		if err != nil {
			return err
		}
	}

	bin := c.CMD.COutStream.Bytes()
	//Save Out to File With HashName rw-r--r--
	err = ioutil.WriteFile(path.Join(wdir, c.CacheKey), bin, 0644)
	// g.COutStream.Write(bin)

	return err

}

// genCacheKey
func (c *CMDWrap) genCacheKey() {
	//copy command
	tarr := []string{}

	for _, v := range c.CMD.Cmdsplit {
		//if it is a tempory file skip it from hash
		if !strings.Contains(v, shared.DefaultSettings.CacheDIR) {
			tarr = append(tarr, v)
		}
	}

	//sort to avoid duplicates
	sort.Strings(tarr)

	//Lets use # as separator
	suffix := strings.Join(tarr, "#")

	data := []byte(suffix)

	bin := md5.Sum(data)

	c.CacheKey = c.CMD.Cmdsplit[0] + "-" + hex.EncodeToString(bin[:])
}

// addReason : It wraps multiple errors if case of a failure
func (c *CMDWrap) addReason(format string, a ...any) {
	c.IsInvalid = true

	if c.CauseofFailure == nil {
		c.CauseofFailure = []string{}
	}

	c.CauseofFailure = append(c.CauseofFailure, fmt.Sprintf(format+"\n", a...))

}

// Disolve : Disolves Command Into Multiple Commands
func (c *CMDWrap) Disolve() ([]CMDWrap, error) {
	cmdarr := []CMDWrap{}
	//get for statement
	getfrom := ""
	dynvar := ""
	// filtered := []string{}

	for _, v := range strings.Split(c.Raw, " ") {
		if strings.Contains(v, "#for:") {
			datx := strings.TrimLeft(v, "#for:")
			splitdat := strings.Split(datx, ":")
			if len(splitdat) == 2 {
				getfrom = splitdat[0]
				dynvar = splitdat[1]
			} else {
				return cmdarr, fmt.Errorf("malformed for loop check syntax")

			}
		}
	}

	value, er1 := internal.ProcessVariable(getfrom)
	if er1 != nil {
		return cmdarr, er1
	}
	value = strings.TrimSpace(value)
	if len(value) < 1 {
		return cmdarr, fmt.Errorf("Variable Has No data . Hence No need for execution of this command")
	}

	if dynvar == "" {
		return cmdarr, fmt.Errorf("Something went wrong this was not supposed to happen %v", c.Raw)

	}

	for _, v := range strings.Split(value, "\n") {
		req := strings.TrimSpace(v)
		if v != "" {

			newcmd := strings.ReplaceAll(c.Raw, "#for:"+getfrom+":"+dynvar, "")

			newcmd = strings.ReplaceAll(newcmd, dynvar, req)

			//replace for loop statement
			// tmparr := strings.Split(newcmd)

			wrap := NewCMDWrap(newcmd, c.Comment+" Loop")

			cmdarr = append(cmdarr, wrap)
		}
	}

	return cmdarr, nil
}

// NewCMDWrap
func NewCMDWrap(newcmd string, comment string) CMDWrap {

	nf, rawcmd := NewNotification(newcmd, comment)
	xz := &SimpleCMD{}

	tcmd := internal.NewCommand(rawcmd, comment)
	tcmd.Resolve()

	wrap := CMDWrap{
		CMD:     xz,
		Raw:     rawcmd,
		Comment: tcmd.Comment,
		Alerts:  nf,
		UID:     tcmd.UID,
	}

	return wrap
}
