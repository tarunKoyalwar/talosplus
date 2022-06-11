package coretest_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/tarunKoyalwar/talosplus/pkg/core"
	"github.com/tarunKoyalwar/talosplus/pkg/ioutils"
	"github.com/tarunKoyalwar/talosplus/pkg/shared"
	"github.com/tarunKoyalwar/talosplus/pkg/shell"
)

// Tests related to parsing a script
func Test_Script(t *testing.T) {

	/*
		basicscript is a very lightweight bash script
		that is expected to run on every linux
		intended to only check if all features of `talos`
		are operating correctly
	*/

	basicscript := `
	@purge = true
	@limit = 3

	// run basic linux commands to test all features

	// get short linux kernel details
	cat /proc/version | cut -d " " -f 1-3 #as:@kernel

	// save data to file and get filename
	echo @kernel{file} #as:@fileaddr

	// add random data to file 
	echo "talos file" >> @fileaddr

	// read details from file and save to buffer
	cat @fileaddr | tee @outfile #as:@alldata

	// pass data to stdin and do some filtering
	grep -i Linux #from:@alldata | cut -d " " -f 1 #as:@reqout

	// save to temp directory and filename as moriningstar
	cp @reqout{file} /tmp/morningstar
	`

	/*
		Commands are dependent on one another
		if everything executes perfectly and there is a file at /tmp/morningstar
		with data `Linux` then script parsing , scheduling and all other functions are working
		correctly
	*/

	s := core.NewScripter()
	ioutils.Cout.Verbose = true
	shell.Settings.Purge = true

	s.Compile(basicscript)

	s.Summarize()

	s.Schedule()

	s.PrintAllCMDs()

	s.Execute()

	if _, err := os.Stat("/tmp/morningstar"); err != nil {
		t.Fatalf("Test for bash script  failed\n any/all functions are not working")
	}

	dat, err2 := ioutil.ReadFile("/tmp/morningstar")
	if err2 != nil {
		t.Fatalf("Test for bash script failed\n any/all functions are not working %v", err2)
	}

	expected := string(bytes.TrimSpace(dat))

	if expected != "Linux" {
		t.Fatalf("Content Mismatch\n Some Functions are not working properly\ngot `%v` expected `Linux`", expected)
	} else {
		t.Log("Core Test Passed , talos is working properly")
	}

	// if everything is fine delete the file
	os.Remove("/tmp/morningstar")

}

// Test For Loop
func Test_For(t *testing.T) {

	ioutils.Cout.Verbose = true

	sample := ""
	for i := 0; i < 10; i++ {
		sample = sample + strconv.Itoa(i) + "\n"
	}

	shared.SharedVars.Set("@sequence", sample, true)

	bscript := `
	// Just Testing For Loop
	echo @z |  tee /tmp/@z #for:@sequence:@z #as:@forout{add}
	`

	shell.Settings.Purge = true
	s := core.NewScripter()
	s.ShowOutput = true

	s.Compile(bscript)

	s.Summarize()

	s.Schedule()

	s.Execute()

	out, er := shared.SharedVars.Get("@forout")

	if er != nil {
		panic(er)
	}

	arr := strings.Fields(out)

	sort.Strings(arr)

	fmt.Printf("GOt : %v\n", arr)

	orig := strings.Fields(sample)

	if !reflect.DeepEqual(orig, arr) {
		t.Errorf("Expected sequence not matched got %v", arr)
	} else {
		t.Logf("test of for directive working successful")
	}

}

func Test_SingleLevel(t *testing.T) {

	basicscript := `
	echo "Testing Level" #as:@slevel{add}

	echo "Level zero" #as:@slevel{add}

	tee /dev/null #from:@slevel #as:@result

	`

	s := core.NewScripter()
	ioutils.Cout.Verbose = true
	shell.Settings.Purge = true

	s.Compile(basicscript)

	s.Schedule()

	s.Execute()

	out, er1 := shell.Buffers.Get("@result")
	if er1 != nil {
		HandleErrors(er1, t, "Failed to get Exported Value of a cmd")
	}

	expected := strings.TrimSpace(`
	Testing Level
	Level zero
	`)

	if out != expected {
		t.Errorf("Failed to met expectations got\n %v", strings.TrimSpace(out))
	}

}

func HandleErrors(er error, t *testing.T, msg string, a ...any) {
	if er != nil {
		t.Logf(msg, a...)
		t.Errorf(msg, a...)
		t.Fatal(er)
	}
}
