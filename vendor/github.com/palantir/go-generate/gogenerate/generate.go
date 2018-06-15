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

package gogenerate

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/palantir/pkg/matcher"
	"github.com/pkg/errors"
)

// Run runs the generate task specified by the provided parameters. Returns an error if running the verification task
// fails.
func Run(rootDir string, projectParam ProjectParam, stdout io.Writer) error {
	_, err := runGenerate(rootDir, projectParam, stdout)
	return err
}

// Verify runs the generate task specified by the provided parameters and return true if the verification is successful
// (that is, running the generator did not change the declared outputs), false otherwise. If verification is not
// successful, the reason is written as output to the provided writer. Returns an error if an error is encountered when
// running the verify task itself.
func Verify(rootDir string, projectParam ProjectParam, stdout io.Writer) (bool, error) {
	diff, err := runGenerate(rootDir, projectParam, stdout)
	if err != nil {
		return false, err
	}

	if len(diff) == 0 {
		return true, nil
	}

	var sortedKeys []string
	for k := range diff {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	var outputParts []string
	outputParts = append(outputParts, fmt.Sprintf("Generators produced output that differed from what already exists: %v", sortedKeys))
	for _, k := range sortedKeys {
		outputParts = append(outputParts, fmt.Sprintf("  %s:", k))
		for _, currGenLine := range strings.Split(diff[k].String(), "\n") {
			outputParts = append(outputParts, fmt.Sprintf("    %s", currGenLine))
		}
	}
	fmt.Fprintln(stdout, strings.Join(outputParts, "\n"))
	return false, nil
}

func runGenerate(rootDir string, projectParam ProjectParam, stdout io.Writer) (map[string]ChecksumsDiff, error) {
	diffs := make(map[string]ChecksumsDiff)
	for _, k := range projectParam.Generators.SortedKeys() {
		v := projectParam.Generators[k]
		m := v.GenPaths
		origChecksums, err := checksumsForMatchingPaths(rootDir, m)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to compute checksums")
		}

		genDir := path.Join(rootDir, v.GoGenDir)
		cmd := exec.Command("go", "generate")
		cmd.Dir = genDir
		cmd.Stdout = stdout
		cmd.Stderr = stdout

		var envVars []string
		for k, v := range projectParam.Generators[k].Environment {
			envVars = append(envVars, fmt.Sprintf("%s=%v", k, v))
		}
		cmd.Env = append(envVars, os.Environ()...)

		if err := cmd.Run(); err != nil {
			return nil, errors.Wrapf(err, "failed to run go generate in %q", genDir)
		}

		newChecksums, err := checksumsForMatchingPaths(rootDir, m)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to compute checksums")
		}

		diff := origChecksums.compare(newChecksums)
		if len(diff) > 0 {
			diffs[k] = diff
		}
	}
	return diffs, nil
}

type checksumSet map[string]*fileChecksumInfo

func (c checksumSet) sortedKeys() []string {
	var sorted []string
	for k := range c {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)
	return sorted
}

type ChecksumsDiff map[string]string

func (c ChecksumsDiff) String() string {
	var sortedKeys []string
	for k := range c {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	var parts []string
	for _, k := range sortedKeys {
		parts = append(parts, fmt.Sprintf("%s: %s", k, c[k]))
	}
	return strings.Join(parts, "\n")
}

func (c checksumSet) compare(other checksumSet) ChecksumsDiff {
	diffs := make(map[string]string)

	// determine missing and extra entries
	for k := range c {
		if _, ok := other[k]; !ok {
			diffs[k] = "existed before, no longer exists"
		}
	}
	for k := range other {
		if _, ok := c[k]; !ok {
			diffs[k] = "did not exist before, now exists"
		}
	}

	// compare content
	for k, v := range c {
		otherV, ok := other[k]
		if !ok {
			continue
		}

		if v.isDir != otherV.isDir {
			if v.isDir {
				diffs[k] = "was previously a directory, is now a file"
			} else {
				diffs[k] = "was previously a file, is now a directory"
			}
			continue
		}
		if v.sha256checksum != otherV.sha256checksum {
			diffs[k] = fmt.Sprintf("previously had checksum %s, now has checksum %s", v.sha256checksum, otherV.sha256checksum)
		}
	}

	return diffs
}

type fileChecksumInfo struct {
	path           string
	isDir          bool
	sha256checksum string
}

func checksumsForMatchingPaths(rootDir string, m matcher.Matcher) (checksumSet, error) {
	pathsToChecksums := make(map[string]*fileChecksumInfo)
	if err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		relPath, err := filepath.Rel(rootDir, path)
		if err != nil {
			return err
		}
		if m.Match(relPath) {
			checksum, err := newChecksum(path, info)
			if err != nil {
				return err
			}
			pathsToChecksums[relPath] = checksum
		}
		return nil
	}); err != nil {
		return nil, errors.Wrapf(err, "failed to walk directory %q", rootDir)
	}
	return pathsToChecksums, nil
}

func newChecksum(filePath string, info os.FileInfo) (*fileChecksumInfo, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		// file is opened for reading only, so safe to ignore errors on close
		_ = f.Close()
	}()

	if info.IsDir() {
		return &fileChecksumInfo{
			path:  filePath,
			isDir: true,
		}, nil
	}

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}
	return &fileChecksumInfo{
		path:           filePath,
		sha256checksum: fmt.Sprintf("%x", h.Sum(nil)),
	}, nil
}
