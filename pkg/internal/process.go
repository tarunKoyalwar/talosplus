package internal

import (
	"io/ioutil"
	"path"
	"regexp"
	"strings"

	"github.com/tarunKoyalwar/talosplus/pkg/shared"
	"github.com/tarunKoyalwar/talosplus/pkg/stringutils"
)

// ProcessVariable: Process Variable and get value required by processor
func ProcessVariable(varname string) (string, error) {

	SharedVars := shared.SharedVars
	DefaultSettings := shared.DefaultSettings

	//first check if varname has anyoperation linked
	re := regexp.MustCompile("{.*}")

	matched := re.FindStringSubmatchIndex(varname)

	if len(matched) != 2 {

		got, er := SharedVars.Get(varname)
		return got, er

	}

	// Find and process operators
	// Ex: add , unique etc
	ops := varname[matched[0]+1 : matched[1]-1]
	varname = varname[:matched[0]]

	for _, v := range strings.Split(strings.TrimSpace(ops), ",") {
		if v == "file" {
			//save data of variable to file and get address of file
			addr, er1 := DefaultSettings.CreateDirectoryIfNotExist(DefaultSettings.ProjectExportName)
			if er1 != nil {
				return "", er1
			}
			reqpath := path.Join(addr, stringutils.RandomString(8))

			got, err := SharedVars.Get(varname)
			if err != nil {
				return "", err
			}

			//write to file
			er2 := ioutil.WriteFile(reqpath, []byte(got), 0644)
			if er2 != nil {
				return reqpath, er2
			}
			return reqpath, nil

			//Allow Even if Input is empty
		} else if v == "!file" {
			//save data of variable to file and get address of file
			addr, er1 := DefaultSettings.CreateDirectoryIfNotExist(DefaultSettings.ProjectExportName)
			if er1 != nil {
				return "", er1
			}
			reqpath := path.Join(addr, stringutils.RandomString(8))

			got, err := SharedVars.Get(varname)
			if err != nil {
				if err.Error() != "empty value" {
					return varname, err
				}
			}

			//write to file
			er2 := ioutil.WriteFile(reqpath, []byte(got), 0644)
			if er2 != nil {
				return reqpath, er2
			}
			return reqpath, nil
		}
	}

	return "", nil
}

// CompleteOperation : Update Value and Complete Operation such as {add} etc
func CompleteOperation(varname string, currvalue string) {
	shared.SharedVars.AcquireLock()
	defer shared.SharedVars.ReleaseLock()

	// ioutils.Cout.PrintWarning("got %v and %v\n", varname, currvalue)

	SharedVars := shared.SharedVars

	//first check if varname has anyoperation linked
	re := regexp.MustCompile("{.*}")

	matched := re.FindStringSubmatchIndex(varname)

	if len(matched) < 2 {
		// SharedVars.AcquireLock()
		SharedVars.Set(varname, currvalue, false)
		// SharedVars.ReleaseLock()
		return
	}

	ops := varname[matched[0]+1 : matched[1]-1]
	varname = varname[:matched[0]]

	for _, v := range strings.Split(strings.TrimSpace(ops), ",") {
		if v == "add" {

			//check if value exists then add itself
			if !SharedVars.Exists(varname) {
				// no such value exists
				SharedVars.Set(varname, currvalue, false)

			} else {
				// SharedVars.AcquireLock()
				//get old value and then add new value
				localvalue, _ := SharedVars.Get(varname)
				newval := localvalue + "\n" + currvalue
				SharedVars.Set(varname, newval, false)
				// SharedVars.ReleaseLock()

			}

		} else if v == "unique" {
			if !SharedVars.Exists(varname) {
				// no such value exists
				SharedVars.Set(varname, stringutils.UniqueElements(currvalue), false)
			} else {
				// SharedVars.AcquireLock()
				localvalue, _ := SharedVars.Get(varname)
				SharedVars.Set(varname, stringutils.UniqueElements(localvalue, currvalue), false)
				// SharedVars.ReleaseLock()
			}
		}
	}

}
