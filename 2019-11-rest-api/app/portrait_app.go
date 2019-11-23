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
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/quoeamaster/golang_blogs/repo"
	"net/http"
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
	api.Use(rest.DefaultDevStack...)

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
	err = http.ListenAndServe(":8100", api.MakeHandler())
	if err != nil {
		return
	}
	return
}


// involve multipart and form-data
func (p *PortraitApp) PostAddPost(w rest.ResponseWriter, req *rest.Request) {

}

// involve form-data or request-body only
func (p *PortraitApp) PostAddComment(w rest.ResponseWriter, req *rest.Request) {

}

// involve path-param
func (p *PortraitApp) GetPostById(w rest.ResponseWriter, req *rest.Request) {

}

// get a random top 10 post(s)
func (p *PortraitApp) GetRandom10Posts(w rest.ResponseWriter, req *rest.Request) {

}