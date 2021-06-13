package scrapper

import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

const (
	KeyConfigFolder  string = "folder"
	KeyConfigPattern string = "pattern"
)

// Scrapper - a folder scrapper, scraps for file(s) that matches search criteria.
type Scrapper struct {
	// pattern - patterns of the file(s) to scrap for.
	pattern string

	// folder - the folder for scrapping.
	folder string
}

// Config - config the scrapper.
func (s *Scrapper) Config(configFile string) (err error) {
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
	s.folder = m[KeyConfigFolder].(string)
	s.pattern = m[KeyConfigPattern].(string)
	return
}

// Scrap - logics to scrap out the files that match the pattern.
// (implements a non reg-exp substring match)
func (s *Scrapper) Scrap() (files []string, err error) {
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
