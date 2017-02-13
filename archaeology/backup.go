package archaeology

import (
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

func Backup(includes, ignores []string, dest string) {
	log.Info("includes: ", includes)
	log.Info("ignores: ", ignores)
	log.Info("destination: ", dest)

	// Walk over all paths and assemble a list of files to back up
	toBackup := []string{}

	addFile := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Print(err)
			return nil
		}

		if info.IsDir() {
			dir := filepath.Base(path)
			for _, ignore := range ignores {
				if ignore == dir {
					log.Info("Ignoring ", path, " (rule ", ignore, ")")
					return filepath.SkipDir
				}
			}
		}

		log.Info("Walking ", path)
		toBackup = append(toBackup, path)
		return nil
	}

	for _, path := range includes {
		err := filepath.Walk(path, addFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Info("Found ", len(toBackup), " files")

}

func backupPath(path string, ignores []string, dest string) {

}
