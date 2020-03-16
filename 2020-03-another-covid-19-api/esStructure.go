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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	CaseConfirmed = "confirmed"
	CaseDeaths    = "deaths"
	CaseRecovered = "recovered"
)

type ESCaseDocumentModel struct {
	Country string				`json:"country"`
	CountryCode2 string			`json:"country_code"`
	Province string				`json:"province"`

	Location CountryLocation	`json:"location"`

	Case CaseTuple				`json:"case"`
	// es document_id
	docId string
}
type CaseTuple struct {
	DateInString string 	`json:"timestamp"` 		// date in history
	Value float64			`json:"count"` 			// actual number of cases
	CaseType string			`json:"case_type"`
}
type CountryLocation struct {
	Longitude float32			`json:"lon"`
	Latitude float32			`json:"lat"`
}

func (c *CaseTuple) ConvertDateFormat() {
	if strings.Compare(strings.Trim(c.DateInString, " "), "") != 0 {
		if parts := strings.Split(c.DateInString, "/"); len(parts) == 3 {
			c.DateInString = fmt.Sprintf("20%v-%02v-%02v", parts[2], parts[0], parts[1])
		} // end -- if parts are split correctly plus EXACTLY 3 parts...
	}
}

// generate a unique doc_id so that the existing data could be overwritten without a problem (assume format of the country code and province name stays the same)
// if not, could simply delete the existing index content and then insert all entries again
func (e *ESCaseDocumentModel) GenerateDocId() {
	proName := strings.Replace(e.Province, " ", "_", -1)
	proName = strings.Replace(proName, ",", "_", -1)
	proName = strings.Replace(proName, ".", "_", -1)

	e.docId = fmt.Sprintf("%v_%v_%v_%v",
		e.CountryCode2,
		proName,
		strings.Replace(e.Case.DateInString, "-", "_", -1),
		e.Case.CaseType)
}


func GenerateESDocsForBulk(list []ESCaseDocumentModel) (err error, bulkString string) {
	var sb strings.Builder

	for _, model := range list {
		sb.WriteString(fmt.Sprintf("{ \"index\": { \"_id\": \"%v\" } }\n", model.docId))
		bContent, err2 := json.Marshal(model)
		if err2 != nil {
			err = err2
			return
		}
		sb.WriteString(string(bContent))
		sb.WriteString("\n")
	}
	sb.WriteString("\n")
	bulkString = sb.String()
	return
}





const (
	contentTypeNDJson = "application/x-ndjson"
	contentTypeJson = "application/json"
)

type ESConnector struct {
	Host string
	Port int
	Index string
}
func NewESConnector(host string, port int, index string) (inst *ESConnector) {
	inst = new(ESConnector)
	inst.Host = host
	inst.Port = port
	inst.Index = index
	return
}

func (e *ESConnector) getTargetUrl() (url string) {
	url = fmt.Sprintf("http://%v:%v/%v", e.Host, e.Port, e.Index)
	return
}

// create the index + mapping if necessary
func (e *ESConnector) CreateIndexIfNotAvailable() (err error)  {
	targetUrl := e.getTargetUrl()

	// do a GET {index}/_mapping api first; would yield exception when index not found
	resp, err := http.Get(targetUrl+"/_mapping")
	/*
	if err == nil {
		// means already there and can leave
		fmt.Println("index exists~")
		return
	}
	*/
	// check response
	defer resp.Body.Close()
	bContents, err := ioutil.ReadAll(resp.Body)
	sContents := string(bContents)
	if strings.Index(sContents, "index_not_found_exception") == -1 {
		fmt.Println("index exists~")
		err = nil
		return
	}

	// create missing index
	var mapping = `{
  "mappings": {
    "properties": {
      "country": { "type": "text", "fields": { "raw": { "type": "keyword" } } },
      "country_code": { "type": "keyword" },
      "province": { "type": "text", "fields": { "raw": { "type": "keyword" } } },
      "location": { "type": "geo_point" },
      "case": {
        "properties": {
          "timestamp": { "type": "date" },
          "count": { "type": "integer" }
        }
      }
    }
  }
}`
	// 1. create client;
	// 2. set header (a MUST for es 6.x or above)
	// 3. check not just error object but also the contents within the response
	// 	(usually no error unless connection problem); hence need to check the content of the response
	hClient := &http.Client{}
	request, err := http.NewRequest(http.MethodPut, targetUrl, bytes.NewBuffer([]byte(mapping)))
	if err != nil {
		return
	}
	request.Header.Set("Content-Type", contentTypeJson)

	resp2, err := hClient.Do(request)
	if err != nil {
		return
	}
	defer resp2.Body.Close()
	bContents, err = ioutil.ReadAll(resp2.Body)
	if err != nil {
		return
	}
	fmt.Println("created index with following response:", string(bContents))

	return
}

// clean up the existing contents (delete_by_query)
func (e *ESConnector) CleanupIndex() (err error) {
	query := `{
  "query": { "match_all": {} }
}`

	resp, err := http.Post(e.getTargetUrl()+"/_delete_by_query", contentTypeJson, bytes.NewBuffer([]byte(query)))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	bContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	fmt.Println("clean up index with following response:", string(bContent))

	return
}

func (e *ESConnector) BulkIngest(query string) (err error) {
	resp, err := http.Post(e.getTargetUrl()+"/_bulk", contentTypeNDJson, bytes.NewBuffer([]byte(query)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bContent, err := ioutil.ReadAll(resp.Body)
	// error check
	if strings.Index(string(bContent), "\"errors\":true") != -1 {
		fmt.Println(string(bContent), "bulk finished WITH error")
	} else {
		fmt.Println("bulk ingest OK")
	}
	return
}


// TODO: simply delete the existing index content and then insert all entries again