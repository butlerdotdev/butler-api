/*
Copyright 2026 The Butler Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package resolution provides shared test fixtures for the User-Team
// membership resolution contract. Both butler-server and butler-controller
// implement the resolution logic independently (Factor-1 decision). These
// fixtures encode the expected behavior: same inputs must produce identical
// outputs from both implementations.
//
// Consumer tests import this package, call LoadAll(), and assert their
// resolution function produces output matching each fixture's Expected
// slice. Drift between implementations surfaces as a test failure in
// whichever component diverged.
//
// Expected output is sorted by team name ascending. Consumer tests must
// sort their resolution output before comparing against fixtures.
package resolution

import (
	"embed"
	"encoding/json"
	"path/filepath"
	"sort"
)

//go:embed *.json
var fixtureFS embed.FS

// FixtureTeamUser represents a user entry in a Team's access configuration.
type FixtureTeamUser struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

// FixtureTeamGroup represents a group entry in a Team's access configuration.
type FixtureTeamGroup struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

// FixtureTeamAccess represents the access configuration of a Team.
type FixtureTeamAccess struct {
	Users  []FixtureTeamUser  `json:"users"`
	Groups []FixtureTeamGroup `json:"groups"`
}

// FixtureTeam represents a Team CRD with only the fields relevant to
// resolution: name and access configuration.
type FixtureTeam struct {
	Name   string            `json:"name"`
	Access FixtureTeamAccess `json:"access"`
}

// FixtureInput contains the inputs to the resolution function.
type FixtureInput struct {
	Email          string        `json:"email"`
	LastSeenGroups []string      `json:"lastSeenGroups"`
	Teams          []FixtureTeam `json:"teams"`
}

// ExpectedMembership represents a single team membership in the expected
// output. Matches the shape of UserTeamMembership in the API types.
type ExpectedMembership struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

// Fixture represents a single test case with inputs and expected output.
type Fixture struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Input       FixtureInput         `json:"input"`
	Expected    []ExpectedMembership `json:"expected"`
}

// LoadAll reads all fixture JSON files from the embedded testdata and
// returns them sorted by fixture name. Panics on malformed fixtures
// (intended for test-time use, not runtime).
func LoadAll() []Fixture {
	entries, err := fixtureFS.ReadDir(".")
	if err != nil {
		panic("resolution fixtures: reading directory: " + err.Error())
	}

	var fixtures []Fixture
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		data, err := fixtureFS.ReadFile(entry.Name())
		if err != nil {
			panic("resolution fixtures: reading " + entry.Name() + ": " + err.Error())
		}

		var f Fixture
		if err := json.Unmarshal(data, &f); err != nil {
			panic("resolution fixtures: parsing " + entry.Name() + ": " + err.Error())
		}

		if f.Name == "" {
			panic("resolution fixtures: " + entry.Name() + " has empty name")
		}

		fixtures = append(fixtures, f)
	}

	sort.Slice(fixtures, func(i, j int) bool {
		return fixtures[i].Name < fixtures[j].Name
	})

	return fixtures
}
