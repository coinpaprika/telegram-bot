// Copyright © 2018 coinpaprika.com
//
// Licensed under the Apache License, version 2.0 (the "License");
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
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows the application version and git revision",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("\nBuilt : %s \nVersion: %s (with : %s) \nDate: %s\n\n", commit, version, runtime.Version(), date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
