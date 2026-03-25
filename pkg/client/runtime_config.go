package client

import (
	"os"
	"strings"

	cfgpkg "github.com/sam1el/snyk-api/pkg/config"
)

func resolveRuntimeConfig() cfgpkg.Resolved {
	res := cfgpkg.Resolved{
		APIURL:     "https://api.snyk.io",
		RestAPIURL: "https://api.snyk.io/rest",
		APIVersion: DefaultAPIVersion.String(),
		Output:     "json",
		PageSize:   100,
	}

	path, err := cfgpkg.DefaultPath()
	if err != nil {
		return res
	}
	fileCfg, err := cfgpkg.LoadFile(path)
	if err != nil {
		return res
	}

	env := map[string]string{}
	for _, kv := range os.Environ() {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) == 2 {
			env[parts[0]] = parts[1]
		}
	}

	return cfgpkg.Resolve(fileCfg, cfgpkg.FlagOverrides{}, env)
}
