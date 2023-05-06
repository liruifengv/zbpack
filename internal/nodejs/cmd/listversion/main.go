// listversion generates the current Node.js version statically.

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/Masterminds/semver"
	"github.com/samber/lo"
)

const TEMPLATE = `// Code generated by listversion; DO NOT EDIT.
//go:generate go run cmd/listversion/main.go

package nodejs

import (
	"github.com/Masterminds/semver/v3"
)

// nodeVersions is a list of all the Node.js versions.
// You should never change values here.
var nodeVersions = []*semver.Version{
	{{- range .}}
	semver.MustParse("{{.}}"),
	{{- end }}
}
`

// getNodeVersionsList fetches the major version of node from
// GitHub (https://github.com/nodejs/Release).
func getNodeVersionsList() ([]*semver.Version, error) {
	const releaseUrl = "https://raw.githubusercontent.com/nodejs/Release/master/schedule.json"

	request, err := http.NewRequest("GET", releaseUrl, nil)
	if err != nil || request == nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil || response == nil {
		return nil, err
	}
	defer response.Body.Close()

	/*  "v20": {
		"start": "2023-04-18",
		"lts": "2023-10-24",
		"maintenance": "2024-10-22",
		"end": "2026-04-30",
		"codename": ""
	}
	*/
	var versionsMap map[string]struct{}
	if err := json.NewDecoder(response.Body).Decode(&versionsMap); err != nil {
		return nil, err
	}

	// convert the unordered map to a slice of string
	versionsList := make([]*semver.Version, 0, len(versionsMap))
	for version := range versionsMap {
		// we ignore all the versions which is unstable (v0.x)
		if strings.HasPrefix(version, "v0.") {
			continue
		}

		// parse the version
		semverVersion, err := semver.NewVersion(version)
		if err != nil {
			return nil, err
		}

		versionsList = append(versionsList, semverVersion)
	}

	// sort the versions
	sort.Slice(versionsList, func(i, j int) bool {
		// we expect descending sort here, so we pass the
		// greater-than result to `less` function
		return versionsList[i].GreaterThan(versionsList[j])
	})

	return versionsList, nil
}

func main() {
	var fileHandle *os.File
	defer fileHandle.Close()

	file, ok := os.LookupEnv("GOFILE")
	if ok {
		_fileHandle, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC, 0o644)
		if err != nil {
			panic(err)
		}

		fileHandle = _fileHandle
	} else {
		fileHandle = os.Stdout // output to stdout
	}

	versions, err := getNodeVersionsList()
	if err != nil {
		panic(err)
	}

	// convert the version to string
	versionStr := lo.Map(versions, func(v *semver.Version, _ int) string {
		return v.String()
	})

	// generate the code
	tmpl, err := template.New("nodejs-version").Parse(TEMPLATE)
	if err != nil {
		panic(err)
	}

	if err := tmpl.Execute(fileHandle, versionStr); err != nil {
		panic(err)
	}

	log.Println("done!")
}