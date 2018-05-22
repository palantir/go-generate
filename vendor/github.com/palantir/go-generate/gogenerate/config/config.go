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

package config

import (
	"github.com/palantir/go-generate/gogenerate"
	"github.com/palantir/go-generate/gogenerate/config/internal/v0"
)

type ProjectConfig v0.ProjectConfig

func (cfg *ProjectConfig) ToParam() gogenerate.ProjectParam {
	generators := make(gogenerate.Generators)
	for k, v := range cfg.Generators {
		v := GeneratorConfig(v)
		generators[k] = v.ToParam()
	}
	return gogenerate.ProjectParam{
		Generators: generators,
	}
}

type GeneratorConfig v0.GeneratorConfig

func (cfg *GeneratorConfig) ToParam() gogenerate.GeneratorParam {
	return gogenerate.GeneratorParam{
		GoGenDir:    cfg.GoGenDir,
		GenPaths:    cfg.GenPaths.Matcher(),
		Environment: cfg.Environment,
	}
}
