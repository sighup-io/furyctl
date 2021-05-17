package distribution

import (
	"path/filepath"

	"github.com/hashicorp/go-getter"
	log "github.com/sirupsen/logrus"
)

const (
	furyFile                    = "Furyfile.yml"
	kustomizationFile           = "kustomization.yaml"
	httpsDistributionRepoPrefix = "http::https://github.com/sighupio/fury-distribution/releases/download/"
)

var fileNames = [...]string{furyFile, kustomizationFile}

func Download(version, path string) error {
	for _, fileName := range fileNames {
		url := httpsDistributionRepoPrefix + version + "/" + fileName
		err := downloadFile(url, fileName, path)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

func downloadFile(url string, outputFileName string, path string) error {
	err := get(url, outputFileName, path, getter.ClientModeFile)
	if err != nil {
		log.Error(err)
	}
	return err
}

func get(src, dest string, path string, mode getter.ClientMode) error {
	log.Debugf("complete url downloading: %s -> %s", src, dest)
	client := &getter.Client{
		Src:  src,
		Dst:  filepath.Join(path, dest),
		Mode: mode,
	}
	log.Debugf("let's get %s -> %s", src, dest)
	err := client.Get()
	if err != nil {
		log.Error(err)
		return err
	}
	log.Debugf("done %s -> %s", src, dest)
	return err
}
