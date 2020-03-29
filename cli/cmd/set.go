// Copyright © 2018 Synergic Mobility Project (https://synergic.mobi)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

type Server struct {
	ServerAddr string
}

var (
	serverstr = []string{}
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		server := []*Server{}
		for _, str := range serverstr {
			server = append(server, &Server{
				ServerAddr: str,
			})
		}
		bytes, err := json.Marshal(server)
		if err != nil {
			fmt.Println(err)
		} else {
			ioutil.WriteFile("setup.json", bytes, os.ModePerm)
			fmt.Printf("Setup Json")
		}
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
	setCmd.Flags().StringArrayVarP(&serverstr, "synerex", "s", []string{}, "server address")
}
