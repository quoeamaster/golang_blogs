package scrapperbidependency

import (
	"io/ioutil"
	"strings"

	config "github.com/quoeamaster/golang_blogs/ggcycleimport/scrapperbidependency/config"
)

type ScrapperBi struct {
	pattern string
	folder  string
}

func (s *ScrapperBi) Config(configFile string) (err error) {
	// [before]
	//s.folder, s.pattern, err = config.Config(*s, configFile)
	s.folder, s.pattern, err = config.Config(s, configFile)
	return
}

// [cycle-import]
// - reason: since this function would be called from package scrapperbidependency/config
// 	in a reverse manner (bi-directional dependency)
//	- solution B: create an interface exposing just enough func points on the
//		package scrapperbidependency/config and "cast" our ScrapperBi into that interface.
//
// IsValueString - checks the provided [v] is a string or not.
func (s *ScrapperBi) IsValueString(v interface{}) (isString bool, sValue string) {
	sValue, isString = v.(string)
	return
}

func (s *ScrapperBi) Scrap() (files []string, err error) {
	files = make([]string, 0)
	fArray, err := ioutil.ReadDir(s.folder)
	if err != nil {
		return
	}
	// search / filter
	for _, f := range fArray {
		if !f.IsDir() {
			if strings.Contains(f.Name(), s.pattern) {
				files = append(files, f.Name())
			}
		}
	}
	return
}
