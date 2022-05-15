package internal

import (
	"crypto/md5"
	"encoding/hex"
	"sort"
	"strings"

	"github.com/tarunKoyalwar/talosplus/pkg/stringutils"
)

// Command : This struct holds raw command
type Command struct {
	Raw     string
	Comment string
	UID     string
}

// This will generate UID's for each command
func (c *Command) GenUID() {

	cmdarr := stringutils.SplitAtSpace(c.Raw)

	//sort to avoid duplicates
	sort.Strings(cmdarr)

	//Lets use # as separator
	suffix := strings.Join(cmdarr, "#")

	data := []byte(suffix)

	bin := md5.Sum(data)

	c.UID = hex.EncodeToString(bin[:])

}

// Resolve : Resolve all directives from command
func (c *Command) Resolve() {
	/*
		Each directive and variables in command have different meanings

		By resolving these directives and variables

		system would be able to create a dependency graph

		do some static analysis and stuff

		In this step command remains as it is directives/ variables are not replaced
		at this stage . that only happens at runtime


	*/

	// Remove all operators from cmd

	tempcmd := RemoveOperators(c.Raw)

	// If command has #for directive process it first

	if strings.Contains(tempcmd, "#for") {
		tempcmd = ProcessForDirective(c.UID, c.Raw)
	}

	for _, word := range stringutils.SplitAtSpace(tempcmd) {

		if strings.Contains(word, "#from") {
			ProcessFromDirective(word, c.UID)
		} else if strings.Contains(word, "#as") {
			ProcessAsDirective(word, c.UID)
		} else if strings.Contains(word, "@") || word == "@outfile" {
			ProcessVariables(word, c.UID)
		}
	}

}

func NewCommand(raw string, comment string) Command {

	z := Command{
		Raw:     raw,
		Comment: comment,
	}

	z.GenUID()

	return z
}
