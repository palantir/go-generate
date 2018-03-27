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

package commoncmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/palantir/go-generate/gogenerate"
	"github.com/palantir/go-generate/gogenerate/config"
)

func NewRunCmd(use string, projectDirFlagVal, cfgFlagVal *string, verifyFlagVal *bool) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: "Run generators specified in configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectParam, err := loadConfig(*cfgFlagVal)
			if err != nil {
				return err
			}
			if *verifyFlagVal {
				if ok, err := gogenerate.Verify(*projectDirFlagVal, projectParam, cmd.OutOrStdout()); err != nil {
					return err
				} else if !ok {
					// if verification failed, return empty error -- the "Verify" call itself will have already written
					// the output to stdout and returning an empty error signals to handlers that no other output needs
					// to be printed.
					return fmt.Errorf("")
				}
				return nil
			}
			return gogenerate.Run(*projectDirFlagVal, projectParam, cmd.OutOrStdout())
		},
	}
}

func loadConfig(cfgFile string) (gogenerate.ProjectParam, error) {
	cfgYML, err := ioutil.ReadFile(cfgFile)
	if os.IsNotExist(err) {
		return gogenerate.ProjectParam{}, nil
	}
	if err != nil {
		return gogenerate.ProjectParam{}, errors.Wrapf(err, "failed to read file %s", cfgFile)
	}
	upgradedCfg, err := config.UpgradeConfig(cfgYML)
	if err != nil {
		return gogenerate.ProjectParam{}, err
	}

	var cfg config.ProjectConfig
	if err := yaml.Unmarshal(upgradedCfg, &cfg); err != nil {
		return gogenerate.ProjectParam{}, errors.Wrapf(err, "failed to unmarshal go-generate configuration")
	}
	return cfg.ToParam(), nil
}
