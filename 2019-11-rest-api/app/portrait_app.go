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
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/quoeamaster/golang_blogs/repo"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
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

/*
type PortraitModel struct {
	PhotoLocation string
	CreateDate time.Time
	Description string
	Photographer string
}
*/

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
		rest.Post("/addComment", p.PostAddPost),
		rest.Get("/getPost/:id", p.GetPostById),
		rest.Get("/getRandom10Posts", p.GetRandom10Posts),
	)
	if err != nil {
		return err
	}

	api.SetApp(router)
	http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./webapp"))))

	err = http.ListenAndServe(":8100", nil)
	if err != nil {
		return
	}
	return
}


// involve multipart and form-data
func (p *PortraitApp) PostAddPost(w rest.ResponseWriter, req *rest.Request) {
	// valueMap := make(map[string]interface{})

	defer req.Body.Close()
	/*
	 *	exception: no multipart boundary param in Content-Type
	 */

	// need to parse everything by myself... orz

	bC, _ := ioutil.ReadAll(req.Body)
	fmt.Println(string(bC))

	reader, err := req.MultipartReader()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(req.ContentLength)
	fmt.Println(reader)


}

// involve form-data or request-body only
func (p *PortraitApp) PostAddComment(w rest.ResponseWriter, req *rest.Request) {

}

// involve path-param
func (p *PortraitApp) GetPostById(w rest.ResponseWriter, req *rest.Request) {

}

// get a random top 10 post(s)
// expected json result =>
// { portraits: [ { id: "axdfd", photo_location: "abc.jpg", create_date: "2019-01-01", desc: "hi", photographe: "victor freeze" }, { ... } ] }
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