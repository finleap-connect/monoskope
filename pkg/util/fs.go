package util

import (
	"io/ioutil"
	"os"
)

// FileExists check if the directory exists
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

// CreateDir creates a directory if it does no exist
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
