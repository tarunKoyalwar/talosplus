package shelltest_test

import (
	"os"
	"strings"
	"testing"

	"github.com/tarunKoyalwar/talosplus/pkg/alerts"
	"github.com/tarunKoyalwar/talosplus/pkg/db"
	"github.com/tarunKoyalwar/talosplus/pkg/shell"
	"github.com/tarunKoyalwar/talosplus/pkg/stringutils"
)

func Test_Echo(t *testing.T) {

	g := shell.NewCMDWrap("echo 'hello world'", "")
	g.Process()

	if err := g.Execute(); err != nil {
		HandleErrors(err, t, "Failed to Run Command %v\n returned %v", g.Raw, err)
	}

	out := strings.TrimSpace(g.CMD.COutStream.String())
	if out != "hello world" {
		t.Errorf("Failed to run echo\n received %v Instead of hello world", out)
	} else {
		t.Logf("Passed Echo Command Test")
	}
}

func Test_File_Export(t *testing.T) {
	g := shell.NewCMDWrap(" cat /etc/passwd | grep ':0:' | cut -d ':' -f 1| tee @outfile #as:@etctest ", "")
	g.Process()

	if err := g.Execute(); err != nil {
		HandleErrors(err, t, "Failed to run cmd %v", g.Raw)
	}

	reqout, er1 := shell.Buffers.Get("@etctest")
	if er1 != nil {
		HandleErrors(er1, t, "Failed to get Exported Value of a cmd")
	}

	if reqout != "root" {
		t.Errorf("Failed to test file export\n Expected root but received %v", reqout)
	} else {
		t.Logf("Passed File Export test Using `@outfile`")
	}

}

func Test_Supply_Stdin(t *testing.T) {
	shell.Buffers.Set("@stdincheck", "Passed to Stdin", true)
	g := shell.NewCMDWrap("tee /dev/null #from:@stdincheck", "")
	g.Process()

	if err := g.Execute(); err != nil {
		HandleErrors(err, t, "Failed to Run Command %v\n returned %v", g.Raw, err)
	}

	out := strings.TrimSpace(g.CMD.COutStream.String())
	if out != "Passed to Stdin" {
		t.Errorf("Expected `Passed to Stdin` but received `%v`", out)
	} else {
		t.Logf("Passed Pipe Stdin")
	}
}

func Test_CMD_Export(t *testing.T) {
	g := shell.NewCMDWrap("echo 'hello world' #as:@echoexport", "")
	g.Process()

	if err := g.Execute(); err != nil {
		HandleErrors(err, t, "Failed to Run Command %v\n returned %v", g.Raw, err)
	}

	out, er1 := shell.Buffers.Get("@echoexport")
	if er1 != nil {
		HandleErrors(er1, t, "Failed to get Exported Value of a cmd")
	}
	if out != "hello world" {
		t.Errorf("Failed Basic_Export_Command Test received %v Instead of hello world", out)
	} else {
		t.Log("Passed Command Export Test")
	}
}

func Test_Alerts(t *testing.T) {
	id := os.Getenv("DISCORD_WID")
	tok := os.Getenv("DISCORD_WTOKEN")

	alerts.Alert = alerts.NewDiscordHook(id, tok)

	alerts.Alert.Title = "Testing"

	g := shell.NewCMDWrap("echo 'Hello Luci!!' #notify{Test Output}", "")
	g.Process()

	if err := g.Execute(); err != nil {
		HandleErrors(err, t, "Failed to Run Command %v\n returned %v", g.Raw, err)
	}

	// g.Export()

}

func HandleErrors(er error, t *testing.T, msg string, a ...any) {
	if er != nil {
		t.Logf(msg, a...)
		t.Errorf(msg, a...)
		t.Fatal(er)
	}
}

func TestMain(m *testing.M) {
	val := os.Getenv("RUN_MONGODB_TEST")

	if val == "" {
		db.UseBBoltDB(os.TempDir(), stringutils.RandomString(6)+".db", "Complex_test")
	} else {
		db.UseMongoDB("", stringutils.RandomString(6)+".db", "Complex_test")
	}

	m.Run()
}
