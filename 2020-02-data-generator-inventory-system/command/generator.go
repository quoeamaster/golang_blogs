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
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/spf13/cobra"
	"strings"
)

var genCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate supermarket transactions or inventory records; results would be written to elasticsearch directly. Use 'file' command to write results to files instead.",
	Long: `
generate supermarket transactions or inventory records; results would be written to files or elasticsearch directly. 
Use 'file' command to write results to files instead.
`,
	Run: func(cmd *cobra.Command, args []string) {
		c := new(GenerateCmdStruct)
		c.execute(cmd, args)
	},
}

const (
	genProfileInventory = "inventory"
	genProfileSales     = "sales"
	genProfileAll       = "all"
)

func init() {
	genCmd.AddCommand(genToFileCmd)
	rootCmd.AddCommand(genCmd)


	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// heyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// heyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	genCmd.PersistentFlags().StringP("source", "s", "datasource", "the folder containing the kml OR geojson files")
	genCmd.PersistentFlags().StringP("filename", "f", "location", "the name of the kml OR geojson files, e.g. filename=abc then abc.kml OR abc.geojson is expected")
	genCmd.PersistentFlags().StringP("profile", "p", genProfileInventory, "profile refers to which dataset to generate, valid option are 'inventory', 'sales' OR 'all'.")
	genCmd.PersistentFlags().Int16("size", 5, "number of records to create on SALES trx only; for inventory profile, this value is ignored")

	genCmd.Flags().StringP("elastichost", "", "http://localhost:9200", "elasticsearch host to connect to")
}

type GenerateCmdStruct struct {

}
func (c *GenerateCmdStruct) execute(cmd *cobra.Command, args []string)  {
	gUtil := NewGeneratorUtil()
	s, err := cmd.PersistentFlags().GetString("source")
	CommonPanic(err)

	f, err := cmd.PersistentFlags().GetString("filename")
	CommonPanic(err)

	p, err := cmd.PersistentFlags().GetString("profile")
	CommonPanic(err)

	size, err := cmd.PersistentFlags().GetInt16("size")
	CommonPanic(err)

	// generate the entries
	entryResponse := gUtil.GenTrx(s, f, p, size)
	switch p {
	case genProfileInventory:
		c.esInventoryIndex(entryResponse.InventoryList)
	}


	return
}

const (
	esIndexInventory = "m_supermarket_inventory"
)

func (c *GenerateCmdStruct) esInventoryIndex(data []InventoryTrxStruct) {
	es, err := elasticsearch.NewDefaultClient()
	CommonPanic(err)
	// assume index template already available in the elasticsearch cluster
	var bulkBuf strings.Builder

	for _, iVal := range data {
		bulkBuf.WriteString(`{"index":{}}
`)
		bulkBuf.WriteString(fmt.Sprintf(`{ "stock_in_price": %v, "stock_in_quantity": %v, "stock_in_date": "%v", "expiry_date": "%v", "product": { "id": "%v", "desc": "%v", "batch_id": "%v" }, "location": { "id": "%v", "name": "%v", "post_code": "%v", "coord": { "lat": %v, "lon": %v }}}
`,
			iVal.StockInCost, iVal.StockInQuantity,  GetESFormattedDate(iVal.StockInDate), GetESFormattedDate(iVal.ExpiryDate),
			iVal.Product.Id, iVal.Product.Desc, iVal.Product.BatchId,
			iVal.Location.Id, iVal.Location.Name, iVal.Location.PostCode, iVal.Location.Lat, iVal.Location.Lng))
	}
	bulkBuf.WriteString(`
`)
	req := esapi.BulkRequest{
		Index: esIndexInventory,
		Body: strings.NewReader(bulkBuf.String()),
	}
	response, err := req.Do(context.Background(), es)
	CommonPanic(err)
	defer response.Body.Close()
	if response.IsError() {
		panic(fmt.Sprintf("could not bulk ingest => %V", response.Status()))
	}
	rMap := ConvertESResponseToMap(*response)
	fmt.Println(fmt.Sprintf("DONE bulk ingest on %v documents, took %vms with errors? %v", len(data), rMap["took"], rMap["errors"]))
}







// **** write to file command ****

var genToFileCmd = &cobra.Command{
	Use:   "file",
	Short: "generate trx and write to file(s)",
	Long: `
generate trx and write to file(s)
`,
	Run: func(cmd *cobra.Command, args []string) {
		c := new(GenToFileCmdStruct)
		c.execute(cmd, args)
	},
}

type GenToFileCmdStruct struct {

}
func (c *GenToFileCmdStruct) execute(cmd *cobra.Command, args []string) {
	return
}

