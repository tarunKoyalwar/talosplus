package shell

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/tarunKoyalwar/talosplus/pkg/alerts"
	"github.com/tarunKoyalwar/talosplus/pkg/ioutils"
	"github.com/tarunKoyalwar/talosplus/pkg/stringutils"
)

// Notifications : Sends and parses notifications for each cmd
type Notifications struct {
	NotifyEnabled bool   //to check if notification should be sent
	NotifyLen     bool   //Only Notify Len and Not Entire Text
	NotifyMsg     string //Msg to be sent along with Notify
	Raw           string
	Comment       string // Comment
}

// Parse : Parse Notification Settings
func (n *Notifications) Parse(raw string) {
	/*

		Syntax1:
		#notify{Some Notification Text}

		Response:
		Some Notification Text : %stdout

		Default:
		Results for %comment : %stdout


		Syntax2:

		#notifylen{Some Notification Text}

		Response:
		Some Notification Text len(%stdout)

		Default:
		Found Total for %comment : len(%stdout)


		#Parse Steps
		1. Identify if #notifiy exists
		2. Extract notify msg
		3. Set Default Messages

	*/

	temp := raw

	if strings.Contains(temp, "#notify") {
		n.NotifyEnabled = true

		// This will extract message
		regex := `(#notify|#notifylen){(?P<msg>.*?)}`

		re := regexp.MustCompile(regex)

		matches := re.FindStringSubmatch(temp)

		if len(matches) == 3 {
			group1 := matches[1]
			n.NotifyMsg = matches[2]

			if group1 == "#notifylen" {
				n.NotifyLen = true
			}

			//remove #notify from string
			n.Raw = strings.ReplaceAll(temp, matches[0], "")
			return

		} else {
			// cmd only has #notify no msg was specified
			n.NotifyMsg = ""

			if strings.Contains(temp, "#notifylen") {
				n.NotifyLen = true
				n.Raw = strings.ReplaceAll(temp, "#notifylen", "")
			} else {
				n.Raw = strings.ReplaceAll(temp, "#notify", "")
			}

			if n.NotifyLen {
				n.NotifyMsg = "Found Total for\n" + n.Comment + " : "
			} else {
				n.NotifyMsg = "Results for \n" + n.Comment + "\n"
			}

		}
	} else {
		n.Raw = raw
	}

}

// Notify : Notifies completion of command
func (n *Notifications) Notify(dat string) {

	// fmt.Printf("\n\nCalled Notify with %v\n\n", dat)

	if dat == "" {
		ioutils.Cout.PrintWarning("Empty String")
		return
	}

	var senderr error

	if !n.NotifyEnabled {
		return
	}

	if alerts.Alert == nil {
		return
	}

	if n.NotifyLen {
		final := fmt.Sprintf("%v %v\n", n.NotifyMsg, len(stringutils.Split(dat, '\n')))

		senderr = alerts.Alert.SendEmbed(final, map[string]string{
			"CMD": n.Raw,
		})

	} else {
		final := fmt.Sprintf("%v\n%v\n", n.NotifyMsg, dat)

		senderr = alerts.Alert.SendEmbed(final, map[string]string{
			"CMD": n.Raw,
		})
	}

	if senderr != nil {
		ioutils.Cout.PrintWarning("Failed to send notification\n%v", senderr)
	}

}

// NotifyErr : Notifies Err
func (n *Notifications) NotifyErr(er error) {}

func NewNotification(raw string, comment string) (*Notifications, string) {
	z := Notifications{
		Comment: comment,
	}

	z.Parse(raw)

	return &z, z.Raw
}
