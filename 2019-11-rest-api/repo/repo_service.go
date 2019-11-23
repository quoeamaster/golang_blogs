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
package repo

import (
	"encoding/json"
	"github.com/pborman/uuid"
	"io/ioutil"
	"os"
)

const (
	REPO_DEF_PATH = "repo/"
	FILE_INFO     = "info.me"
	FILE_COMMENT  = "comment.me"
	FILE_IMG      = "portrait_main."
)

type FileRepoService struct {

}

func NewFileRepoService() (instance *FileRepoService) {
	instance = new(FileRepoService)
	if err := instance.Init(); err != nil {
		panic(err)
	}
	return
}

func (r *FileRepoService) Init() (err error) {

	return
}

// generate a new UUID to act as the portrait's id
func (r *FileRepoService) generatePortraitId() (id string) {
	id = uuid.New()
	return
}

func (r *FileRepoService) CreateFolder() (folderName string, err error) {

	return
}

func (r *FileRepoService) WriteFileFromBytes(bContent []byte, folder, filename string) (err error) {
	return
}

// return the available folder-name list
func (r *FileRepoService) GetFolderList() (list []string, err error) {
	list = make([]string, 0)

	fileInfos, err := ioutil.ReadDir(REPO_DEF_PATH)
	if err != nil {
		return
	}
	// filter out only "folders"
	for _, fInfo := range fileInfos {
		if fInfo.IsDir() {
			list = append(list, fInfo.Name())
		}
	}
	return
}

// return the folder-info based on the folder-list provided
func (r *FileRepoService) GetFolderInfo(list []string) (err error, metaList []map[string]interface{}) {
	metaList = make([]map[string]interface{}, 0)

	for _, folderName := range list {
		// create the meta-info file
		metaFileLocation := REPO_DEF_PATH + folderName + "/" + FILE_INFO
		if _, err2 := os.Stat(metaFileLocation); !os.IsNotExist(err2) {
			bContent, err3 := ioutil.ReadFile(metaFileLocation)
			if err3 != nil {
				err = err3
				return
			}
			metaMap := make(map[string]interface{})
			err3 = json.Unmarshal(bContent, metaMap)
			if err3 != nil {
				err = err3
				return
			}
			metaList = append(metaList, metaMap)
		}
	}
	return
}