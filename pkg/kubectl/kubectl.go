package kustomize

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/hashicorp/go-getter"
	log "github.com/sirupsen/logrus"
)

func Download(version string, path string) error {
	this_os, err := osDetect()
	if err != nil {
		return err
	}
	this_arch, err := archDetect()
	if err != nil {
		return err
	}
	client := &getter.Client{
		Src:  fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/%v/bin/%v/%v/kubectl", version, this_os, this_arch),
		Dst:  path,
		Mode: getter.ClientModeAny,
	}
	err = client.Get()
	if err != nil {
		return err
	}
	err = os.Chmod(filepath.Join(path, "kubectl"), 0755)
	if err != nil {
		return err
	}
	return nil
}

func osDetect() (string, error) {
	os := runtime.GOOS
	switch os {
	case "darwin", "linux":
		log.Debugf("%v is supported", os)
	default:
		errMsg := fmt.Sprintf("%v is not supported", os)
		log.Errorf(errMsg)
		return "", fmt.Errorf(errMsg)
	}
	return os, nil
}

func archDetect() (string, error) {
	arch := runtime.GOARCH
	switch arch {
	case "amd64":
		log.Debugf("%v is supported", arch)
	default:
		errMsg := fmt.Sprintf("%v is not supported", arch)
		log.Errorf(errMsg)
		return "", fmt.Errorf(errMsg)
	}
	return arch, nil
}
