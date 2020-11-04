package util

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

// FileExists check if file exists
func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// CreateDir creates a directory if it does not exist
func CreateDir(dirname string, permission os.FileMode) error {
	exists, err := FileExists(dirname)
	if err != nil {
		return err
	}

	if !exists {
		err := os.MkdirAll(dirname, permission)
		if err != nil {
			return err
		}
	}

	return nil
}

// CreateFileIfNotExists creates an empty file if it does not exist
func CreateFileIfNotExists(filename string, permission os.FileMode) error {
	exists, err := FileExists(filename)
	if err != nil {
		return err
	}
	if !exists {
		err = ioutil.WriteFile(filename, []byte{}, permission)
		if err != nil {
			return err
		}
	}
	return nil
}

// HomeDir returns the home directory for the current user.
// On Windows:
// 1. the first of %HOME%, %HOMEDRIVE%%HOMEPATH%, %USERPROFILE% containing a `.kube\config` file is returned.
// 2. if none of those locations contain a `.kube\config` file, the first of %HOME%, %USERPROFILE%, %HOMEDRIVE%%HOMEPATH% that exists and is writeable is returned.
// 3. if none of those locations are writeable, the first of %HOME%, %USERPROFILE%, %HOMEDRIVE%%HOMEPATH% that exists is returned.
// 4. if none of those locations exists, the first of %HOME%, %USERPROFILE%, %HOMEDRIVE%%HOMEPATH% that is set is returned.
func HomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOME")
		homeDriveHomePath := ""
		if homeDrive, homePath := os.Getenv("HOMEDRIVE"), os.Getenv("HOMEPATH"); len(homeDrive) > 0 && len(homePath) > 0 {
			homeDriveHomePath = homeDrive + homePath
		}
		userProfile := os.Getenv("USERPROFILE")

		// Return first of %HOME%, %HOMEDRIVE%/%HOMEPATH%, %USERPROFILE% that contains a `.kube\config` file.
		// %HOMEDRIVE%/%HOMEPATH% is preferred over %USERPROFILE% for backwards-compatibility.
		for _, p := range []string{home, homeDriveHomePath, userProfile} {
			if len(p) == 0 {
				continue
			}
			if _, err := os.Stat(filepath.Join(p, ".monoskope", "config")); err != nil {
				continue
			}
			return p
		}

		firstSetPath := ""
		firstExistingPath := ""

		// Prefer %USERPROFILE% over %HOMEDRIVE%/%HOMEPATH% for compatibility with other auth-writing tools
		for _, p := range []string{home, userProfile, homeDriveHomePath} {
			if len(p) == 0 {
				continue
			}
			if len(firstSetPath) == 0 {
				// remember the first path that is set
				firstSetPath = p
			}
			info, err := os.Stat(p)
			if err != nil {
				continue
			}
			if len(firstExistingPath) == 0 {
				// remember the first path that exists
				firstExistingPath = p
			}
			if info.IsDir() && info.Mode().Perm()&(1<<(uint(7))) != 0 {
				// return first path that is writeable
				return p
			}
		}

		// If none are writeable, return first location that exists
		if len(firstExistingPath) > 0 {
			return firstExistingPath
		}

		// If none exist, return first location that is set
		if len(firstSetPath) > 0 {
			return firstSetPath
		}

		// We've got nothing
		return ""
	}
	return os.Getenv("HOME")
}
