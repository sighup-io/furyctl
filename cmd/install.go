// Copyright © 2018 Sighup SRL support@sighup.io
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
	"log"

	"github.com/sighup-io/furyctl/internal/furyconfig"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.PersistentFlags().BoolVarP(&parallel, "parallel", "p", true, "if true enables parallel downloads")
}

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Download dependencies specified in Furyfile.yml",
	Long:  "Download dependencies specified in Furyfile.yml",
	Run: func(cmd *cobra.Command, args []string) {
		c := furyconfig.MustReadFuryFile("./")
		packages := c.GetPackages()
		err := download(packages)
		if err != nil {
			log.Println("ERROR DOWNLOADING: ", err)
		}
	},
}
