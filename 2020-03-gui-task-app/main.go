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
package main

import (
	"fmt"
	"github.com/quoeamaster/golang_blogs/app"
	"github.com/zserge/lorca"
	"net"
	"net/http"
	"os"
	"os/signal"
)

const eventNotesAllLoaded = "go-notes-all-loaded"

func main()  {
	args := []string{}
	prepareArgsForLorcaBootstrap(args)

	// create and launch the app
	ui, err := lorca.New("", "", 800, 600, args...)
	genericErrHandler(err, "initializing the app UI")
	defer ui.Close()

	// init the app model
	appPtr := initApp(ui)
	fmt.Println("app model:", appPtr)
	// remove it if not for demo
	initDemoApp(ui)

	// connect to FS (fileServer pointing to folder www)
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	genericErrHandler(err, "connecting to the fileServer (e.g. www folder)")
	defer listener.Close()

	// start the server by binding the listener
	go http.Serve(listener, http.FileServer(FS))

	// create the url for running
	url := fmt.Sprintf("http://%v", listener.Addr())
	if len(os.Args) > 1 {
		url = fmt.Sprintf("%v/%v", url, os.Args[1])
	}
	err = ui.Load(url)
	genericErrHandler(err, "load the index.html")

	// os signal handling
	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-ui.Done():
	}
	// can exit now
	fmt.Println("Thanks for using the app!")
}

func initApp(ui lorca.UI) (appPtr *app.App) {
	appPtr = app.NewApp()

	err := ui.Bind("onStart", func() {
		appPtr.OnStart()
	})
	genericErrHandler(err, "binding onStart event")

	// OnCreateNoteTask
	err = ui.Bind("onCreateNoteTask", func(content, todayInString string, x, y, angle string) {
		// pass also parameters from javascript side...?
		err2, notesInString := appPtr.OnCreateNoteTask(content, todayInString, x, y, angle)
		genericErrHandler(err2, "create note / task")
		// eval and emit a global event window.eventBus.$emit('xxx-event', object) => use ui.Eval()
		jsCommand := fmt.Sprintf(
			"window.eventBus.$emit('%v', JSON.parse('%v'));", eventNotesAllLoaded, notesInString)
		// PS. debug -> fmt.Println(jsCommand)
		//fmt.Println(jsCommand)

		ui.Eval(jsCommand)
	})
	genericErrHandler(err, "binding onStart event")


	err = ui.Bind( "onGetNotes", func() {
		notesInString := appPtr.GetNotesRepoInString()
		jsCommand := fmt.Sprintf(
			"window.eventBus.$emit('%v', JSON.parse('%v'));", eventNotesAllLoaded, notesInString)
		ui.Eval(jsCommand)
	})


	return
}

func initDemoApp(ui lorca.UI)  {
	err := ui.Bind("onStart", func() {
		// perform the server side task here
		fmt.Println("app started")
	})
	if err != nil {
		// or any exception handling code here
		panic(err)
	}

}
