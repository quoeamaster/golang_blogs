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
	"strings"
)

const defaultRepoLocation = "other_resources/notes.json"


type App struct {
	// map [string] -> list of map [string] -> interface (any value)
	notes map[string][]map[string]interface{}
}

func NewApp() (appPtr *App) {
	appPtr = new(App)
	appPtr.notes = make(map[string][]map[string]interface{})

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
	var listOfMap []map[string]interface{}
	listing := n.notes[todayInString]

	// try / catch on ANY error happening later
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	if listing == nil {
		// create an array of interface{}
		listOfMap = make([]map[string]interface{}, 0)

	} else {
		// convert the object to array of map[string]interface{}
		//listOfMap = listing.([]map[string]interface{}) (no need for conversion... probably)
		listOfMap = listing
	}
	// append
	entry := make(map[string]interface{})
	entry["content"] = content

	listOfMap = append(listOfMap, entry)
	n.notes[todayInString] = listOfMap

	// write to file
	notesInString := n.notesRepoToString(n.notes)
	err = ioutil.WriteFile(defaultRepoLocation, []byte(notesInString), 0644)
	
	return
}

func (n *App) notesRepoToString(repo map[string][]map[string]interface{}) (value string) {
	var b strings.Builder

	fmt.Fprintf(&b, "{ ")
	for key, val := range repo {
		fmt.Fprintf(&b, "\"%v\": [", key)
		// list of map (notes content)
		for idx, noteObject := range val {
			if idx > 0 {
				fmt.Fprintf(&b, ", ")
			}
			fmt.Fprintf(&b, n.noteToString(noteObject, true))
		} // end -- for ( [] -> map[string]interface{} level)
		fmt.Fprintf(&b, "]")
	} // end -- for ( date -> []map[string]interface{} level)

	fmt.Fprintf(&b, " }")
	value = b.String()

	return
}

/**
 *	convert the note (map[string]interface{}) to string format
 */
func (n *App) noteToString(entry map[string]interface{}, needWrapping bool) (value string) {
	//[]byte(value)

	var b strings.Builder
	idx := 0

	if needWrapping {
		fmt.Fprint(&b, "{ ")
	}
	for key, val := range entry {
		if idx > 0 {
			fmt.Fprintf(&b, ", ")
		}
		// key part
		fmt.Fprintf(&b, "\"%v\": ", key)
		// value part
		switch val.(type) {
		case int, int8, int16, int32, int64:
			fmt.Fprintf(&b, "%v", val)
		case float32, float64:
			fmt.Fprintf(&b, "%v", val)
		default:
			fmt.Fprintf(&b, "\"%v\"", val)
		}
		idx++
	}
	if needWrapping {
		fmt.Fprint(&b, " }")
	}
	value = b.String()

	return
}