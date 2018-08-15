// Copyright 2018. box.la authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"

	"github.com/boxproject/agent/commands"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	commands.InitLogger()
	app := newApp()
	app.Run(os.Args)
}

func newApp() *cli.App {
	app := cli.NewApp()
	app.Version = PrintVersion(gitCommit, stage, version)
	app.Name = "Blockchain agent"
	app.Usage = "The blockchain monitor command line interface"
	app.Author = "2SE Group"
	app.Copyright = "Copyright 2017-2018 The exchange Authors"
	app.Email = "support@2se.com"
	app.Description = "blockchain agent"

	app.Commands = []cli.Command{
		// 启动
		{
			Name:   "start",
			Usage:  "start the monitor",
			Action: commands.StartCmd,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config,c",
					Usage: "Path of the config.json file",
					Value: "",
				},
				cli.StringFlag{
					Name:  "block-file,b",
					Usage: "Check point block number",
					Value: "",
				},
			},
		},
		// 停止
		{
			Name:   "stop",
			Usage:  "stop the monitor",
			Action: commands.StopCmd,
			Flags:  []cli.Flag{},
		},
	}

	return app
}

func PrintVersion(gitCommit, stage, version string) string {
	if gitCommit != "" {
		return fmt.Sprintf("%s-%s-%s", stage, version, gitCommit)
	}
	return fmt.Sprintf("%s-%s", stage, version)
}
