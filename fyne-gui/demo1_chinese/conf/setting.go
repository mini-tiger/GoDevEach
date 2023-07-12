package conf

import (
	"encoding/json"
	"fyne.io/fyne/v2/app"
	"os"
	"path/filepath"
)

type Setting struct {
	Index int
}

var SettingPath string

var AppSetting Setting

func LoadSetting() error {
	s := &app.SettingsSchema{}
	path := s.StoragePath()
	SettingPath = path
	file, err := os.Open(path) // #nosec
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(filepath.Dir(path), 0700)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	decode := json.NewDecoder(file)

	return decode.Decode(&AppSetting)
}

func (s *Setting) SaveToFile() error {
	path := SettingPath
	err := os.MkdirAll(filepath.Dir(path), 0700)
	if err != nil { // this is not an exists error according to docs
		return err
	}

	data, err := json.Marshal(&s)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}