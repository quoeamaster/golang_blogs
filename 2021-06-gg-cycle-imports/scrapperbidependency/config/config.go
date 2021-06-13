package config

import (
	"encoding/json"
	"io/ioutil"
)

const (
	KeyConfigFolder  string = "folder"
	KeyConfigPattern string = "pattern"
)

// IValidator - a minimal interface required to "cast" our target
// scrapperbidependency.ScrapperBi instance.
type IValidator interface {
	// IsValueString - validates whether the given [v] is a string.
	IsValueString(v interface{}) (bool, string)
}

// [cycle-import]
// - reason: since this function need an instance of scrapperbidependency.ScrapperBi
// 	in a reverse manner (bi-directional dependency); a cycle-import would be introduced
//	- solution B: create an interface exposing just enough func points (IValidator) on this package
//		and "cast" our scrapperbidependency/ScrapperBi into this interface.
//
// Config - take in an instance of [scrapperbidependency.Scrapper] and
// would employ its [IsValueString] function to validate whether the config value is string.
//
// [before] func Config(s scrapper.ScrapperBi, configFile string) (folder, pattern string, err error) {
func Config(s IValidator, configFile string) (folder, pattern string, err error) {
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
	// validation on the value of the config (whether it is a string)
	var v interface{}

	v = m[KeyConfigFolder]
	if isString, sValue := s.IsValueString(v); isString {
		folder = sValue
	}
	v = m[KeyConfigPattern]
	if isString, sValue := s.IsValueString(v); isString {
		pattern = sValue
	}
	return
}
