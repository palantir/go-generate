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
	"github.com/palantir/godel/pkg/versionedconfig"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/palantir/go-generate/gogenerate"
	"github.com/palantir/go-generate/internal/legacy"
)

func UpgradeConfig(cfgBytes []byte) ([]byte, error) {
	if legacy.IsLegacyConfig(cfgBytes) {
		return legacy.UpgradeLegacyConfig(cfgBytes)
	}

	version, err := versionedconfig.ConfigVersion(cfgBytes)
	if err != nil {
		return nil, err
	}
	switch version {
	case "", "0":
		var cfg gogenerate.ProjectConfig
		if err := yaml.UnmarshalStrict(cfgBytes, &cfg); err != nil {
			return nil, errors.Wrapf(err, "failed to unmarshal input as v0 YAML")
		}
		// input is valid current configuration: return exactly
		return cfgBytes, nil
	default:
		return nil, errors.Errorf("unsupported version: %s", version)
	}
}
