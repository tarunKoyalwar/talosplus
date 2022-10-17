package coretest_test

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/tarunKoyalwar/talosplus/pkg/core"
	"github.com/tarunKoyalwar/talosplus/pkg/shared"
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

	s := core.NewVolatileEngine()

	s.Compile(basicscript)

	s.Evaluate()

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

	s := core.NewVolatileEngine()

	sample := ""
	for i := 0; i < 10; i++ {
		sample = sample + strconv.Itoa(i) + "\n"
	}

	shared.SharedVars.Set("@sequence", sample, true)

	bscript := `
	// Just Testing For Loop
	echo @z |  tee /tmp/@z #for:@sequence:@z #as:@forout{add}
	`

	s.Compile(bscript)

	s.Evaluate()

	s.Schedule()

	s.Execute()

	out, er := shared.SharedVars.Get("@forout")

	if er != nil {
		panic(er)
	}

	arr := strings.Fields(out)

	sort.Strings(arr)

	log.Printf("Got : %v\n", arr)

	orig := strings.Fields(sample)

	if !reflect.DeepEqual(orig, arr) {
		t.Errorf("Expected sequence not matched got %v", arr)
	} else {
		t.Logf("test of for directive working successful")
	}

}

// func Test_sizes(t *testing.T) {
// 	z := shell.CMDWrap{CMD: &shell.SimpleCMD{
// 		COutStream: *bytes.NewBufferString("Yup tons of data"),
// 	}}

// 	z2 := &shell.CMDWrap{
// 		CMD: &shell.SimpleCMD{
// 			COutStream: *bytes.NewBufferString("Yup tons of data"),
// 		},
// 	}

// 	t.Logf("Size of struct is %v\n", unsafe.Sizeof(z))

// 	t.Logf("Size of struct is %v\n", unsafe.Sizeof(z2))
// }
