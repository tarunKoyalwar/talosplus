package subcmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/tarunKoyalwar/talosplus/pkg/db"
	"golang.design/x/clipboard"
)

// PrepareDB : Connects DB and required collection
func PrepareDB() error {
	return db.UseMongoDB(DefaultSettings.ActiveURL, DefaultSettings.ActiveDB, DefaultSettings.ActiveColl)
}

// HasStdin : Check if Stdin is present
func HasStdin() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	mode := stat.Mode()

	isPipedFromChrDev := (mode & os.ModeCharDevice) == 0
	isPipedFromFIFO := (mode & os.ModeNamedPipe) != 0

	return isPipedFromChrDev || isPipedFromFIFO
}

// Readclipboard : Read Data From Clipboard
func Readclipboard() string {
	x := clipboard.Read(clipboard.FmtText)

	return string(bytes.TrimSpace(x))
}

// Writeclipboard : Write Data to Clipboard
func Writeclipboard(data string) {
	go func() {
		<-clipboard.Write(clipboard.FmtText, []byte(data))
	}()

	time.Sleep(time.Duration(2) * time.Second)
}

func manageinput(data *string, args []string) {
	if HasStdin() {
		*data = GetStdin()
	} else if len(args) > 0 {
		*data = strings.TrimSpace(args[0])
	} else if ClipBoardIn {
		*data = Readclipboard()
	} else {
		panic("No Data Source!!!")

	}
}

// GetStdin : Get all Data present on stdin
func GetStdin() string {
	bin, _ := ioutil.ReadAll(os.Stdin)
	return string(bin)
}
