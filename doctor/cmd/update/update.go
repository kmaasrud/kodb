package update

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kmaasrud/doctor/msg"
	"github.com/kmaasrud/doctor/utils"
	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

func Update(ver string) error {
	if ver == "DEV" {
		fmt.Println("You are on the DEV version of Doctor (downloaded with the Go tool.)\nRun 'go install github.com/kmaasrud/doctor@latest' to update.")
		return nil
	}

	latest, found, err := selfupdate.DetectLatest("kmaasrud/doctor")
	if err != nil {
		return errors.New("Error occurred while detecting version: " + err.Error())
	}

	v := semver.MustParse(ver[1:])
	if !found || latest.Version.LTE(v) {
		msg.Success("Current version is the latest!")
		return nil
	}

	exe, err := os.Executable()
	if err != nil {
		return errors.New(`Could not locate executable path.
There might be an issue with the permissions of you Doctor binary, or you might have symlinked Doctor.`)
	}

	var confirmString string
	fmt.Printf("Do you want to update to version %s? (Y/n) ", msg.Style(latest.Version.String(), "Bold"))
	fmt.Scanln(&confirmString)
	if strings.ToLower(confirmString) == "n" {
		msg.Info("Keeping current version.")
		return nil
	}

	if err := selfupdate.UpdateTo(latest.AssetURL, exe); err != nil {
		return errors.New("Error occurred while updating binary: " + err.Error())
	}

    // Remove the cached embedded files, so they can be rewritten with any changes in the update.
    dataDir, err := utils.FindDoctorDataDir()
    if err == nil {
        err = os.RemoveAll(filepath.Join(dataDir, "embedded"))
        if err != nil {
            msg.Warning("Tried to remove cached help files, but encountered an error: " + err.Error())
        }
    }

	msg.Success("Successfully updated to version " + latest.Version.String() + "!")
	return nil
}
