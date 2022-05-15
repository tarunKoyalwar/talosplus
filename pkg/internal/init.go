package internal

import (
	"regexp"
	"strings"

	"github.com/tarunKoyalwar/talosplus/pkg/shared"
	"github.com/tarunKoyalwar/talosplus/pkg/stringutils"
)

// ParseScript : This is entry point and all parsing is done from here
func ParseScript(s string) []Command {

	lines := stringutils.Split(s, '\n')

	//convert multiline cmds to single line

	singlelines := parseMultiLine(lines)

	cmdarr := parseScriptLines(singlelines)

	// Parse Settings from Global Vars
	vars := shared.SharedVars.GetGlobalVars()
	shared.DefaultSettings.ParseSettings(vars)

	//all parsing done

	return cmdarr

}

// parseCommands : Parse Commands and Load Settings + Env Variables
func parseScriptLines(arr []string) []Command {

	cmds := []Command{}

	//Get Comment Commands and Global Variables
	for i := 0; i < len(arr); i++ {

		v := arr[i]

		// Check if Line has any comments
		z := strings.Index(v, "//")

		if z != -1 {
			//remove anything after //
			v = strings.TrimSpace(v[:z])

			// If line only had a comment skip
			if v == "" {
				continue
			}
		}

		//Parse Variables at top of script
		if strings.HasPrefix(v, "@") {
			splitd := stringutils.Split(v, '=')

			if len(splitd) == 2 {
				key := strings.TrimSpace(splitd[0])
				value := strings.TrimSpace(splitd[1])

				// remove outer quotes from value
				value = strings.Trim(value, "\"")

				// Add variable to shared vars

				// Malformed variables are ignored
				shared.SharedVars.Set(key, value, true)

			}
			continue
		}

		x := NewCommand(v, "")

		// Extract any comment given to command by
		// parsing previous line
		if i-1 >= 0 {

			lastline := strings.TrimSpace(arr[i-1])

			if strings.HasPrefix(lastline, "//") {
				x.Comment = strings.TrimLeft(lastline, "//")
			}
		}

		cmds = append(cmds, x)

	}

	return cmds

}

// arseMutliLine : Convert MultLine Commands to single line
func parseMultiLine(arr []string) []string {

	/*
	   Changes a multi line command to single line

	   Ex :

	   for i in {1..9}; do; \
	   printf "no $i"\
	   done;

	    for i in {1..9}; do; printf "no $i"; done;
	*/

	var sb strings.Builder

	newarr := []string{}

	for _, v := range arr {

		v = strings.TrimSpace(v)

		if strings.HasSuffix(v, `\`) {
			sb.WriteString(v[:len(v)-1])
		} else {
			// if buffer is not empty
			if sb.Len() != 0 {
				newarr = append(newarr, sb.String()+" "+v)
				sb.Reset()
			} else {

				newarr = append(newarr, v)
			}
		}

	}

	return newarr

}

func GetVariableName(varname string) string {
	//first check if varname has anyoperation linked
	re := regexp.MustCompile("{.*}")

	matched := re.FindStringSubmatchIndex(varname)

	if len(matched) != 2 {

		return varname
	}

	// Find and process operators
	// Ex: add , unique etc
	varname = varname[:matched[0]]

	return varname

}
