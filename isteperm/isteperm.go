package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"go.mattglei.ch/timber"
)

const DIRECTORY = "/Users/matt/src/gleich/rit-cs-labs/iste-140"

func main() {
	timber.Timezone(time.Local)
	timber.TimeFormat("03:04:05")

	err := os.Chdir(DIRECTORY)
	if err != nil {
		timber.Fatal(err, "failed to change directory to", DIRECTORY)
	}

	updated := 0
	err = filepath.WalkDir(DIRECTORY, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("%w failed to walk directory at path %s", err, path)
		}

		info, err := d.Info()
		if err != nil {
			return fmt.Errorf("%w failed to get info for %s", err, path)
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
				return fmt.Errorf("%w failed to set file permissions of %d for %s", err, mode, path)
			}
			relativePath, err := filepath.Rel(DIRECTORY, path)
			if err != nil {
				return fmt.Errorf("%w failed to get relative path for %s", err, path)
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
