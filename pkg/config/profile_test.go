package config

import "testing"

func TestResolvePrecedence(t *testing.T) {
	file := &File{
		Current: "default",
		Profiles: map[string]Profile{
			"default": {
				Token:      "file-token",
				APIURL:     "https://file.api",
				RestAPIURL: "https://file.api/rest",
				APIVersion: "2024-01-01",
				OrgID:      "file-org",
				GroupID:    "file-group",
				Output:     "yaml",
				PageSize:   50,
				Debug:      true,
			},
		},
	}

	env := map[string]string{
		"SNYK_API":        "https://env.api",
		"SNYK_OUTPUT":     "json",
		"SNYK_PAGE_SIZE":  "10",
		"SNYK_PROJECT_ID": "env-project",
	}

	overrideOutput := "table"
	flags := FlagOverrides{
		APIVersion: "2025-11-05",
		Output:     overrideOutput,
		ProjectID:  "flag-project",
	}

	res := Resolve(file, flags, env)

	if res.APIURL != "https://env.api" {
		t.Fatalf("expected APIURL from env, got %s", res.APIURL)
	}
	if res.RestAPIURL != "https://file.api/rest" {
		t.Fatalf("expected RestAPIURL from file, got %s", res.RestAPIURL)
	}
	if res.APIVersion != "2025-11-05" {
		t.Fatalf("expected APIVersion from flags, got %s", res.APIVersion)
	}
	if res.Output != "table" {
		t.Fatalf("expected Output from flags, got %s", res.Output)
	}
	if res.PageSize != 10 {
		t.Fatalf("expected PageSize from env, got %d", res.PageSize)
	}
	if !res.Debug {
		t.Fatalf("expected Debug to remain true from file")
	}
	if res.Token != "file-token" {
		t.Fatalf("expected Token from file, got %s", res.Token)
	}
	if res.ProjectID != "flag-project" {
		t.Fatalf("expected ProjectID from flags, got %s", res.ProjectID)
	}
	if res.ProfileName != "default" {
		t.Fatalf("expected ProfileName default, got %s", res.ProfileName)
	}
}

func TestResolveProfileSelection(t *testing.T) {
	file := &File{
		Current: "file",
		Profiles: map[string]Profile{
			"file": {
				APIURL: "https://file.api",
			},
			"env": {
				APIURL: "https://env.api",
			},
		},
	}

	env := map[string]string{
		"SNYK_API_PROFILE": "env",
	}

	res := Resolve(file, FlagOverrides{}, env)
	if res.ProfileName != "env" {
		t.Fatalf("expected env profile, got %s", res.ProfileName)
	}
	if res.APIURL != "https://env.api" {
		t.Fatalf("expected APIURL from env profile, got %s", res.APIURL)
	}

	// Flag overrides should beat env selection.
	flags := FlagOverrides{Profile: "file"}
	res2 := Resolve(file, flags, env)
	if res2.ProfileName != "file" {
		t.Fatalf("expected file profile via flags, got %s", res2.ProfileName)
	}
	if res2.APIURL != "https://file.api" {
		t.Fatalf("expected APIURL from file profile, got %s", res2.APIURL)
	}
}
