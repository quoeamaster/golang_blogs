package scrappercycleimport

import (
	"io/ioutil"
	"strings"

	config "github.com/quoeamaster/golang_blogs/ggcycleimport/scrappercycleimport/config"
)

// [cycle-import]
// - reason: these constants would be accessed by the scrappercycleimport/config package
// 	which becomes a circular dependency.
//	- solution A: move the constants to the scrappercycleimport/config package
// 	instead of the current package.
/*
const (
	KeyConfigFolder  string = "folder"
	KeyConfigPattern string = "pattern"
)*/

// ScrapperCycleImport - the same scrapper, but introduced a cycle import issue.
// To create the cycle import uncomment the methods and code blocks marked [cycle-import] and run unit tests.
type ScrapperCycleImport struct {
	// pattern - patterns of the file(s) to scrap for.
	pattern string

	// folder - the folder for scrapping.
	folder string
}

// [cycle-import]
func (s *ScrapperCycleImport) Config(configFile string) (err error) {
	s.folder, s.pattern, err = config.Config(configFile)
	return
}

func (s *ScrapperCycleImport) Scrap() (files []string, err error) {
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
