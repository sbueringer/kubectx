// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	"fmt"

	"github.com/spf13/cobra"
)

// contextCmd represents the context command
var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Gets and sets the current context",
	Long: `Gets and sets the current context

Examples:
  # Get current context
  kubectx context

  # Set current context
  kubectx context kubernetes-admin@kubernetes
	`,
	Run: func(cmd *cobra.Command, args []string) {
		loadKubeContext()

		switch len(args) {
		case 0:
			context, err := getCurrentContext()
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println(context)
			}
		case 1:
			setCurrentContext(args[0])
			saveKubeContext()
		default:
			cmd.Help()
		}
	},
}

func init() {
	rootCmd.AddCommand(contextCmd)
}
