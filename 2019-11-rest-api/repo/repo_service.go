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
	"strings"
)

const (
	REPO_DEF_PATH     = "webapp/portrait_repo/"
	FILE_INFO         = "info.me"
	FILE_COMMENT      = "comment.me"
	COMMENT_DELIMITER = "--delim--"
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
	// any repo connector setup here
	return
}

// generate a new UUID to act as the portrait's id
func (r *FileRepoService) generatePortraitId() (id string) {
	id = uuid.New()
	return
}

// create the folder (e.g. repo/{UUID} ) where the UUID is generated on demand
func (r *FileRepoService) CreateFolder() (folderName string, err error) {
	folderName = r.generatePortraitId()
	err = os.MkdirAll(REPO_DEF_PATH +  folderName, 0777)
	return
}

// write / overwrite content to a file; if you prefer to append contents instead...
// check the AppendStringToFile for string data
func (r *FileRepoService) WriteFileFromBytes(bContent []byte, folder, filename string) (err error) {
	finalFilename := REPO_DEF_PATH + folder + "/" + filename
	err = ioutil.WriteFile(finalFilename, bContent, 0777)
	return
}

// append contents to an existing file; if the file is not available, create it first
func (r *FileRepoService) AppendStringToFile(content, folder, filename, delimiter string) (err error) {
	finalFilename := REPO_DEF_PATH + folder + "/" + filename

	// open file for append
	fileHandle, err := os.OpenFile(finalFilename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	defer fileHandle.Close()
	if err != nil {
		return
	}
	_, err = fileHandle.WriteString(content + delimiter)
	if err != nil {
		return
	}
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
			// the "interface{}" must be an address reference
			err3 = json.Unmarshal(bContent, &metaMap)
			if err3 != nil {
				err = err3
				return
			}
			metaList = append(metaList, metaMap)
		}
	}
	return
}

func (r *FileRepoService) GetCommentsByPortraitId(id string) (err error, comments []string) {
	comments = []string{}

	commentFileLocation := REPO_DEF_PATH + id + "/" + FILE_COMMENT
	if fInfo, err3 := os.Stat(commentFileLocation); !os.IsNotExist(err3) && !fInfo.IsDir() {
		if bContent, err2 := ioutil.ReadFile(commentFileLocation); err2 != nil {
			return
		} else {
			comments = strings.Split(string(bContent), COMMENT_DELIMITER)
		}
	}
	// fmt.Println("length of comments:", len(comments), comments)
	return
}