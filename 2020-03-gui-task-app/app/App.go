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
package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const defaultRepoLocation = "other_resources/notes.json"


type App struct {
	notes map[string]interface{}
}

func NewApp() (appPtr *App) {
	appPtr = new(App)
	appPtr.notes = make(map[string]interface{})

	err := appPtr.loadNotesFromRepo()
	if err != nil {
		panic(err)
	}
fmt.Println(appPtr.notes, "###")
	return
}

// load the notes repos / data-structure
func (n *App) loadNotesFromRepo() (err error) {
	// is the targeted repo file available?
	_, err2 := os.Stat(defaultRepoLocation)
	if os.IsNotExist(err2) == true {
		err = ioutil.WriteFile(defaultRepoLocation, nil, 0644)
		if err != nil {
			return
		}
	} else {
		// load the repo file
		bContent, err2 := ioutil.ReadFile(defaultRepoLocation)
		if err2 != nil {
			err = err2
			return
		}

		if len(bContent) > 0 {
			// load / unmarshal
			err2 = json.Unmarshal(bContent, &n.notes)
			if err2 != nil {
				err = err2
				return
			}
		}
	}
	return
}


/**
 *	onStart event / hook
 */
func (n *App) OnStart() (err error) {
	fmt.Println("app ready~")
	return
}


/* ---------------------- */
/*   business event(s)    */
/* ---------------------- */

/**
 *	create note / task
 */
func (n *App) OnCreateNoteTask(content, todayInString string) (err error) {
	//fmt.Println("tbd - save the note / task:", content, todayInString)
	listing := n.notes[todayInString]
	if listing == nil {
		// TODO: create an array of interface{}

	} else {
		// append
		// TODO: convert the object to array of interface{}
	}


	return
}