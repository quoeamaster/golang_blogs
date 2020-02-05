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
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// heyCmd represents the hey command
var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "parse and prepare supermarket data from the corresponding kml OR geojson files",
	Long: `
parse and prepare supermarket data from the corresponding kml OR geojson files
`,
	Run: func(cmd *cobra.Command, args []string) {
		p := new(ParserCmdStruct)
		p.execute(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(parseCmd)

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// heyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// heyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	parseCmd.Flags().StringP("source", "s", "datasource", "the folder containing the kml OR geojson files")
	parseCmd.Flags().StringP("filename", "f", "location", "the name of the kml OR geojson files, e.g. filename=abc then abc.kml OR abc.geojson is expected")
}

const (
	kmlTagPlacemarkStart = "<Placemark"
	kmlTagPlacemarkEnd   = "</Placemark"

	kmlTagNameRoot    = "<name>"
	kmlTagNameRootEnd = "</name>"

	kmlTagSimpleData    = "<SimpleData name=\""
	kmlTagSimpleDataEnd = "</SimpleData>"

	kmlTagCommonCloseChar = ">"

	kmlTagCoord    = "<coordinates>"
	kmlTagCoordEnd = "</coordinates>"

	kmlKeyLicName    = "LIC_NAME"
	kmlKeyBlockHouse = "BLK_HOUSE"
	kmlKeyStreetName = "STR_NAME"
	kmlKeyPostalCode = "POSTCODE"
	kmlKeyLicenseNum = "LIC_NO"
	kmlKeyIncCrc     = "INC_CRC"
	kmlKeyFmelUpdD   = "FMEL_UPD_D"
)

type PlacemarkStruct struct {
	ID string  `json:"id"`    // id
	Name string `json:"name"` // LIC_NAME
	BlockHouse string `json:"block_house"` // BLK_HOUSE
	StreetName string `json:"street_name"` // STR_NAME
	Postcode string  `json:"postcode"` // POSTCODE
	LicenseNum string  `json:"license_num"` // LIC_NO
	IncCrc string  `json:"inc_crc"` // INC_CRC
	FmelUpdD string `json:"fmel_upd_d"` // FMEL_UPD_D
	Lng float32 `json:"lng"` // lng,lat...
	Lat float32 `json:"lat"` // lng,lat...
}

type ParserCmdStruct struct {

}

func (p *ParserCmdStruct) execute(cmd *cobra.Command, args []string) {
	source := cmd.Flag("source").Value.String()
	filename := cmd.Flag("filename").Value.String()

	// try to look for kml first
	p.parseKML(source, filename)

	// TODO: ignore the parsing of geojson for this release
}

func (p *ParserCmdStruct) parseKML(source, filename string) {
	kml := fmt.Sprintf("%v%v%v.kml", source, string(os.PathSeparator), filename)
	fInfo, err := os.Stat(kml)
	if os.IsNotExist(err) {
		panic(err)
	}
	if fInfo.IsDir() {
		panic(fmt.Sprintf("the given folder and filename path is a DIRECTORY instead! %v", kml))
	}

	fPtr, err := os.OpenFile(kml, os.O_RDONLY, 0755)
	if err != nil {
		panic(err)
	}
	defer fPtr.Close()

	var dataLines []string
	var locationList []PlacemarkStruct
	isDataCollecting := false

	scannerPtr := bufio.NewScanner(fPtr)
	for scannerPtr.Scan() {
		line := scannerPtr.Text()
		if strings.Index(line, kmlTagPlacemarkStart) != -1 {
			isDataCollecting = true
			dataLines = append(dataLines, line)
		} else if strings.Index(line, kmlTagPlacemarkEnd) != -1 {
			isDataCollecting = false
			dataLines = append(dataLines, line)
			// need to parse them into data-structures
			locationList = append(locationList, p.parsePlacemarkSlice(dataLines))
			dataLines = []string{} // reset

		} else if isDataCollecting == true {
			dataLines = append(dataLines, line)
		}
	}
	// write to file (prepared)
	bContent, err := json.Marshal(locationList)
	if err != nil  {
		panic(err)
	}
	targetFilename := fmt.Sprintf("%v%v%v_prepared.json", source, string(os.PathSeparator), filename)
	err = ioutil.WriteFile(targetFilename, bContent, 0644)
	if err != nil  {
		panic(err)
	}
}

func (p *ParserCmdStruct) parsePlacemarkSlice(lines []string) (location PlacemarkStruct) {
	location = PlacemarkStruct{}

	for _, line := range lines {
		if strings.Index(line, kmlTagNameRoot) == 0 {
			location.ID = p.getTagValue(line, kmlTagNameRoot, kmlTagNameRootEnd)
			continue
		}
		if strings.Index(line, kmlTagSimpleData) == 0 {
			p.parseSimpleDataFromLine(line, &location)
			continue
		}
		if strings.Index(line, kmlTagCoord) == 0 {
			p.parseCoord(line, &location)
			continue
		}
	}
	return
}

func (p *ParserCmdStruct) getTagValue(line, tag, endTag string) (val string) {
	startIdx := strings.Index(line, tag)
	endIdx := strings.Index(line, endTag)

	if startIdx != -1 {
		startIdx = startIdx + len(tag)
	}
	// not found...?!
	if endIdx == -1 || startIdx == -1 {
		return
	}
	// extract
	val = line[startIdx:endIdx]
	return
}
func (p *ParserCmdStruct) parseSimpleDataFromLine(line string, placemarkPtr *PlacemarkStruct) {
	sIdx := strings.Index(line, kmlTagSimpleData) + len(kmlTagSimpleData)

	fieldName := line[sIdx:]
	fieldName = fieldName[0:strings.Index(fieldName,"\"")]

	// extract value...
	sIdx = strings.Index(line, kmlTagCommonCloseChar) + 1
	eIdx := strings.Index(line, kmlTagSimpleDataEnd)

	// not found??!
	if sIdx == -1 || eIdx == -1 {
		return
	}
	val := line[sIdx:eIdx]
	// set value back
	switch fieldName {
	case kmlKeyLicName:
		placemarkPtr.Name = val
	case kmlKeyBlockHouse:
		placemarkPtr.BlockHouse = val
	case kmlKeyStreetName:
		placemarkPtr.StreetName = val
	case kmlKeyPostalCode:
		placemarkPtr.Postcode = val
	case kmlKeyLicenseNum:
		placemarkPtr.LicenseNum = val
	case kmlKeyIncCrc:
		placemarkPtr.IncCrc = val
	case kmlKeyFmelUpdD:
		placemarkPtr.FmelUpdD = val
	}
}
func (p *ParserCmdStruct) parseCoord(line string, placemarkPtr *PlacemarkStruct) {
	coordParts := strings.Split(p.getTagValue(line, kmlTagCoord, kmlTagCoordEnd), ",")
	if coordParts != nil && len(coordParts) == 3 {
		placemarkPtr.Lat, _ = p.convertStringToFloat32(coordParts[1])
		placemarkPtr.Lng, _ = p.convertStringToFloat32(coordParts[0])
	}
}
func (p *ParserCmdStruct) convertStringToFloat32(val string) (fVal float32, err error) {
	f64, err := strconv.ParseFloat(val, 32)
	if err != nil {
		fVal = 0
	}
	fVal = float32(f64)
	return
}
