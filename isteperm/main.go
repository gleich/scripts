package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"pkg.mattglei.ch/timber"
)

const DIRECTORY = "/Users/matt/src/gleich/rit-cs-labs/iste-140"

func main() {
	timber.SetTimezone(time.Local)
	timber.SetTimeFormat("03:04:05")

	wd, err := os.Getwd()
	if err != nil {
		timber.Fatal(err, "failed to get current directory")
	}
	if wd != DIRECTORY {
		timber.FatalMsg("please run from", DIRECTORY)
	}

	updated := 0
	err = filepath.WalkDir(DIRECTORY, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("%v failed to walk directory at path %s", err, path)
		}

		info, err := d.Info()
		if err != nil {
			return fmt.Errorf("%v failed to get info for %s", err, path)
		}

		perm := info.Mode().Perm()
		var mode fs.FileMode
		if d.IsDir() && perm != 0750 {
			mode = 0750
		} else if !d.IsDir() && perm != 0650 {
			mode = 0650
		}
		if mode != 0 {
			err = os.Chmod(path, mode)
			if err != nil {
				return fmt.Errorf("%v failed to set file permissions of %d for %s", err, mode, path)
			}
			relativePath, err := filepath.Rel(DIRECTORY, path)
			if err != nil {
				return fmt.Errorf("%v failed to get relative path for %s", err, path)
			}
			timber.Done(relativePath, "set to", mode)
			updated++
		}

		return nil
	})
	if err != nil {
		timber.Fatal(err, "failed to walk directory")
	}
	timber.Done("updated permissions for", updated, "files/folders")
}
