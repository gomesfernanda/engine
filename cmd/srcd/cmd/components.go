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
	"context"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/src-d/engine-cli/components"
)

// componentsCmd represents the components command
var componentsCmd = &cobra.Command{
	Use:   "components",
	Short: "Manage source{d} components and their installations",
}

// componentsListCmd represents the components list command
var componentsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List source{d} components",
	Run: func(cmd *cobra.Command, args []string) {
		imgs, err := components.List(context.Background())
		if err != nil {
			log.Printf("could not list images: %v", err)
			os.Exit(1)
		}
		for _, img := range imgs {
			fmt.Println(img)
		}
	},
}

// componentsCmd represents the components install command
var componentsInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install source{d} component",
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			log.Printf("installing %s", arg)
			err := components.Install(context.Background(), arg)
			if err != nil {
				log.Printf("could not install %s: %v", arg, err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(componentsCmd)
	componentsCmd.AddCommand(componentsListCmd)
	componentsCmd.AddCommand(componentsInstallCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// componentsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// componentsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
