package test

import (
	"strings"
	"testing"

	scrapper "github.com/quoeamaster/golang_blogs/ggcycleimport/scrappercycleimport"
)

const (
	configFileForUnitTest string = "./cfgTest.json"
	matchedFilename       string = "scrapper_test.go"
)

func TestScrapper(t *testing.T) {
	instance := new(scrapper.ScrapperCycleImport)
	if err := instance.Config(configFileForUnitTest); err != nil {
		t.Fatalf("error in configuring the Scrapper, [%v]", err)
	}
	files, err := instance.Scrap()
	if err != nil {
		t.Fatalf("error in scrapping, [%v]", err)
	}
	// verification
	// should have "scrapper_test.go" as the only match.
	if len(files) == 0 {
		t.Fatalf("expect exactly 1 match, actual is 0")
	}
	numMatches := 0
	for _, file := range files {
		if strings.Compare(file, matchedFilename) == 0 {
			numMatches++
		}
	}
	if numMatches != 1 {
		t.Fatalf("should have an EXACT 1 match for [%v], actual [%v]", matchedFilename, numMatches)
	}
}
