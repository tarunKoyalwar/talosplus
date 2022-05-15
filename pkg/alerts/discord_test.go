package alerts_test

import (
	"os"
	"testing"

	"github.com/tarunKoyalwar/talosplus/pkg/alerts"
)

func Test_Alerts(t *testing.T) {
	id := os.Getenv("DISCORD_WID")
	tok := os.Getenv("DISCORD_WTOKEN")

	alerts.Alert = alerts.NewDiscordHook(id, tok)

	alerts.Alert.Title = "Testing Hook"

	err := alerts.Alert.SendEmbed("Just a new simple test", map[string]string{})

	if err != nil {
		t.Errorf(err.Error())
	}
}
