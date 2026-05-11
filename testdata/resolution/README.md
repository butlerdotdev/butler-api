# Resolution Test Fixtures

Shared test fixtures for User-Team membership resolution. Both butler-server
and butler-controller implement the resolution logic independently (per the
Factor-1 decision in the User-Team membership spec). These fixtures encode
the resolution contract so that drift between implementations is detected
by CI rather than discovered in production.

## How it works

Each JSON file describes one test case: a set of inputs (user email,
lastSeenGroups, Team CRDs) and the expected output (sorted list of
UserTeamMembership). Consumer tests import this package and assert their
resolution function produces output matching each fixture.

## Fixture format

```json
{
  "name": "descriptive_case_name",
  "description": "What this fixture tests",
  "input": {
    "email": "user@example.com",
    "lastSeenGroups": ["group-id-1"],
    "teams": [
      {
        "name": "team-name",
        "access": {
          "users": [{"name": "user@example.com", "role": "admin"}],
          "groups": [{"name": "group-id-1", "role": "viewer"}]
        }
      }
    ]
  },
  "expected": [
    {"name": "team-name", "role": "admin"}
  ]
}
```

## Contract assertions

- Expected output is sorted by team name ascending
- When manual and group access overlap on the same team, the highest role wins
- Empty role strings default to "viewer"
- Group matching is symmetric: both user-side and config-side groups are normalized
- Normalization: TrimSpace, LDAP CN extraction, email domain stripping, lowercase

## Adding fixtures

1. Create a new JSON file following the naming convention: `NN_descriptive_name.json`
2. Both butler-server and butler-controller tests consume all fixtures via `LoadAll()`
3. Adding a fixture requires both consumer test suites to pass against it
4. Removing or modifying an existing fixture is a breaking change for downstream tests

## Consuming fixtures

```go
import "github.com/butlerdotdev/butler-api/testdata/resolution"

func TestResolutionDrift(t *testing.T) {
    for _, fixture := range resolution.LoadAll() {
        t.Run(fixture.Name, func(t *testing.T) {
            // Map fixture inputs to your domain types, call your
            // resolution function, sort the output by team name,
            // compare against fixture.Expected.
        })
    }
}
```

## Consumer tests must sort output

The expected output in each fixture is sorted by team name ascending.
Consumer tests must sort their resolution function's output before
comparing against the fixture expectations. This ensures the fixtures
test resolution correctness, not implementation-specific iteration order.
