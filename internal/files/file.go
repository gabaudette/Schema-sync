package helpers

import (
	"os"
	"path/filepath"
	"ssync/internal/config"
	"strings"
)

func ReadDir(directories []string) ([]os.DirEntry, []string, error) {
	var files []os.DirEntry
	var paths []string

	conf := config.GetConfig()

	for _, dir := range directories {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			//log.Printf("Directory %s does not exist", dir)
			continue
		}

		// If the directory is in the ../api/ignored directories, skip it
		for _, ignoredDir := range conf.IgnoredDirs {
			//println("Ignored directory: ", ignoredDir)
			if strings.Contains(dir, ignoredDir) {
				continue
			}

			if filepath.Base(dir) == ignoredDir {
				continue
			}
		}

		err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {

			//println("Currently walking directory: ", path)
			if err != nil {
				return err
			}
			if !d.IsDir() {
				files = append(files, d)
				paths = append(paths, path)
			}
			return nil
		})
		if err != nil {
			//log.Printf("Error walking directory %s: %v", dir, err)
		}
	}

	return files, paths, nil
}
