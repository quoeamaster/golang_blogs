/*
Copyright © 2019 quo master

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
package app

import (
	"encoding/json"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/quoeamaster/golang_blogs/repo"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type PortraitApp struct {
	fileRepoService *repo.FileRepoService
}

func NewPortraitApp() (instance *PortraitApp) {
	instance = new(PortraitApp)
	err := instance.Init()
	if err != nil {
		panic(err)
	}
	return
}

func (p *PortraitApp) Init() (err error) {
	// create repo service
	p.fileRepoService = repo.NewFileRepoService()

	// setup REST api
	api := rest.NewApi()
	//api.Use(rest.DefaultDevStack...)
	api.Use([]rest.Middleware{
		&rest.AccessLogApacheMiddleware{},
		&rest.TimerMiddleware{},
		&rest.RecorderMiddleware{},
		&rest.PoweredByMiddleware{},
		&rest.RecoverMiddleware{
			EnableResponseStackTrace: true,
		},
	}...)

	router, err := rest.MakeRouter(
		rest.Post("/addPost", p.PostAddPost),
		rest.Post("/addComment", p.PostAddComment),
		rest.Get("/getPost/:id", p.GetPostById),
		rest.Get("/getRandom10Posts", p.GetRandom10Posts),
	)
	if err != nil {
		return err
	}

	api.SetApp(router)
	http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./webapp"))))
	// static pages example => http://localhost:8100/static/index.html (index.html was available under /webapp folder)

	err = http.ListenAndServe(":8100", nil)
	if err != nil {
		return
	}
	return
}


// involve multipart and form-data
func (p *PortraitApp) PostAddPost(w rest.ResponseWriter, req *rest.Request) {
	valueMap := make(map[string]interface{})

	defer req.Body.Close()
	/*
	 *	exception: no multipart boundary param in Content-Type means ...
	 * 	normal multipart/form-data should be like this where the "boundary" param is available for parsing
	 * 	multipart/form-data; boundary=----WebKitFormBoundaryIox6yH4Vs0P82Y5O
	 *
	 *	for jQuery; if you set the headers to Content-Type: multipart/form-data;
	 *	then the boundary would not be supplied at all; hence got this exception
	 *
	 *	jQuery bug has a workaround but setting BOTH: (check add.html => jQueryFileUpload)
	 *	1. enctype => multipart/form-data,
	 *	2. contentType => false
	 */
	defer req.Body.Close()
	if err := req.ParseMultipartForm(10 << 20); err != nil {
		valueMap["error"] = err.Error()
		w.WriteJson(valueMap)
		return
	}
	// create folder for portrait
	uuidFolder, err := p.fileRepoService.CreateFolder()
	if err != nil {
		valueMap["error"] = err.Error()
		w.WriteJson(valueMap)
		return
	}
	// write the uploaded photo file
	fileParts := req.MultipartForm.File["file"]
	if fileParts != nil && len(fileParts) > 0 {
		// since only 1 file should be uploaded via the param named "file"... get only the 1st file-part
		uploadedFile, err2 := fileParts[0].Open()
		if err2 != nil {
			valueMap["error"] = err2.Error()
			w.WriteJson(valueMap)
			return
		}
		bContent, err2 := ioutil.ReadAll(uploadedFile)
		if err2 != nil {
			valueMap["error"] = err2.Error()
			w.WriteJson(valueMap)
			return
		}
		err2 = p.fileRepoService.WriteFileFromBytes(bContent, uuidFolder, fileParts[0].Filename)
		if err2 != nil {
			valueMap["error"] = err2.Error()
			w.WriteJson(valueMap)
			return
		}
		defer uploadedFile.Close()
	}
	// write meta-info (desc, create_date, photographer, id, photo_location)
	model := new(PortraitModel)
	model.Id = uuidFolder
	cDate, err := time.Parse("2006-01-02", req.Form.Get("createDate"))
	if err != nil {
		cDate = time.Unix(0,0)
	}
	model.CreateDate = cDate
	model.Description = req.Form.Get("desc")
	model.Photographer = req.Form.Get("photographer")
	model.PhotoLocation = fileParts[0].Filename

	bContent, err := json.Marshal(model)
	if err != nil {
		valueMap["error"] = err.Error()
		w.WriteJson(valueMap)
		return
	}
	if err = p.fileRepoService.WriteFileFromBytes(bContent, uuidFolder, repo.FILE_INFO); err != nil {
		valueMap["error"] = err.Error()
		w.WriteJson(valueMap)
		return
	}
	// all good!
	valueMap["status"] = "success"
	valueMap["message"] = "portrait successfully created"
	valueMap["portrait_id"] = uuidFolder
	w.WriteJson(valueMap)
}

// model / struct for json marshalling purpose
type PortraitModel struct {
	Id string				`json:"id"`
	PhotoLocation string	`json:"photo_location"`
	CreateDate time.Time	`json:"create_date"`
	Description string		`json:"desc"`
	Photographer string		`json:"photographer"`
}

// involve form-data or request-body only
func (p *PortraitApp) PostAddComment(w rest.ResponseWriter, req *rest.Request) {
	valueMap := make(map[string]interface{})

	if err := req.ParseForm(); err != nil {
		valueMap["error"] = err.Error()
		w.WriteJson(valueMap)
		return
	}
	// append contents to comment file
	if err := p.fileRepoService.AppendStringToFile(req.Form.Get("comment"), req.Form.Get("id"), repo.FILE_COMMENT, repo.COMMENT_DELIMITER); err != nil {
		valueMap["error"] = err.Error()
		w.WriteJson(valueMap)
		return
	}
	valueMap["status"] = "success"
	valueMap["message"] = "comment successfully added~ Reload page to refresh the comments"
	w.WriteJson(valueMap)
}

// involve path-param
func (p *PortraitApp) GetPostById(w rest.ResponseWriter, req *rest.Request) {
	valueMap := make(map[string]interface{})

	portraitId := req.PathParams["id"]
	// parse-form
	if err := req.ParseForm(); err != nil {
		valueMap["error"] = err.Error()
		w.WriteJson(valueMap)
		return
	}
	// get comments by portrait_id
	if err, comments := p.fileRepoService.GetCommentsByPortraitId(portraitId); err != nil {
		valueMap["error"] = err.Error()
		w.WriteJson(valueMap)
		return
	} else {
		valueMap["comments"] = comments
		valueMap["status"] = "success"
		w.WriteJson(valueMap)
	}
}

// get a random top 10 post(s)
// expected json result =>
// { portraits: [ { id: "axdfd", photo_location: "abc.jpg", create_date: "2019-01-01", desc: "hi", photographer: "victor freeze" }, { ... } ] }
func (p *PortraitApp) GetRandom10Posts(w rest.ResponseWriter, req *rest.Request) {
	valueMap := make(map[string]interface{})

	// a. get random 10 portraits from repository
	// b. create the response json: portraits: [ { "a": "xx", "b": "yy" } ]
	folderList, err := p.fileRepoService.GetFolderList()
	if err != nil {
		valueMap["error"] = err.Error()
		w.WriteJson(valueMap)
		return
	}

	// random pick 10
	if len(folderList) == 0 {
		valueMap["portraits"] = []string{}
		w.WriteJson(valueMap)
		return

	} else if len(folderList) < 10 {
		err2, metaInfoList := p.fileRepoService.GetFolderInfo(folderList)
		if err2 != nil {
			valueMap["error"] = err2.Error()
			w.WriteJson(valueMap)
			return
		}
		valueMap["portraits"] = metaInfoList
		w.WriteJson(valueMap)
		return

	} else {
		// random pick 10
		folderMap := make(map[int]string)
		fInfoListLen := len(folderList)
		for {
			if len(folderMap) < 10 {
				idx := rand.Intn(fInfoListLen)
				mapVal := folderMap[idx]
				if strings.Compare(strings.Trim(mapVal, " "), "") == 0 {
					folderMap[idx] = folderList[idx]
				}
			} else {
				break
			}
		}
		if len(folderMap) > 0 {
			idx := 0
			randomList := make([]string, 10)
			for _, folderName := range folderMap {
				randomList[idx] = folderName
				idx++
			}

			err2, metaInfoList := p.fileRepoService.GetFolderInfo(randomList)
			if err2 != nil {
				valueMap["error"] = err2.Error()
				w.WriteJson(valueMap)
				return
			}
			valueMap["portraits"] = metaInfoList
			w.WriteJson(valueMap)
			return
		}
		// default
		valueMap["portraits"] = []string{}
		w.WriteJson(valueMap)
		return
	}
}