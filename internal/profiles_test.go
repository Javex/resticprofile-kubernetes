package internal_test

import (
	"resticprofilek8s/internal"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestProfilesFromYaml(t *testing.T) {
	tcs := []struct {
		name          string
		profileString string
		expProfile    map[string]internal.Profile
	}{
		{
			"defaultStop",
			`
        test:
          name: testName
          namespace: testNamespace
          deployment: testDeployment
          host: test.example.com
          Folders: 
            - /test/folder
      `,

			map[string]internal.Profile{
				"test": {
					Name:       "testName",
					Namespace:  "testNamespace",
					Deployment: "testDeployment",
					Stop:       false,
					Host:       "test.example.com",
					Folders:    []string{"/test/folder"},
				},
			},
		},
		{
			"allFields",
			`
        test:
          name: testName
          namespace: testNamespace
          deployment: testDeployment
          host: test.example.com
          stop: true
          Folders: 
            - /test/folder
      `,

			map[string]internal.Profile{
				"test": {
					Name:       "testName",
					Namespace:  "testNamespace",
					Deployment: "testDeployment",
					Stop:       true,
					Host:       "test.example.com",
					Folders:    []string{"/test/folder"},
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			profiles, err := internal.ProfilesFromYamlString(tc.profileString)
			if err != nil {
				t.Fatalf("Error parsing yaml: %v", err)
			}

			if !cmp.Equal(profiles, tc.expProfile) {
				t.Fatalf("Unexpected difference in profiles: %v", cmp.Diff(profiles, tc.expProfile))
			}
		})
	}
}

func TestProfilesFromYamlErrors(t *testing.T) {
	tcs := []struct {
		name          string
		profileString string
		// expErr        error
	}{
		{
			name: "emptyFolders",
			profileString: `
        test:
          name: testName
          namespace: testNamespace
          deployment: testDeployment
          host: test.example.com
      `,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			_, err := internal.ProfilesFromYamlString(tc.profileString)
			if err == nil {
				t.Fatalf("Expected an error, but got nil")
			}
		})
	}
}

func TestProfileToTOML(t *testing.T) {
	tcs := []struct {
		name    string
		profile internal.Profile
		repo    internal.ResticRepository
		expTOML string
	}{
		{
			name: "default",
			profile: internal.Profile{
				Name:       "testName",
				Namespace:  "testNamespace",
				Deployment: "testDeployment",
				Host:       "test.example.com",
				Folders:    []string{"/test/folder", "/test/folder2"},
			},
			repo: internal.ResticRepository{
				Name:   "testRepo",
				Suffix: "-testSuffix",
			},
			expTOML: `
[testName-testSuffix]
  inherit = "testRepo"
  [testName-testSuffix.backup]
    tag = ["testName"]
    source = [
    "/test/folder",
    "/test/folder2",
    ]
    host = "test.example.com"
  [testName-testSuffix.snapshots]
    tag = ["testName"]
    host = "test.example.com"
`,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			toml, err := tc.profile.ToTOML(tc.repo)
			if err != nil {
				t.Fatalf("Error generating TOML: %v", err)
			}

			if toml != tc.expTOML {
				t.Fatalf("Unexpected difference in TOML: %v", cmp.Diff(toml, tc.expTOML))
			}
		})
	}
}
