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
package command

import (
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"time"
)

func CommonPanic(err error) {
	if err != nil {
		panic(err)
	}
}

// common method to parse a esapi response into map[string]interface{}
func ConvertESResponseToMap(r esapi.Response) map[string]interface{} {
	var m map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		CommonPanic(err)
	}
	return m
}

func GetESFormattedDate(d time.Time) string {
	return d.Format("2006-01-02T15:04:05")
}