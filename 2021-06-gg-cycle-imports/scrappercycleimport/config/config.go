package config

import (
	"encoding/json"
	"io/ioutil"
)

const (
	KeyConfigFolder  string = "folder"
	KeyConfigPattern string = "pattern"
)

// [cycle-import]
// - reason: accssing the package scrappercycleimport constants.
//	- solution A: move the constants to this package instead of the scrappercycleimport package.
//
// Config - the extracted configuration logic originated from the scrapper struct.
func Config(configFile string) (folder, pattern string, err error) {
	bContent, err := ioutil.ReadFile(configFile)
	if err != nil {
		return
	}
	// assume the file contents are json
	m := make(map[string]interface{})
	err = json.Unmarshal(bContent, &m)
	if err != nil {
		return
	}
	folder = m[KeyConfigFolder].(string)
	pattern = m[KeyConfigPattern].(string)
	return
}
