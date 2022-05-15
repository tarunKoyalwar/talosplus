package internal

import (
	"regexp"
	"strings"

	"github.com/tarunKoyalwar/talosplus/pkg/shared"
	"github.com/tarunKoyalwar/talosplus/pkg/stringutils"
)

//This will contain actions for each directives

// ProcessForDirective : Process For Directive
func ProcessForDirective(uid string, cmd string) string {
	/*
		#for
		For directive works similar to for each loop

		[+] Syntax
		#for:@arr:@z
		Here

		@z : Loop Variable
		@arr : Array to Loop Over


		[+] Steps
		1. identify array variable(@arr) and loop variable(@z)
		2. register uid as dependent of @arr in the registry in shared package
		3. remove all loop variables from cmd
		4. return cmd to process other directives
	*/

	returncmd := ""
	loopvar := ""

	for _, v := range stringutils.SplitAtSpace(cmd) {

		if strings.Contains(v, "#for") {
			datx := strings.TrimLeft(v, "#for:")
			splitfor := strings.Split(datx, ":")

			if len((splitfor)) == 2 {
				arrayvar := splitfor[0]
				loopvar = splitfor[1]

				arrayvar = GetVariableName(arrayvar)
				loopvar = GetVariableName(loopvar)

				// register uid as dependent
				shared.DefaultRegistry.Registerdep(uid, arrayvar)
			}
		} else {
			returncmd += v + " "
		}
	}

	// remove all loop vars
	returncmd = strings.ReplaceAll(returncmd, loopvar, "")

	return returncmd

}

// ProcessFromDirective : Process From Directive
func ProcessFromDirective(word string, uid string) {
	/*
		#from
		Use given variable data as stdin for this command
		and mark dependency

		[+] Syntax
		#from:@gvar

		Here

		@gvar : data source for this command


		[+] Steps
		1. Just mark this cmd as dependent of that variable

	*/
	gvar := strings.TrimLeft(word, "#from:")

	// register uid as dependent
	shared.DefaultRegistry.Registerdep(uid, gvar)
}

// ProcessAsDirective : Process As Directive
func ProcessAsDirective(word string, uid string) {
	/*
		#as
		marks command as exporter of this given variable

		[+] Syntax
		#as:@somevar

		Here

		@somevar : output of this command will be exported to @somevar


		[+] Steps
		1. Just mark this cmd as Exporter of that variable

	*/

	word = GetVariableName(word)

	somevar := strings.TrimLeft(word, "#as:")

	// register uid as exporter
	shared.DefaultRegistry.RegisterExport(uid, somevar)

}

// ProcessVariables : Process Variables @something
func ProcessVariables(word string, uid string) {
	/*
		[+] Syntax
		@varname or @outfile

		Variables have two different categories

		Here

		@outfile :
		There are many commands that store the filtered and required output
		to a file instead of printing it on stdout. In such cases using
		@outfile will use content of file as output instead of stdout

		@varname : use value of the variable
		Ex : curl -v @exploiturl


		[+] Steps
		1. Just mark this cmd as dependent of that variable

	*/
	word = GetVariableName(word)

	if word == "@outfile" {
		//ignore nothing to do
		return
	}

	// register uid as exporter
	shared.DefaultRegistry.Registerdep(uid, word)
}

// RemoveOperators : Removes Ops from variables
func RemoveOperators(raw string) string {
	/*
		variables can have operators like

		@gvar{file} @gvar{unique} @gvar{add}


		remove all operators from command using regex

	*/

	re := regexp.MustCompile("{[^(}|@)]*}")

	return re.ReplaceAllLiteralString(raw, "")
}
