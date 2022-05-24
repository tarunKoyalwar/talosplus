package shell

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/tarunKoyalwar/talosplus/pkg/ioutils"
)

// BashBuiltInCmds : Bash Interpreter Built in commands like for ,if etc
var BashBuiltInCmds map[string]bool = map[string]bool{
	"for": true,
	"if":  true,
}

// SimpleCMD : Wrapper for std lib exec.CMD With Some other features
type SimpleCMD struct {
	Cmdsplit   []string     //Command split using space
	CErrStream bytes.Buffer // Error Stream
	COutStream bytes.Buffer // Stdout Stream
	CPipeIn    bytes.Buffer // Pass Data Using Pipe
	DIR        string       // DIR in which CMD should be run
	Failed     bool         // If Command Failed due to a reason
}

// CheckInstall : Check if Program Binary Exist Either in Path or DIR
func (s *SimpleCMD) CheckInstall() (string, error) {

	//Check if command is built in shell interpreter
	if BashBuiltInCmds[s.Cmdsplit[0]] {
		return s.Cmdsplit[0], nil
	}

	//check for cmd in
	cpath, err := exec.LookPath(s.Cmdsplit[0])
	if err == nil {
		return cpath, err
	}
	if s.DIR != "" {
		// Program Is Not Installed In Path
		p := path.Join(s.DIR, s.Cmdsplit[0])
		if _, err := os.Stat(p); err == nil {

			return p, err
		}
	}

	return "", fmt.Errorf("%v Not Found in path or dir", s.Cmdsplit[0])

}

// Run : Run Command Get Error
func (s *SimpleCMD) Run() error {
	// ioutils.Cout.PrintWarning("running %v from cmdwrap", s.Cmdsplit)

	// Just a precaution to avoid FP in results
	_, er := s.CheckInstall()

	if er != nil {
		ioutils.Cout.PrintWarning("did not find %v", er)
		s.Failed = true
		return er
	}

	var c exec.Cmd

	//Its better to use sh -c Instead of Normal
	//Doing this allows piping and all features with it
	//Without that we can only run 1 command

	if s.DIR != "" {

		c = *exec.Command("sh", "-c", strings.Join(s.Cmdsplit, " "))
		c.Dir = s.DIR
	} else {
		c = *exec.Command("sh", "-c", strings.Join(s.Cmdsplit, " "))
	}

	c.Stdout = &s.COutStream
	c.Stderr = &s.CErrStream

	if s.CPipeIn.Len() > 1 {

		stdin, err := c.StdinPipe()
		if err != nil {
			ioutils.Cout.PrintWarning("[error] %v", err)
			s.Failed = true
		}
		io.WriteString(stdin, s.CPipeIn.String())
		stdin.Close()
	}

	err := c.Run()

	if s.COutStream.Len() > 0 || s.CErrStream.Len() > 0 {
		err = nil
	}

	// Identify Cases When Command returns error or panics
	if s.CErrStream.Len() > 0 {
		trimmedcout := strings.TrimSpace(s.COutStream.String())

		if len(trimmedcout) == 0 {
			dat := s.CErrStream.String()

			dat = strings.TrimSpace(dat)

			if len(dat) > 0 {
				err = fmt.Errorf("%v", dat)
			}

		}

	}

	if err != nil {
		s.Failed = true
		return err
	}

	return err

}

// UseStdin : Uses provided string as stdin
func (s *SimpleCMD) UseStdin(d string) {

	//skipping errs since stdin usually is safe
	s.CPipeIn.WriteString(d)
}
