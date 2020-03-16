/*
Copyright © 2020 quo master

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	apiUrl           = "https://coronavirus-tracker-api.herokuapp.com/"
	apiPathAll       = "all"
	apiPathConfirmed = "confirmed"
	apiPathDeaths    = "deaths"
	apiPathRecovered = "recovered"

	datasetFolder = "dataset"
)

func main()  {
	var dataMap map[string]interface{}

	// get today's date
	todayString, _ := getTodayDateInUTC()
	// check if the data already there or not (reduce network traffic if accessing data from the same url all the time within a day)
	if available := isDataFileAvailable(todayString); available {
		// read and load it
		err2, dataMap2 := loadFromDatafile(todayString)
		genericErrHandler(err2, "loading data from data file")
		dataMap = dataMap2

	} else {
		err2, msg, dataMap2, contentsInByte := getJsonData(apiUrl + apiPathAll)
		genericErrHandler(err2, msg)

		// save it to disk
		err2 = writeDatafile(todayString, contentsInByte)
		genericErrHandler(err2, "write data file")

		dataMap = dataMap2
	}

	// a. parse confirmed
	var confirmedTotal float64
	var deathTotal float64
	var recoveredTotal float64

	err, confirmedTotal, pConfirmedList := parseCases(dataMap[CaseConfirmed].(map[string]interface{}), CaseConfirmed)
	genericErrHandler(err, "parsing confirmed cases")
	fmt.Println("["+todayString+"]\nconfirmed:", confirmedTotal)

	err, deathTotal, pDeathList := parseCases(dataMap[CaseDeaths].(map[string]interface{}), CaseDeaths)
	genericErrHandler(err, "parsing deaths cases")
	fmt.Println(len(pDeathList), deathTotal)

	err, recoveredTotal, pRecoveredList := parseCases(dataMap[CaseRecovered].(map[string]interface{}), CaseRecovered)
	genericErrHandler(err, "parsing recovered cases")
	fmt.Println(len(pRecoveredList), recoveredTotal)

	// es-connector
	es := NewESConnector("localhost", 9200, "covid_19_api")
	// create index definition if necessary
	err = es.CreateIndexIfNotAvailable()
	if err != nil {
		fmt.Println("something wrong when creating index:", err)
	}
	// clean up existing data
	err = es.CleanupIndex()
	if err != nil {
		fmt.Println("something wrong when cleanup index:", err)
	}
	// bulk ingest for confirmed
	// write to es in bulk ingest format
	err, sBulk := GenerateESDocsForBulk(pConfirmedList)
	genericErrHandler(err, "generating bulk syntax for confirmed list")
	//fmt.Println(sBulk, len(sBulk), confirmedTotal, deathTotal, recoveredTotal)
	//fmt.Println(confirmedTotal, deathTotal, recoveredTotal)
	err = es.BulkIngest(sBulk)
	if err != nil {
		fmt.Println("something wrong with bulk ingest confirmed cases", err)
	}
	err, sBulk = GenerateESDocsForBulk(pDeathList)
	genericErrHandler(err, "generating bulk syntax for death list")
	err = es.BulkIngest(sBulk)
	if err != nil {
		fmt.Println("something wrong with bulk ingest death cases", err)
	}
	err, sBulk = GenerateESDocsForBulk(pRecoveredList)
	genericErrHandler(err, "generating bulk syntax for recovered list")
	err = es.BulkIngest(sBulk)
	if err != nil {
		fmt.Println("something wrong with bulk ingest recovered cases", err)
	}
}

// generic error handler
func genericErrHandler(err error, message string) {
	if err != nil {
		fmt.Println("*** error:", message)
		panic(err)
	}
}

// simple method to get today's date plus its string format (UTC)
func getTodayDateInUTC() (dateInString string, date time.Time) {
	// need to truncate??? (no for now)
	date = time.Now().UTC()
	dateInString = date.Format("2006-01-02")

	return
}

// check if the data file exists (no need to get from url again and again; good for development)
func isDataFileAvailable(today string) (available bool) {
	available = false
	dataFile := fmt.Sprintf("./%v%v%v.json", datasetFolder, string(os.PathSeparator), today)
	_, err := os.Stat(dataFile)
	if err != nil && os.IsNotExist(err) {
		available = false
	} else {
		available = true
	}
	return
}

// write all json content to a file
func writeDatafile(today string, data []byte) (err error)  {
	err = ioutil.WriteFile(fmt.Sprintf("./%v%v%v.json", datasetFolder, string(os.PathSeparator), today), data, 0x777)
	return
}

// load contents from data file and convert to map object
func loadFromDatafile(today string) (err error, dataMap map[string]interface{})  {
	bContent, err := ioutil.ReadFile(fmt.Sprintf("./%v%v%v.json", datasetFolder, string(os.PathSeparator), today))
	if err != nil {
		return
	}
	dataMap = make(map[string]interface{})
	err = json.Unmarshal(bContent, &dataMap)

	return
}

// retrieve json data from the api
func getJsonData(url string) (err error, msg string, dataMap map[string]interface{}, contentInBytes []byte) {
	response, err := http.Get(url)
	if err != nil {
		msg = "get API on url:" + url
		return
	}

	defer response.Body.Close()
	contentInBytes, err = ioutil.ReadAll(response.Body)
	if err != nil {
		msg = "reading content of api url"
		return
	}

	dataMap = make(map[string]interface{})
	err = json.Unmarshal(contentInBytes, &dataMap)
	if err != nil {
		msg = "unmarshalling json content to a map"
		return
	}
	//contentsInString = string(bContent)
	//contentsInBytes = bContent

	return
}



func parseCases(dataMap map[string]interface{}, caseType string) (err error, grandTotal float64, modelList []ESCaseDocumentModel) {
	if dataMap == nil || len(dataMap) == 0 {
		err = errors.New("data-map provided is empty")
		return
	}
	// grand total for all country-province pair(s)
	grandTotal = 0.0
	modelList = make([]ESCaseDocumentModel, 0)

	pCProvinceLst := dataMap["locations"].([]interface{})
	for _, itemInterface := range pCProvinceLst {
		item := itemInterface.(map[string]interface{})

		// parse LatLon
		latLonMap := item["coordinates"].(map[string]interface{})

		tmpStringVal := latLonMap["long"].(string)
		tmpF64Val, _ := strconv.ParseFloat(tmpStringVal, 64)
		tLongitude := float32(tmpF64Val)

		tmpStringVal = latLonMap["lat"].(string)
		tmpF64Val, _ = strconv.ParseFloat(tmpStringVal, 64)
		tLatitude := float32(tmpF64Val)

		// parse country, province
		tCountry := item["country"].(string)
		tProvince := item["province"].(string)
		if strings.Compare(strings.Trim(tProvince, " "), "") == 0 {
			tProvince = tCountry
		}
		tCountryCode2 := item["country_code"].(string)

		// grandtotal country level sum-up
		grandTotal += item["latest"].(float64)

		// histories
		historyInterface := item["history"].(interface{})
		historyMap := historyInterface.(map[string]interface{})
		for tDate, value := range historyMap {
			pESModel := new(ESCaseDocumentModel)
			pCase := new(CaseTuple)
			pLoc := new(CountryLocation)

			pCase.Value = value.(float64)
			pCase.DateInString = tDate // convert from 1/22/20 to 2020-01-22
			pCase.ConvertDateFormat()
			pCase.CaseType = caseType
			pESModel.Case = *pCase

			pESModel.CountryCode2 = tCountryCode2
			pESModel.Country = tCountry
			pESModel.Province = tProvince

			pLoc.Longitude = tLongitude
			pLoc.Latitude = tLatitude
			pESModel.Location = *pLoc

			// generate doc_id
			pESModel.GenerateDocId()
			//fmt.Println(pESModel.docId)

			modelList = append(modelList, *pESModel)
		} // end -- for (history of a LOCATION)
	} // end -- for (all LOCATION..S)
	return
}

