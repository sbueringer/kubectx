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
	"time"
	"strings"
	"regexp"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

// KubeContext comment
type KubeContext struct {
	GlobalEnv Environment            `json:"globalEnv"`
	Envs      map[string]Environment `json:"envs"`
}

// Environment comment
type Environment struct {
	UpdatedOn	  time.Time
	CurrentConfig string             `json:"currentConfig"`
	Contexts      map[string]Context `json:"contexts"`
}

// Context comment
type Context struct {
	CurrentContext string            `json:"currentContext"`
	Namespaces     map[string]string `json:"namespaces"`
}

var kubectxFile = os.Getenv("HOME") + "/.kube/kubectx.json"
var kubeContext KubeContext

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kubectx",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		loadKubeContext()

		binName := os.Args[0]

		var cmdName string
		var value string
		if binName == "kcfg" {
			cmdName = "config"
			if len(args) == 1 {
				value = args[0]
			}
		} else if binName == "kctx" {
			cmdName = "context"
			if len(args) == 1 {
				value = args[0]
			}
		} else if binName == "kns" {
			cmdName = "namespace"
			if len(args) == 1 {
				value = args[0]
			}
		} else {
			if len(args) == 1 {
				cmdName = args[0]
			}
			if len(args) == 2 {
				value = args[1]
			}
		}
		
		if cmdName == "" {
			fmt.Print("TODO print usage")
			return
		}		

		if cmdName != "" && value == "" {
			if cmdName == "config" {
				fmt.Println(getCurrentConfig())
			} else if cmdName == "context" {
				fmt.Println(getCurrentContext())
			} else if cmdName == "namespace" {
				fmt.Println(getCurrentNamespace())
			}
		} else if cmdName != "" && value != "" {
			if cmdName == "config" {
				setCurrentConfig(value)
				saveKubeContext()
			} else if cmdName == "context" {
				setCurrentContext(value)
				saveKubeContext()
			} else if cmdName == "namespace" {
				setCurrentNamespace(value)
				saveKubeContext()
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getCurrentNamespace() string {
	var currentNamespace string
	if terminalID := os.Getenv("TERMINAL_ID"); terminalID != "" {
		currentConfig := kubeContext.Envs[terminalID].CurrentConfig
		currentContext := kubeContext.Envs[terminalID].Contexts[currentConfig].CurrentContext
		currentNamespace = kubeContext.Envs[terminalID].Contexts[currentConfig].Namespaces[currentContext]
	}
	if currentNamespace != "" {
		return currentNamespace
	}
	currentConfig := kubeContext.GlobalEnv.CurrentConfig
	currentContext := kubeContext.GlobalEnv.Contexts[currentConfig].CurrentContext
	currentNamespace = kubeContext.GlobalEnv.Contexts[currentConfig].Namespaces[currentContext]
	if currentNamespace != "" {
		return currentNamespace
	}
	return "default"
}

func setCurrentNamespace(namespace string) {
	currentConfigName := getCurrentConfig()
	currentContextName := getCurrentContext()

	configContext, ok := kubeContext.GlobalEnv.Contexts[currentConfigName]
	if !ok {
		configContext = Context{
			Namespaces: make(map[string]string),
		}
	}
	if configContext.Namespaces == nil {
		configContext.Namespaces = make(map[string]string)
	}
	configContext.Namespaces[currentContextName] = namespace
	kubeContext.GlobalEnv.Contexts[currentConfigName] = configContext
	kubeContext.GlobalEnv.UpdatedOn = time.Now()

	if terminalID := os.Getenv("TERMINAL_ID"); terminalID != "" {
		env, ok := kubeContext.Envs[terminalID]
		if !ok {
			env = Environment{}
		}
		env.CurrentConfig = currentConfigName
		if env.Contexts == nil {
			env.Contexts = make(map[string]Context)
		}

		configContext, ok := env.Contexts[currentConfigName]
		if !ok {
			configContext = Context{
				Namespaces: make(map[string]string),
			}
		}
		configContext.CurrentContext = currentContextName

		if configContext.Namespaces == nil {
			configContext.Namespaces = make(map[string]string)
		}
		configContext.Namespaces[currentContextName] = namespace

		env.Contexts[currentConfigName] = configContext
		env.UpdatedOn = time.Now()
		kubeContext.Envs[terminalID] = env
	}
}

func getCurrentContext() string {
	var currentContext string
	if terminalID := os.Getenv("TERMINAL_ID"); terminalID != "" {
		currentConfig := kubeContext.Envs[terminalID].CurrentConfig
		currentContext = kubeContext.Envs[terminalID].Contexts[currentConfig].CurrentContext
	}
	if currentContext != "" {
		return currentContext
	}
	currentConfig := kubeContext.GlobalEnv.CurrentConfig
	currentContext = kubeContext.GlobalEnv.Contexts[currentConfig].CurrentContext
	if currentContext != "" {
		return currentContext
	}
	return getDefaultContext()
}

func getDefaultContext() string {
	currentConfig := getCurrentConfig()
	file, err := ioutil.ReadFile(currentConfig)
	if err != nil { panic(err) }

	currentContext := regexp.MustCompile(`current-context: .*`)
	match := currentContext.Find(file)

	return strings.Split(string(match), ": ")[1]
}

func setCurrentContext(context string) {
	currentConfigName := getCurrentConfig()

	if kubeContext.GlobalEnv.Contexts == nil {
		kubeContext.GlobalEnv.Contexts = make(map[string]Context)
	}

	configContext, ok := kubeContext.GlobalEnv.Contexts[currentConfigName]
	if !ok {
		configContext = Context{
			Namespaces: make(map[string]string),
		}
	}
	configContext.CurrentContext = context

	kubeContext.GlobalEnv.Contexts[currentConfigName] = configContext
	kubeContext.GlobalEnv.UpdatedOn = time.Now()

	if terminalID := os.Getenv("TERMINAL_ID"); terminalID != "" {
		env, ok := kubeContext.Envs[terminalID]
		if !ok {
			env = Environment{}
		}
		env.CurrentConfig = currentConfigName
		if env.Contexts == nil {
			env.Contexts = make(map[string]Context)
		}

		configContext, ok := env.Contexts[currentConfigName]
		if !ok {
			configContext = Context{
				Namespaces: make(map[string]string),
			}
		}
		configContext.CurrentContext = context

		env.Contexts[currentConfigName] = configContext
		env.UpdatedOn = time.Now()
		kubeContext.Envs[terminalID] = env
	}
}

func setCurrentConfig(config string) {
	kubeContext.GlobalEnv.CurrentConfig = config
	kubeContext.GlobalEnv.UpdatedOn = time.Now()

	if terminalID := os.Getenv("TERMINAL_ID"); terminalID != "" {
		env, ok := kubeContext.Envs[terminalID]
		if !ok {
			env = Environment{}
		}
		env.CurrentConfig = config
		env.UpdatedOn = time.Now()
		kubeContext.Envs[terminalID] = env
	}
}

func getCurrentConfig() string {
	var currentConfig string
	if terminalID := os.Getenv("TERMINAL_ID"); terminalID != "" {
		currentConfig = kubeContext.Envs[terminalID].CurrentConfig
	}
	if currentConfig != "" {
		return currentConfig
	}
	currentConfig = kubeContext.GlobalEnv.CurrentConfig
	if currentConfig != "" {
		return currentConfig
	}
	return os.Getenv("HOME") + "/.kube/config"
}

func loadKubeContext() {
	if _, err := os.Stat(kubectxFile); err == nil {
		file, err := ioutil.ReadFile(kubectxFile)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(file, &kubeContext)
		if err != nil {
			panic(err)
		}
	}
	if kubeContext.Envs == nil {
		kubeContext.Envs = make(map[string]Environment)
	}
}

func saveKubeContext() {

	yesterday := time.Now().Add(-24 * time.Hour)
	for envName, env := range kubeContext.Envs {
		if (env.UpdatedOn.Before(yesterday)){
			delete(kubeContext.Envs, envName)
		}
	}

	kubeContextJSON, err := json.Marshal(kubeContext)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(kubectxFile, kubeContextJSON, 0644)
	if err != nil {
		panic(err)
	}
}
