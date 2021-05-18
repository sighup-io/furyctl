package vendoring

import (
	"os"
	"path/filepath"

	"github.com/sighupio/furyctl/internal/configuration"
	log "github.com/sirupsen/logrus"
)

func Download(path string, furyconf configuration.Furyconf, parallel bool) error {
	list, err := furyconf.Parse("", true)
	if err != nil {
		log.WithError(err).Error("Error parsing Furyfile")
		return err
	}
	if path == "" {
		path, err = os.Getwd()
		if err != nil {
			return err
		}
	}
	newList := make([]configuration.Package, len(list))
	for i, pkg := range list {
		pkg.Dir = filepath.Join(path, pkg.Dir)
		newList[i] = pkg
	}
	err = DownloadPackages(newList, parallel)
	if err != nil {
		log.WithError(err).Error("Error downloading dependencies on Furyfile")
		return err
	}
	return nil
}
