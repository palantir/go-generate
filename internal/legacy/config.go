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

package legacy

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/palantir/go-generate/gogenerate"
)

type legacyConfigStruct struct {
	Legacy                   bool `yaml:"legacy-config"`
	gogenerate.ProjectConfig `yaml:",inline"`
}

func IsLegacyConfig(cfgBytes []byte) bool {
	var cfg legacyConfigStruct
	if err := yaml.Unmarshal(cfgBytes, &cfg); err != nil {
		return false
	}
	return cfg.Legacy
}

func UpgradeLegacyConfig(cfgBytes []byte) ([]byte, error) {
	var legacyCfg legacyConfigStruct
	if err := yaml.UnmarshalStrict(cfgBytes, &legacyCfg); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal legacy configuration")
	}
	// succeed in unmarshalling legacy configuration. Legacy configuration is compatible with v0 configuration, so
	// simply return the provided input.
	return cfgBytes, nil
}
