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

import "github.com/spf13/cobra"

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

	genCmd.Flags().StringP("elastichost", "", "http://localhost:9200", "elasticsearch host to connect to")
}

type GenerateCmdStruct struct {

}
func (c *GenerateCmdStruct) execute(cmd *cobra.Command, args []string)  {
	gUtil := NewGeneratorUtil()
	gUtil.GenInventoryTrx()

	return
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

