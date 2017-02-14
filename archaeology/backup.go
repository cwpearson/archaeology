package archaeology

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	glob "github.com/mattn/go-zglob"
)

const (
	backup = iota
	delete
)

// Action models an action that the application will take
type Action struct {
	src   string
	op    int
	store Store
}

// Do executes the action
func (a *Action) Do() error {
	return errors.New("Action.Do Unimplemented")
}

func (a *Action) String() string {
	var opStr string
	switch a.op {
	case backup:
		opStr = "backup"
	case delete:
		opStr = "delete"
	default:
		log.Fatal("Unexpected op")
	}
	return "[" + opStr + "] " + a.src
}

func scan(includes, ignores []string) ([]string, error) {
	matchedPaths := []string{}

	addFile := func(path string, info os.FileInfo, err error) error {
		if os.IsPermission(err) {
			if info.IsDir() {
				log.Warn("[perm] Skipping dir ", path)
				return filepath.SkipDir
			}
			log.Warn("[perm] Skipping file ", path)
			return nil
		}

		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			dir := filepath.Base(path)
			for _, ignore := range ignores {
				matched, err := glob.Match(ignore, dir)
				if err != nil {
					return err
				}
				if matched {
					log.Info("Ignoring ", dir, " (rule ", ignore, ")")
					return filepath.SkipDir
				}
			}
		} else { // skip files
			for _, ignore := range ignores {
				matched, err := glob.Match(ignore, path)
				if err != nil {
					return err
				}
				if matched {
					log.Info("Ignoring ", path, " (rule ", ignore, ")")
					return nil
				}
			}
			log.Info("Adding ", path)
			matchedPaths = append(matchedPaths, path)

		}
		return nil
	}

	for _, path := range includes {
		err := filepath.Walk(path, addFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	return matchedPaths, nil
}

func sameContents(a, b io.Reader) bool {

	// Compare contents
	blockSize := 4096
	aBuf := make([]byte, 0, blockSize)
	bBuf := make([]byte, 0, blockSize)

	for {
		aBuf = aBuf[:cap(aBuf)]
		bBuf = bBuf[:cap(bBuf)]
		n, err := a.Read(aBuf)
		aBuf = aBuf[:n]
		if 0 == n {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		n, err = b.Read(bBuf)
		bBuf = bBuf[:n]
		if 0 == n {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		// Compare buffers
		if len(aBuf) != len(bBuf) {
			return false
		}
		if bytes.Equal(aBuf, bBuf) == false {
			return false
		}

		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

	}
	return true
}

func Backup(includes, ignores []string, dest string) error {
	log.Info("includes: ", includes)
	log.Info("ignores: ", ignores)
	log.Info("destination: ", dest)

	// Walk over all paths and assemble a list of files eligible for backup
	toBackup, err := scan(includes, ignores)
	if err != nil {
		return err
	}

	log.Info("Found ", len(toBackup), " files")

	// For each of those files, determine if it needs to be backed up
	store := &LocalStore{root: dest}
	actions := []Action{}
	for _, path := range toBackup {
		storeReader, err := store.MostRecent(path)
		if err != nil {
			return err
		}
		currentReader, err := os.Open(path)
		if err != nil {
			return err
		}
		if false == sameContents(storeReader, currentReader) {
			actions = append(actions, Action{path, backup, store})
		}
	}

	// Execute the assembled list of actions, and record any problems
	for _, action := range actions {
		err := action.Do()
		if err != nil {
			return err
		}
	}

	return nil
}

func backupPath(path string, ignores []string, dest string) {

}
