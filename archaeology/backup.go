package archaeology

import log "github.com/Sirupsen/logrus"

func Backup(includes, ignores []string, dest string) {
	log.Info("includes: ", includes)
	log.Info("ignores: ", ignores)
	log.Info("destination: ", dest)
}
