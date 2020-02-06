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
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type generatorUtil struct {
	invList []string				// PS. prepared during GenTrx
	locationList []PlacemarkStruct	// PS. loaded on demand -> u.loadPlacemartList(source, filename)

	clientDemoList []clientDemo		// prepared during init -> prepareRandomData()
	occupationList []string			// prepared during init -> prepareRandomData()

}
func NewGeneratorUtil() (instance *generatorUtil) {
	instance = new(generatorUtil)
	instance.prepareRandomData()
	return
}

const (
	srcInventoryBaseFilename = "inventory_base.txt"
	srcOnlineSalesCsv = "sourceOnlineSales.csv"
)

func (u *generatorUtil) GenTrx(source, filename, profile string, size int32) (resp *EntryResponse) {
	resp = new(EntryResponse)
	resp.Profile = profile

	if strings.Compare(genProfileInventory, profile) == 0 {
		invList := u.generateInventoryTrx(source, filename)
		resp.InventoryList = invList
		// TODO: write to file... for TESTING only
		//bContent, err := json.Marshal(invList)
		//CommonPanic(err)

		//err = ioutil.WriteFile(fmt.Sprintf("%v%vtest.json", source, string(os.PathSeparator)), bContent, 0755)
		//CommonPanic(err)
	} else if strings.Compare(genProfileSales, profile) == 0 {
		resp.SalesList = u.generateSalesTrx(source, filename, size)
	}
	// TODO: other profiles (all)
	return
}
func (u *generatorUtil) generateInventoryTrx(source, filename string) (inventoryList []InventoryTrxStruct) {
	// read all entries from Online Sales.csv... create the inventory list
	// save the list to "inventory_base.txt" => 1 line with "," separated
	if len(u.invList) == 0 {
		u.invList = u.getInventoryList(source, filename)
	}
	// load locations
	if len(u.locationList) == 0 {
		u.locationList = u.loadPlacemartList(source, filename)
	}

	for _, inv := range u.invList {
		eInv := new(InventoryTrxStruct)
		eInv.StockInCost = u.getRandomFloat32(20, 160)
		eInv.StockInQuantity = int32(u.getRandomInteger(500, 10000))
		eInv.StockInDate = u.getRandomDate(180, 365)
		eInv.ExpiryDate = eInv.StockInDate.Add( time.Hour * time.Duration(24 * u.getRandomInteger(365, 730)) )

		product := new(ProductStruct)
		prodIdParts := strings.Split(inv, "--")

		product.Desc = prodIdParts[0]
		product.Id = prodIdParts[1]
		product.BatchId = fmt.Sprintf("%v-%06d", product.Id, u.getRandomInteger(1, 10))
		eInv.Product = *product

		// random get location
		iLocIdx := u.getRandomInteger(0, len(u.locationList))
		loc := u.locationList[iLocIdx]

		location := new(LocationStruct)
		location.Id = loc.ID
		location.Name = loc.Name
		location.PostCode = loc.Postcode
		location.Lat = loc.Lat
		location.Lng = loc.Lng
		eInv.Location = *location

		inventoryList = append(inventoryList, *eInv)
	}
	return inventoryList
}
// reusable method for generating Inventory and Sales entries.
// The return list contains the inventory information for building both types of trx
func (u *generatorUtil) getInventoryList(source, filename string) (invList []string) {
	// a. check if a previous run has created teh inventory_base.txt (no need to re-parse the whole Online sales.csv again
	baseFilename := fmt.Sprintf("%v%v%v", source, string(os.PathSeparator), srcInventoryBaseFilename)
	_, err := os.Stat(baseFilename)
	if err == nil || os.IsExist(err) {
		// load the contents back to memory and start working...
		fHandle, err := os.OpenFile(baseFilename, os.O_RDONLY, 0755)
		if err != nil {
			panic(err)
		}
		defer fHandle.Close()

		bContents, err := ioutil.ReadAll(fHandle)
		if err != nil {
			panic(err)
		}
		invList = strings.Split(string(bContents), ",")

	} else {
		sourceFile := fmt.Sprintf("%v%v%v", source, string(os.PathSeparator), srcOnlineSalesCsv)
		invList, err = u.parseSourceOnlineSalesCsv(sourceFile)
		if err != nil {
			panic(err)
		}
		finalContent := strings.Join(invList, ",")
		// write to base file
		err = ioutil.WriteFile(baseFilename, []byte(finalContent), 0755)
		if err != nil {
			panic(err)
		}
	}
	return
}

// parsing the source file sourceOnlineSales.csv
func (u *generatorUtil) parseSourceOnlineSalesCsv(filename string) (invList []string, err error) {
	fHandle, err := os.OpenFile(filename, os.O_RDONLY, 0755)
	if err != nil {
		return
	}
	defer fHandle.Close()

	invMap := make(map[string]bool)

	scanner := bufio.NewScanner(fHandle)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")
		if len(parts) > 3 {
			// only story unique values
			iVal := parts[2]
			if strings.Compare("Description", iVal) == -1 && invMap[iVal] == false {
				invMap[iVal] = true
				invList = append(invList, iVal)
			}
		}
	}
	// add back the generated Id
	prodIdMap := make(map[string]bool)
	for idx, prod := range invList {
		prodId := u.getRandomId("", idx)
		// check uniqueness?
		if prodIdMap[prodId] == false {
			prodIdMap[prodId] = true
			invList[idx] = prod + "--" +prodId
		} else {
			for true {
				prodId := u.getRandomId("", idx)
				if prodIdMap[prodId] == false {
					prodIdMap[prodId] = true
					invList[idx] = prod + "--" + prodId
					break
				}
			}
		}
	}
	return
}

func (u *generatorUtil) getRandomInteger(lower, upper int) (value int) {
	return rand.Intn(upper - lower) + lower
}
func (u *generatorUtil) getRandomFloat32(lower, upper float32) (value float32) {
	fV := rand.Float32()*(upper - lower) + lower
	return float32(math.Round(float64(fV) * 100) / 100)
}
func (u *generatorUtil) getRandomDate(lower, upper int) time.Time {
	cDate := time.Now()
	iDays := u.getRandomInteger(lower, upper)

	cDate = cDate.Add(-1 * (time.Hour * time.Duration(24*iDays)) )
	return cDate
}
func (u *generatorUtil) getRandomDateWithin24Hours() time.Time {
	cDate := time.Now()
	randMinutes := rand.Intn(24 * 60)

	cDate = cDate.Add(-1 * (time.Minute * time.Duration(randMinutes)) )
	return cDate
}
func (u *generatorUtil) getRandomId(category string, seed int) (id string) {
	d := int(time.Now().UnixNano())
	id = fmt.Sprintf("%v", int(math.Round(rand.Float64() * float64(d)))+seed )

	return
}

// InvoiceNo,StockCode,Description,Quantity,InvoiceDate,UnitPrice,CustomerID,Country (Online Sales.csv)
// Invoice ID,Branch,City,Customer type,Gender,Product line,Unit price,Quantity,Tax 5%,Total,Date,Time,Payment,cogs,gross margin percentage,gross income,Rating (supermarket_sales.csv)


func (u *generatorUtil) GenSalesTrx() {

}

// load the placemart info - a.k.a. location list
func (u *generatorUtil) loadPlacemartList(source, file string) (locations []PlacemarkStruct) {
	fname := fmt.Sprintf("%v%v%v_prepared.json", source, string(os.PathSeparator), file)
	_, err := os.Stat(fname)
	if err != nil && os.IsNotExist(err) {
		panic(err)
	}
	fHandle, err := os.OpenFile(fname, os.O_RDONLY, 0755)
	CommonPanic(err)
	defer fHandle.Close()

	bContent, err := ioutil.ReadAll(fHandle)
	CommonPanic(err)

	err = json.Unmarshal(bContent, &locations)
	CommonPanic(err)

	return
}


func (u *generatorUtil) generateSalesTrx(source, filename string, size int32) (salesList []SalesTrxStruct) {
	// check if inventory list and location list ready
	if len(u.invList) == 0 {
		u.invList = u.getInventoryList(source, filename)
	}
	if len(u.locationList) == 0 {
		u.locationList = u.loadPlacemartList(source, filename)
	}

	for i:=0; i<int(size); i++ {
		sS := new(SalesTrxStruct)
		sS.Date = u.getRandomDateWithin24Hours()
		sS.SellingPrice = u.getRandomFloat32(20, 160)
		sS.Quantity = int32(u.getRandomInteger(1, 20))

		sProd := new(ProductStruct)
		invIdx := u.getRandomInteger(0, len(u.invList))
		invAtIdx :=  u.invList[invIdx]
		parts := strings.Split(invAtIdx, "--")
		sProd.Id = parts[1]
		sProd.Desc = parts[0]
		sProd.BatchId = fmt.Sprintf("%v-%06d", sProd.Id, u.getRandomInteger(1, 10))
		sS.Product = *sProd

		sClient := new(ClientStruct)
		sClient.Id = fmt.Sprintf("%06d", u.getRandomInteger(30, 67301))
		// pick a random clientDemo
		cDemo := u.clientDemoList[u.getRandomInteger(0, len(u.clientDemoList))]
		sClient.Name = fmt.Sprintf("%v %v", cDemo.Name, cDemo.Surname)
		sClient.Gender = cDemo.Gender
		// pick a random occupation
		oDemo := u.occupationList[u.getRandomInteger(0, len(u.occupationList))]
		sClient.Occupation = oDemo
		sS.Client = *sClient

		sLoc := new(LocationStruct)
		lDemo := u.locationList[u.getRandomInteger(0, len(u.locationList))]
		sLoc.Name = lDemo.Name
		sLoc.Id = lDemo.ID
		sLoc.PostCode = lDemo.Postcode
		sLoc.Lat = lDemo.Lat
		sLoc.Lng = lDemo.Lng
		sS.Location = *sLoc

		salesList = append(salesList, *sS)
	}
	return
}

// models

type InventoryTrxStruct struct {
	StockInCost float32 `json:"stock_in_cost"`  		// range : 20 ~ 160
	StockInQuantity int32 `json:"stock_in_quantity"`	// range : 500 ~ 10000
	StockInDate time.Time `json:"stock_in_date"`		// range : 180 to 365 days earlier
	ExpiryDate time.Time `json:"expiry_date"`			// above date + 365 ~ 730 days

	Product ProductStruct `json:"product"`

	Location LocationStruct `json:"location"`
}

type LocationStruct struct {
	Id string `json:"id"`
	Name string `json:"name"`
	PostCode string `json:"post_code"`
	Lat float32 `json:"lat"`
	Lng float32 `json:"lng"`
}

type ProductStruct struct {
	Id string `json:"id"`
	Desc string `json:"desc"`
	BatchId string `json:"batch_id"`
}

type SalesTrxStruct struct {
	Date time.Time `json:"date"`				// within 24 hours of current time
	SellingPrice float32 `json:"selling_price"`	// random 20 ~ 160
	Quantity int32 `json:"quantity"`			// random 1 ~ 20

	Product ProductStruct `json:"product"`

	Client ClientStruct `json:"client"`

	Location LocationStruct `json:"location"`
}

type ClientStruct struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Gender string `json:"gender"`
	Occupation string `json:"occupation"`
}

type clientDemo struct {
	Name string `json:"name"`
	Surname string `json:"surname"`
	Gender string `json:"gender"`
}

// response model to encapsulate all generated entries
type EntryResponse struct {
	Profile string
	InventoryList []InventoryTrxStruct
	SalesList []SalesTrxStruct
}

// prepare the static data for the generation
func (u *generatorUtil) prepareRandomData() {
	/*clientDemoList = []clientDemo{
		clientDemo{ name: "a" },
	}*/

	// get random name + gender through api https://uinames.com/api/?amount=100
	resp, err := http.Get("https://uinames.com/api/?amount=200")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	bContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bContent, &u.clientDemoList)
// TODO: add also the clientId... (in the new phase)
	if err != nil {
		panic(err)
	}

	u.occupationList = []string{
		"accountant",
		"actor",
		"actuary",
		"adhesive bonding machine tender",
		"adjudicator",
		"administrative assistant",
		"administrative services manager",
		"adult education teacher",
		"advertising manager",
		"advertising sales agent",
		"aerobics instructor",
		"aerospace engineer",
		"aerospace engineering technician",
		"agent",
		"agricultural engineer",
		"agricultural equipment operator",
		"agricultural grader",
		"agricultural inspector",
		"agricultural manager",
		"agricultural sciences teacher",
		"agricultural sorter",
		"agricultural technician",
		"agricultural worker",
		"air conditioning installer",
		"air conditioning mechanic",
		"air traffic controller",
		"aircraft cargo handling supervisor",
		"aircraft mechanic",
		"aircraft service technician",
		"airline copilot",
		"airline pilot",
		"ambulance dispatcher",
		"ambulance driver",
		"amusement machine servicer",
		"anesthesiologist",
		"animal breeder",
		"animal control worker",
		"animal scientist",
		"animal trainer",
		"animator",
		"answering service operator",
		"anthropologist",
		"apparel patternmaker",
		"apparel worker",
		"arbitrator",
		"archeologist",
		"architect",
		"architectural drafter",
		"architectural manager",
		"archivist",
		"art director",
		"art teacher",
		"artist",
		"assembler",
		"astronomer",
		"athlete",
		"athletic trainer",
		"ATM machine repairer",
		"atmospheric scientist",
		"attendant",
		"audio and video equipment technician",
		"audio-visual and multimedia collections specialist",
		"audiologist",
		"auditor",
		"author",
		"auto damage insurance appraiser",
		"automotive and watercraft service attendant",
		"automotive glass installer",
		"automotive mechanic",
		"avionics technician",
		"baggage porter",
		"bailiff",
		"baker",
		"barback",
		"barber",
		"bartender",
		"basic education teacher",
		"behavioral disorder counselor",
		"bellhop",
		"bench carpenter",
		"bicycle repairer",
		"bill and account collector",
		"billing and posting clerk",
		"biochemist",
		"biological technician",
		"biomedical engineer",
		"biophysicist",
		"blaster",
		"blending machine operator",
		"blockmason",
		"boiler operator",
		"boilermaker",
		"bookkeeper",
		"boring machine tool tender",
		"brazer",
		"brickmason",
		"bridge and lock tender",
		"broadcast news analyst",
		"broadcast technician",
		"brokerage clerk",
		"budget analyst",
		"building inspector",
		"bus mechanic",
		"butcher",
		"buyer",
		"cabinetmaker",
		"cafeteria attendant",
		"cafeteria cook",
		"camera operator",
		"camera repairer",
		"cardiovascular technician",
		"cargo agent",
		"carpenter",
		"carpet installer",
		"cartographer",
		"cashier",
		"caster",
		"ceiling tile installer",
		"cellular equipment installer",
		"cement mason",
		"channeling machine operator",
		"chauffeur",
		"checker",
		"chef",
		"chemical engineer",
		"chemical plant operator",
		"chemist",
		"chemistry teacher",
		"chief executive",
		"child social worker",
		"childcare worker",
		"chiropractor",
		"choreographer",
		"civil drafter",
		"civil engineer",
		"civil engineering technician",
		"claims adjuster",
		"claims examiner",
		"claims investigator",
		"cleaner",
		"clinical laboratory technician",
		"clinical laboratory technologist",
		"clinical psychologist",
		"coating worker",
		"coatroom attendant",
		"coil finisher",
		"coil taper",
		"coil winder",
		"coin machine servicer",
		"commercial diver",
		"commercial pilot",
		"commodities sales agent",
		"communications equipment operator",
		"communications teacher",
		"community association manager",
		"community service manager",
		"compensation and benefits manager",
		"compliance officer",
		"composer",
		"computer hardware engineer",
		"computer network architect",
		"computer operator",
		"computer programmer",
		"computer science teacher",
		"computer support specialist",
		"computer systems administrator",
		"computer systems analyst",
		"concierge",
		"conciliator",
		"concrete finisher",
		"conservation science teacher",
		"conservation scientist",
		"conservation worker",
		"conservator",
		"construction inspector",
		"construction manager",
		"construction painter",
		"construction worker",
		"continuous mining machine operator",
		"convention planner",
		"conveyor operator",
		"cook",
		"cooling equipment operator",
		"copy marker",
		"correctional officer",
		"correctional treatment specialist",
		"correspondence clerk",
		"correspondent",
		"cosmetologist",
		"cost estimator",
		"costume attendant",
		"counseling psychologist",
		"counselor",
		"courier",
		"court reporter",
		"craft artist",
		"crane operator",
		"credit analyst",
		"credit checker",
		"credit counselor",
		"criminal investigator",
		"criminal justice teacher",
		"crossing guard",
		"curator",
		"custom sewer",
		"customer service representative",
		"cutter",
		"cutting machine operator",
		"dancer",
		"data entry keyer",
		"database administrator",
		"decorating worker",
		"delivery services driver",
		"demonstrator",
		"dental assistant",
		"dental hygienist",
		"dental laboratory technician",
		"dentist",
		"derrick operator",
		"designer",
		"desktop publisher",
		"detective",
		"diagnostic medical sonographer",
		"die maker",
		"diesel engine specialist",
		"dietetic technician",
		"dietitian",
		"dinkey operator",
		"director",
		"dishwasher",
		"dispatcher",
		"door-to-door sales worker",
		"drafter",
		"dragline operator",
		"drama teacher",
		"dredge operator",
		"dressing room attendant",
		"dressmaker",
		"drier operator",
		"drilling machine tool operator",
		"dry-cleaning worker",
		"drywall installer",
		"dyeing machine operator",
		"earth driller",
		"economics teacher",
		"economist",
		"editor",
		"education administrator",
		"electric motor repairer",
		"electrical electronics drafter",
		"electrical engineer",
		"electrical equipment assembler",
		"electrical installer",
		"electrical power-line installer",
		"electrician",
		"electro-mechanical technician",
		"elementary school teacher",
		"elevator installer",
		"elevator repairer",
		"embalmer",
		"emergency management director",
		"emergency medical technician",
		"engine assembler",
		"engineer",
		"engineering manager",
		"engineering teacher",
		"english language teacher",
		"engraver",
		"entertainment attendant",
		"environmental engineer",
		"environmental science teacher",
		"environmental scientist",
		"epidemiologist",
		"escort",
		"etcher",
		"event planner",
		"excavating operator",
		"executive administrative assistant",
		"executive secretary",
		"exhibit designer",
		"expediting clerk",
		"explosives worker",
		"extraction worker",
		"fabric mender",
		"fabric patternmaker",
		"fabricator",
		"faller",
		"family practitioner",
		"family social worker",
		"family therapist",
		"farm advisor",
		"farm equipment mechanic",
		"farm labor contractor",
		"farmer",
		"farmworker",
		"fashion designer",
		"fast food cook",
		"fence erector",
		"fiberglass fabricator",
		"fiberglass laminator",
		"file clerk",
		"filling machine operator",
		"film and video editor",
		"financial analyst",
		"financial examiner",
		"financial manager",
		"financial services sales agent",
		"fine artist",
		"fire alarm system installer",
		"fire dispatcher",
		"fire inspector",
		"fire investigator",
		"firefighter",
		"fish and game warden",
		"fish cutter",
		"fish trimmer",
		"fisher",
		"fitness studies teacher",
		"fitness trainer",
		"flight attendant",
		"floor finisher",
		"floor layer",
		"floor sander",
		"floral designer",
		"food batchmaker",
		"food cooking machine operator",
		"food preparation worker",
		"food science technician",
		"food scientist",
		"food server",
		"food service manager",
		"food technologist",
		"foreign language teacher",
		"foreign literature teacher",
		"forensic science technician",
		"forest fire inspector",
		"forest fire prevention specialist",
		"forest worker",
		"forester",
		"forestry teacher",
		"forging machine setter",
		"foundry coremaker",
		"freight agent",
		"freight mover",
		"fundraising manager",
		"funeral attendant",
		"funeral director",
		"funeral service manager",
		"furnace operator",
		"furnishings worker",
		"furniture finisher",
		"gaming booth cashier",
		"gaming cage worker",
		"gaming change person",
		"gaming dealer",
		"gaming investigator",
		"gaming manager",
		"gaming surveillance officer",
		"garment mender",
		"garment presser",
		"gas compressor",
		"gas plant operator",
		"gas pumping station operator",
		"general manager",
		"general practitioner",
		"geographer",
		"geography teacher",
		"geological engineer",
		"geological technician",
		"geoscientist",
		"glazier",
		"government program eligibility interviewer",
		"graduate teaching assistant",
		"graphic designer",
		"groundskeeper",
		"groundskeeping worker",
		"gynecologist",
		"hairdresser",
		"hairstylist",
		"hand grinding worker",
		"hand laborer",
		"hand packager",
		"hand packer",
		"hand polishing worker",
		"hand sewer",
		"hazardous materials removal worker",
		"head cook",
		"health and safety engineer",
		"health educator",
		"health information technician",
		"health services manager",
		"health specialties teacher",
		"healthcare social worker",
		"hearing officer",
		"heat treating equipment setter",
		"heating installer",
		"heating mechanic",
		"heavy truck driver",
		"highway maintenance worker",
		"historian",
		"history teacher",
		"hoist and winch operator",
		"home appliance repairer",
		"home economics teacher",
		"home entertainment installer",
		"home health aide",
		"home management advisor",
		"host",
		"hostess",
		"hostler",
		"hotel desk clerk",
		"housekeeping cleaner",
		"human resources assistant",
		"human resources manager",
		"human service assistant",
		"hunter",
		"hydrologist",
		"illustrator",
		"industrial designer",
		"industrial engineer",
		"industrial engineering technician",
		"industrial machinery mechanic",
		"industrial production manager",
		"industrial truck operator",
		"industrial-organizational psychologist",
		"information clerk",
		"information research scientist",
		"information security analyst",
		"information systems manager",
		"inspector",
		"instructional coordinator",
		"instructor",
		"insulation worker",
		"insurance claims clerk",
		"insurance sales agent",
		"insurance underwriter",
		"intercity bus driver",
		"interior designer",
		"internist",
		"interpreter",
		"interviewer",
		"investigator",
		"jailer",
		"janitor",
		"jeweler",
		"judge",
		"judicial law clerk",
		"kettle operator",
		"kiln operator",
		"kindergarten teacher",
		"laboratory animal caretaker",
		"landscape architect",
		"landscaping worker",
		"lathe setter",
		"laundry worker",
		"law enforcement teacher",
		"law teacher",
		"lawyer",
		"layout worker",
		"leather worker",
		"legal assistant",
		"legal secretary",
		"legislator",
		"librarian",
		"library assistant",
		"library science teacher",
		"library technician",
		"licensed practical nurse",
		"licensed vocational nurse",
		"life scientist",
		"lifeguard",
		"light truck driver",
		"line installer",
		"literacy teacher",
		"literature teacher",
		"loading machine operator",
		"loan clerk",
		"loan interviewer",
		"loan officer",
		"lobby attendant",
		"locker room attendant",
		"locksmith",
		"locomotive engineer",
		"locomotive firer",
		"lodging manager",
		"log grader",
		"logging equipment operator",
		"logistician",
		"machine feeder",
		"machinist",
		"magistrate judge",
		"magistrate",
		"maid",
		"mail clerk",
		"mail machine operator",
		"mail superintendent",
		"maintenance painter",
		"maintenance worker",
		"makeup artist",
		"management analyst",
		"manicurist",
		"manufactured building installer",
		"mapping technician",
		"marble setter",
		"marine engineer",
		"marine oiler",
		"market research analyst",
		"marketing manager",
		"marketing specialist",
		"marriage therapist",
		"massage therapist",
		"material mover",
		"materials engineer",
		"materials scientist",
		"mathematical science teacher",
		"mathematical technician",
		"mathematician",
		"maxillofacial surgeon",
		"measurer",
		"meat cutter",
		"meat packer",
		"meat trimmer",
		"mechanical door repairer",
		"mechanical drafter",
		"mechanical engineer",
		"mechanical engineering technician",
		"mediator",
		"medical appliance technician",
		"medical assistant",
		"medical equipment preparer",
		"medical equipment repairer",
		"medical laboratory technician",
		"medical laboratory technologist",
		"medical records technician",
		"medical scientist",
		"medical secretary",
		"medical services manager",
		"medical transcriptionist",
		"meeting planner",
		"mental health counselor",
		"mental health social worker",
		"merchandise displayer",
		"messenger",
		"metal caster",
		"metal patternmaker",
		"metal pickling operator",
		"metal pourer",
		"metal worker",
		"metal-refining furnace operator",
		"metal-refining furnace tender",
		"meter reader",
		"microbiologist",
		"middle school teacher",
		"milling machine setter",
		"millwright",
		"mine cutting machine operator",
		"mine shuttle car operator",
		"mining engineer",
		"mining safety engineer",
		"mining safety inspector",
		"mining service unit operator",
		"mixing machine setter",
		"mobile heavy equipment mechanic",
		"mobile home installer",
		"model maker",
		"model",
		"molder",
		"mortician",
		"motel desk clerk",
		"motion picture projectionist",
		"motorboat mechanic",
		"motorboat operator",
		"motorboat service technician",
		"motorcycle mechanic",
		"multimedia artist",
		"museum technician",
		"music director",
		"music teacher",
		"musical instrument repairer",
		"musician",
		"natural sciences manager",
		"naval architect",
		"network systems administrator",
		"new accounts clerk",
		"news vendor",
		"nonfarm animal caretaker",
		"nuclear engineer",
		"nuclear medicine technologist",
		"nuclear power reactor operator",
		"nuclear technician",
		"nursing aide",
		"nursing instructor",
		"nursing teacher",
		"nutritionist",
		"obstetrician",
		"occupational health and safety specialist",
		"occupational health and safety technician",
		"occupational therapist",
		"occupational therapy aide",
		"occupational therapy assistant",
		"offbearer",
		"office clerk",
		"office machine operator",
		"operating engineer",
		"operations manager",
		"operations research analyst",
		"ophthalmic laboratory technician",
		"optician",
		"optometrist",
		"oral surgeon",
		"order clerk",
		"order filler",
		"orderly",
		"ordnance handling expert",
		"orthodontist",
		"orthotist",
		"outdoor power equipment mechanic",
		"oven operator",
		"packaging machine operator",
		"painter ",
		"painting worker",
		"paper goods machine setter",
		"paperhanger",
		"paralegal",
		"paramedic",
		"parking enforcement worker",
		"parking lot attendant",
		"parts salesperson",
		"paving equipment operator",
		"payroll clerk",
		"pediatrician",
		"pedicurist",
		"personal care aide",
		"personal chef",
		"personal financial advisor",
		"pest control worker",
		"pesticide applicator",
		"pesticide handler",
		"pesticide sprayer",
		"petroleum engineer",
		"petroleum gauger",
		"petroleum pump system operator",
		"petroleum refinery operator",
		"petroleum technician",
		"pharmacist",
		"pharmacy aide",
		"pharmacy technician",
		"philosophy teacher",
		"photogrammetrist",
		"photographer",
		"photographic process worker",
		"photographic processing machine operator",
		"physical therapist aide",
		"physical therapist assistant",
		"physical therapist",
		"physician assistant",
		"physician",
		"physicist",
		"physics teacher",
		"pile-driver operator",
		"pipefitter",
		"pipelayer",
		"planing machine operator",
		"planning clerk",
		"plant operator",
		"plant scientist",
		"plasterer",
		"plastic patternmaker",
		"plastic worker",
		"plumber",
		"podiatrist",
		"police dispatcher",
		"police officer",
		"policy processing clerk",
		"political science teacher",
		"political scientist",
		"postal service clerk",
		"postal service mail carrier",
		"postal service mail processing machine operator",
		"postal service mail processor",
		"postal service mail sorter",
		"postmaster",
		"postsecondary teacher",
		"poultry cutter",
		"poultry trimmer",
		"power dispatcher",
		"power distributor",
		"power plant operator",
		"power tool repairer",
		"precious stone worker",
		"precision instrument repairer",
		"prepress technician",
		"preschool teacher",
		"priest",
		"print binding worker",
		"printing press operator",
		"private detective",
		"probation officer",
		"procurement clerk",
		"producer",
		"product promoter",
		"production clerk",
		"production occupation",
		"proofreader",
		"property manager",
		"prosthetist",
		"prosthodontist",
		"psychiatric aide",
		"psychiatric technician",
		"psychiatrist",
		"psychologist",
		"psychology teacher",
		"public relations manager",
		"public relations specialist",
		"pump operator",
		"purchasing agent",
		"purchasing manager",
		"radiation therapist",
		"radio announcer",
		"radio equipment installer",
		"radio operator",
		"radiologic technician",
		"radiologic technologist",
		"rail car repairer",
		"rail transportation worker",
		"rail yard engineer",
		"rail-track laying equipment operator",
		"railroad brake operator",
		"railroad conductor",
		"railroad police",
		"rancher",
		"real estate appraiser",
		"real estate broker",
		"real estate manager",
		"real estate sales agent",
		"receiving clerk",
		"receptionist",
		"record clerk",
		"recreation teacher",
		"recreation worker",
		"recreational therapist",
		"recreational vehicle service technician",
		"recyclable material collector",
		"referee",
		"refractory materials repairer",
		"refrigeration installer",
		"refrigeration mechanic",
		"refuse collector",
		"regional planner",
		"registered nurse",
		"rehabilitation counselor",
		"reinforcing iron worker",
		"reinforcing rebar worker",
		"religion teacher",
		"religious activities director",
		"religious worker",
		"rental clerk",
		"repair worker",
		"reporter",
		"residential advisor",
		"resort desk clerk",
		"respiratory therapist",
		"respiratory therapy technician",
		"retail buyer",
		"retail salesperson",
		"revenue agent",
		"rigger",
		"rock splitter",
		"rolling machine tender",
		"roof bolter",
		"roofer",
		"rotary drill operator",
		"roustabout",
		"safe repairer",
		"sailor",
		"sales engineer",
		"sales manager",
		"sales representative",
		"sampler",
		"sawing machine operator",
		"scaler",
		"school bus driver",
		"school psychologist",
		"school social worker",
		"scout leader",
		"sculptor",
		"secondary education teacher",
		"secondary school teacher",
		"secretary",
		"securities sales agent",
		"security guard",
		"security system installer",
		"segmental paver",
		"self-enrichment education teacher",
		"semiconductor processor",
		"septic tank servicer",
		"set designer",
		"sewer pipe cleaner",
		"sewing machine operator",
		"shampooer",
		"shaper",
		"sheet metal worker",
		"sheriff's patrol officer",
		"ship captain",
		"ship engineer",
		"ship loader",
		"shipmate",
		"shipping clerk",
		"shoe machine operator",
		"shoe worker",
		"short order cook",
		"signal operator",
		"signal repairer",
		"singer",
		"ski patrol",
		"skincare specialist",
		"slaughterer",
		"slicing machine tender",
		"slot supervisor",
		"social science research assistant",
		"social sciences teacher",
		"social scientist",
		"social service assistant",
		"social service manager",
		"social work teacher",
		"social worker",
		"sociologist",
		"sociology teacher",
		"software developer",
		"software engineer",
		"soil scientist",
		"solderer",
		"sorter",
		"sound engineering technician",
		"space scientist",
		"special education teacher",
		"speech-language pathologist",
		"sports book runner",
		"sports entertainer",
		"sports performer",
		"stationary engineer",
		"statistical assistant",
		"statistician",
		"steamfitter",
		"stock clerk",
		"stock mover",
		"stonemason",
		"street vendor",
		"streetcar operator",
		"structural iron worker",
		"structural metal fabricator",
		"structural metal fitter",
		"structural steel worker",
		"stucco mason",
		"substance abuse counselor",
		"substance abuse social worker",
		"subway operator",
		"surfacing equipment operator",
		"surgeon",
		"surgical technologist",
		"survey researcher",
		"surveying technician",
		"surveyor",
		"switch operator",
		"switchboard operator",
		"tailor",
		"tamping equipment operator",
		"tank car loader",
		"taper",
		"tax collector",
		"tax examiner",
		"tax preparer",
		"taxi driver",
		"teacher assistant",
		"teacher",
		"team assembler",
		"technical writer",
		"telecommunications equipment installer",
		"telemarketer",
		"telephone operator",
		"television announcer",
		"teller",
		"terrazzo finisher",
		"terrazzo worker",
		"tester",
		"textile bleaching operator",
		"textile cutting machine setter",
		"textile knitting machine setter",
		"textile presser",
		"textile worker",
		"therapist",
		"ticket agent",
		"ticket taker",
		"tile setter",
		"timekeeping clerk",
		"timing device assembler",
		"tire builder",
		"tire changer",
		"tire repairer",
		"title abstractor",
		"title examiner",
		"title searcher",
		"tobacco roasting machine operator",
		"tool filer",
		"tool grinder",
		"tool maker",
		"tool sharpener",
		"tour guide",
		"tower equipment installer",
		"tower operator",
		"track switch repairer",
		"tractor operator",
		"tractor-trailer truck driver",
		"traffic clerk",
		"traffic technician",
		"training and development manager",
		"training and development specialist",
		"transit police",
		"translator",
		"transportation equipment painter",
		"transportation inspector",
		"transportation security screener",
		"transportation worker",
		"trapper",
		"travel agent",
		"travel clerk",
		"travel guide",
		"tree pruner",
		"tree trimmer",
		"trimmer",
		"truck loader",
		"truck mechanic",
		"tuner",
		"turning machine tool operator",
		"typist",
		"umpire",
		"undertaker",
		"upholsterer",
		"urban planner",
		"usher",
		"valve installer",
		"vending machine servicer",
		"veterinarian",
		"veterinary assistant",
		"veterinary technician",
		"vocational counselor",
		"vocational education teacher",
		"waiter",
		"waitress",
		"watch repairer",
		"water treatment plant operator",
		"weaving machine setter",
		"web developer",
		"weigher",
		"welder",
		"wellhead pumper",
		"wholesale buyer",
		"wildlife biologist",
		"window trimmer",
		"wood patternmaker",
		"woodworker",
		"word processor",
		"writer",
		"yardmaster",
		"zoologist",
	}
}

