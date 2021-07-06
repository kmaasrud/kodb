package utils

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/kmaasrud/doctor/msg"
)


// Searches up the directory tree to find a doctor.yaml file and returns the path
// of the directory containing it. If it reaches the root directory without finding
// anything, the function returns an error.
func FindDoctorRoot() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		msg.Error(err.Error())
	}

	for {
		if filepath.Dir(path) == path {
			return "", errors.New("Could not find a Doctor document.")
		} else if _, err := os.Stat(filepath.Join(path, "doctor.toml")); os.IsNotExist(err) {
			path = filepath.Dir(path)
		} else {
			return path, nil
		}
	}
}


// Wrapper function around exec.LookPath. Also consults the Doctor data dir
func CheckPath(program string) (string, error) {
	// Check for program in Doctor's data directory
	doctorPath, err := FindDoctorDataDir()
	if err == nil {
		path := filepath.Join(doctorPath, "bin", program)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	// Check if program in PATH
	path, err := exec.LookPath(program)
	if err != nil {
		return "", errors.New("Could not find " + program + " in your PATH.")
	}
	return path, nil
}

// Ensures a directory exists. If not, creates it (and any needed parents.)
func EnsureDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
