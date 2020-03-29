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
	//"flag"
	//"encoding/json"
	//"flag"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/mtfelian/golang-socketio/transport"

	gosocketio "github.com/mtfelian/golang-socketio"

	"github.com/spf13/cobra"
)

var sioClients []*gosocketio.Client

var sioClient *gosocketio.Client
var Providers []string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sm",
	Short: "Synergic Exchange command launcher",
	Long: `Synergic Exchange command launcher
For example:

se is a CLI launcher for Synergic Exchange.

   se run all
   se status   // show status of provider/servers
`,
	//	Run: func(cmd *cobra.Command, args[]string){},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		//		fmt.Println(err)
		os.Exit(1)
	}
}

func getSynerexAddr() []string {
	bytes, _ := ioutil.ReadFile("setup.json")
	fmt.Printf("addr : %v\n", bytes)
	server := make([]Server, 0)
	if err := json.Unmarshal(bytes, &server); err != nil {
		panic(err)
	}
	fmt.Printf("addr : %v\n", server)
	synerexAddrs := []string{}
	for _, serv := range server {
		synerexAddrs = append(synerexAddrs, serv.ServerAddr)
		fmt.Printf("addr : %v\n", serv.ServerAddr)
	}
	fmt.Printf("addr2 : %v\n", synerexAddrs)
	return synerexAddrs
}

func init() {

	//cobra.OnInitialize(initConfig)

	//synerexAddrs := getSynerexAddr()
	serviceAddr := os.Getenv("SIMULATOR_SERVICE_NAME")
	synerexAddrs := []string{serviceAddr}

	//ch := make(chan bool, len(synerexAddrs))
	for _, addr := range synerexAddrs {
		ioAddr := "ws://" + addr + "/socket.io/?EIO=3&transport=websocket"
		log.Printf("ioAddr: %v", ioAddr)
		go func() {
			var err error
			sioClient, err := gosocketio.Dial(ioAddr, transport.DefaultWebsocketTransport())

			if err != nil {
				fmt.Println("se: Error to connect with se-daemon. You have to start se-daemon first. %v", err)
				return
			}
			sioClient.On(gosocketio.OnConnection, func(c *gosocketio.Channel, param interface{}) {
				//		fmt.Println("Go socket.io connected ",c)
			})
			sioClient.On("providers", func(c *gosocketio.Channel, param interface{}) {
				//fmt.Println("Get Providers ",param)
				// we have to keep this to check parameters
				procs := param.([]interface{})
				Providers = make([]string, len(procs))
				for i, pp := range procs {
					Providers[i] = pp.(string)
				}
			})

			sioClient.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel, param interface{}) {
				fmt.Println("Go socket.io disconnected ", c)
			})

			sioClients = append(sioClients, sioClient)

		}()
	}
	/*if len(synerexAddrs) != 0 {
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			for {
				<-ch
				if len(ch) == len(synerexAddrs) {
					wg.Done()
				}
			}
		}()
		wg.Wait()
	}*/
	time.Sleep(200 * time.Millisecond)
	fmt.Printf("sioClients %v", sioClients)
}
