package subcmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

var SettingsFile string

var DefaultSettings *Settings

// Settings : Settings that are stored in Config File
type Settings struct {
	ActiveURL    string `json:"activeurl"`
	ActiveDB     string `json:"activedb"`
	ActiveColl   string `json:"activecoll"`
	ActiveScript string `json:"activescript"`
}

// Available : Check If Settings are available
func (s *Settings) Available() bool {
	if s.ActiveDB == "" || s.ActiveColl == "" {
		return false
	} else {
		return true
	}
}

// Save
func (s *Settings) Save() error {

	bin, _ := json.Marshal(s)
	err := ioutil.WriteFile(SettingsFile, bin, 0644)

	return err
}

// LoadSettings
func LoadSettings() {
	_, err := os.Stat(SettingsFile)
	if err != nil {
		DefaultSettings = &Settings{}
		return
	}

	bin, _ := ioutil.ReadFile(SettingsFile)

	var z Settings

	json.Unmarshal(bin, &z)

	DefaultSettings = &z

}

// CreateSettings : Create Settings if not exist
func CreateSettings(wdir string) (string, error) {

	//check if cache dir exists and dir name
	_, err := os.Stat(wdir)

	if err != nil {

		//Create New DIrectory
		err := os.Mkdir(wdir, 0755)
		if err != nil {
			return "", err
		}
	}

	return wdir, nil
}

func init() {
	config, _ := os.UserConfigDir()
	dirpath := path.Join(config, "talos")

	CreateSettings(dirpath)

	SettingsFile = path.Join(dirpath, "talos.json")
}
