package lua

import (
	_ "embed"
	"errors"
	"os"
	"path/filepath"

	"github.com/kmaasrud/doctor/log"
	"github.com/kmaasrud/doctor/utils"
)

//go:embed crossref.lua
var crossref []byte

//go:embed abstract.lua
var abstract []byte

//go:embed wordcount.lua
var wordcount []byte

var buildFilters = map[string][]byte{
	"crossref.lua": crossref,
	"abstract.lua": abstract,
}

// Function returns a list of paths too all buildFilters. If the filters are not available in the
// Doctor data directory, then the function will ensure they are written there first.
func BuildFilters() []string {
    msg := log.Get()

	var paths []string

	dataDir, err := utils.FindDoctorDataDir()
	if err != nil {
		msg.Warning("Skipping Lua filters. Could not determine the Doctor data directory:\n    " + err.Error())
		return paths
	}

	embedDir := filepath.Join(dataDir, "embedded")

	err = utils.EnsureDir(embedDir)
	if err != nil {
		msg.Warning("Skipping Lua filters. Could not ensure Doctor data directory exists:\n    " + err.Error())
		return paths
	}

	for filename, filter := range buildFilters {
		path := filepath.Join(embedDir, filename)

		if _, err := os.Stat(path); os.IsNotExist(err) {
			msg.Info("Writing filter '" + filename + "' to '" + embedDir + "'.")
			err := os.WriteFile(path, filter, 0644)
			if err != nil {
				msg.Warning("Could not create '" + filename + "', skipping it. " + err.Error())
				continue
			}
		}

		paths = append(paths, path)
	}

	return paths
}

func WordCountFilter() (string, error) {
    msg := log.Get()

	dataDir, err := utils.FindDoctorDataDir()
	if err != nil {
		return "", errors.New("Could not determine the Doctor data directory.")
	}

	embedDir := filepath.Join(dataDir, "embedded")

	err = utils.EnsureDir(embedDir)
	if err != nil {
		return "", err
	}

	path := filepath.Join(embedDir, "wordcount.lua")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		msg.Info("Writing filter 'wordcount.lua' to '" + embedDir + "'.")
		err := os.WriteFile(path, wordcount, 0644)
		if err != nil {
			return "", errors.New("Could not create 'wordcount.lua'. " + err.Error())
		}
	}

	return path, nil
}
