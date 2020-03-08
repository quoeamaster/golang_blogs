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
	"runtime"
)

func genericErrHandler(err error, description ...string)  {
	if err != nil {
		if description != nil {
			fmt.Println(fmt.Sprintf("oops! something is wrong! %v\n", description[0]))
		}
		panic(err)
	}
}

/**
 *	prepare bootstrap arguments for different OS (for the moment, only Linux)
 */
func prepareArgsForLorcaBootstrap(args []string) []string {
	if runtime.GOOS == "linux" {
		args = append(args, "--class=Lorca")
	}
	return args
}


