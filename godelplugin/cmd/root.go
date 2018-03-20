// Copyright 2016 Palantir Technologies, Inc.
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
	"github.com/palantir/godel/framework/pluginapi"
	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use:   "generate",
		Short: "Run generators specified in configuration",
	}

	projectDirFlagVal string
	cfgFlagVal        string
	verifyFlagVal     bool
)

func init() {
	pluginapi.AddProjectDirPFlagPtr(RootCmd.PersistentFlags(), &projectDirFlagVal)
	RootCmd.PersistentFlags().StringVar(&cfgFlagVal, "config", "", "the YAML configuration file for the generate task")
}
