// Copyright © 2019 The Things Network Foundation, The Things Industries B.V.
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

package ttnmage

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/TheThingsIndustries/magepkg/git"
	"github.com/blang/semver"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Version namespace.
type Version mg.Namespace

var goVersionFile = `// Code generated by Makefile. DO NOT EDIT.

package version

// TTN Version
var TTN = "%s-dev"
`

var currentVersion string

func (Version) getCurrent() error {
	_, _, tag, err := git.Info()
	if err != nil {
		return err
	}
	currentVersion = tag
	return nil
}

// Current returns the current version.
func (Version) Current() error {
	mg.Deps(Version.getCurrent)
	fmt.Println(currentVersion)
	return nil
}

// Files writes the current version to files that contain version info.
func (Version) Files() error {
	_, _, tag, err := git.Info()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("pkg/version/ttn.go", []byte(fmt.Sprintf(goVersionFile, tag)), 0644)
	if err != nil {
		return err
	}
	version := strings.TrimPrefix(tag, "v")
	for _, packageJSONFile := range []string{"package.json", "sdk/js/package.json"} {
		err = sh.Run(
			filepath.Join("node_modules", ".bin", "json"),
			"-f", packageJSONFile,
			"-I",
			"-e", fmt.Sprintf(`this.version="%s"`, version),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func bumpVersion(bump string) error {
	mg.Deps(Version.getCurrent)
	version, err := semver.Parse(strings.TrimPrefix(currentVersion, "v"))
	if err != nil {
		return err
	}
	var newVersion semver.Version
	switch bump {
	case "major":
		newVersion.Major = version.Major + 1
	case "minor":
		newVersion.Major = version.Major
		newVersion.Minor = version.Minor + 1
	case "patch":
		newVersion.Major = version.Major
		newVersion.Minor = version.Minor
		newVersion.Patch = version.Patch + 1
	}
	currentVersion = fmt.Sprintf("v%s", newVersion)
	return nil
}

// BumpMajor bumps a major version.
func (Version) BumpMajor() error { return bumpVersion("major") }

// BumpMinor bumps a minor version.
func (Version) BumpMinor() error { return bumpVersion("minor") }

// BumpPatch bumps a patch version.
func (Version) BumpPatch() error { return bumpVersion("patch") }
