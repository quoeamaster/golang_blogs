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
	Short: "generate supermarket transactions or inventory records; results would be written to files or elasticsearch directly",
	Long: `
generate supermarket transactions or inventory records; results would be written to files or elasticsearch directly
`,
	Run: func(cmd *cobra.Command, args []string) {
		execute(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(genCmd)

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// heyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// heyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	genCmd.Flags().StringP("source", "s", "datasource", "the folder containing the kml OR geojson files")
	genCmd.Flags().StringP("filename", "f", "location", "the name of the kml OR geojson files, e.g. filename=abc then abc.kml OR abc.geojson is expected")
}

func execute(cmd *cobra.Command, args []string)  {
	return
}