package shared

import (
	"os"
	"path"
	"strconv"

	"github.com/tarunKoyalwar/talosplus/pkg/alerts"
	"github.com/tarunKoyalwar/talosplus/pkg/ioutils"
)

// DefaultSettings : Self Explainatory
var DefaultSettings *Settings = NewSettings()

// Settings : All Settings and Env Variables
type Settings struct {
	Purge             bool   // Skip Cached Data
	Limit             int    // Max Concurrent program
	ProjectExportName string // Name for directory containing all project exports
	CacheDIR          string // Cache Directory (default : temp)
	ProjectName       string // DIR Name where cache is saved
}

// ParseSettings : Parse Settings From Variables
func (s *Settings) ParseSettings(m map[string]string) {
	for k, v := range m {

		if v == "" {
			continue
		}

		// Check for global settings
		switch k {
		case "@pname":
			s.ProjectName = v
			s.ProjectExportName = v + "exports"

		case "@cachedir":
			s.CacheDIR = v

		case "@purge":
			z, er := strconv.ParseBool(v)
			if er == nil {
				s.Purge = z
			}

		case "@notifytitle":
			if alerts.Alert != nil {
				alerts.Alert.Title = v
			}

		case "@disablenotify":
			z, er := strconv.ParseBool(v)
			if er == nil {
				if alerts.Alert != nil {
					alerts.Alert.Disabled = z
				}
			}

		case "@limit":
			z, er := strconv.Atoi(v)
			if er == nil {
				if z < 2 {
					z = 2
				}
				s.Limit = z
			}
		}

	}

	if s.ProjectName == "" {
		//warn and set as phontom
		ioutils.Cout.PrintWarning("Project Name Not Set Fallback: talos")
		s.ProjectName = "talos"
		s.ProjectExportName = "talosExports"

	}

	if s.CacheDIR == "" {
		ioutils.Cout.PrintWarning("CacheDIR Not Set Fallback: system temp directory")
		s.CacheDIR = os.TempDir()
	}

}

func NewSettings() *Settings {
	s := Settings{
		Limit:             4,
		ProjectExportName: "talosExports",
		ProjectName:       "talos",
		CacheDIR:          os.TempDir(),
	}

	return &s
}

// CreateDirectoryIfNotExist : Creates Directory IF Not Present
func (s *Settings) CreateDirectoryIfNotExist(directory string) (string, error) {
	wdir := path.Join(s.CacheDIR, directory)

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
