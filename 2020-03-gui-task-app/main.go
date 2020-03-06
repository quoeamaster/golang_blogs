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
	"github.com/zserge/lorca"
	"net"
	"net/http"
	"os"
	"os/signal"
)

func main()  {
	args := []string{}
	prepareArgsForLorcaBootstrap(args)

	// create and launch the app
	ui, err := lorca.New("", "", 480, 320, args...)
	genericErrHandler(err, "initializing the app UI")
	defer ui.Close()

	// connect to FS (fileServer pointing to folder www)
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	genericErrHandler(err, "connecting to the fileServer (e.g. www folder)")
	defer listener.Close()

	// start the server by binding the listener
	go http.Serve(listener, http.FileServer(FS))

	err = ui.Load(fmt.Sprintf("http://%s", listener.Addr()))
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


