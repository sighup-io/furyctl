package provisioners

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sighupio/furyctl/internal/configuration"
	"github.com/sighupio/furyctl/pkg/distribution"
	kube "github.com/sighupio/furyctl/pkg/kubectl"
	kust "github.com/sighupio/furyctl/pkg/kustomize"
	"github.com/sighupio/furyctl/pkg/vendoring"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type KubeProvisionOptions struct {
	KubectlVersion   string
	KustomizeVersion string
	FuryVersion      string
}

type KubeProvision struct {
	KubectlVersion    string
	KustomizeVersion  string
	FuryVersion       string
	BinDirectory      string
	ManifestDirectory string

	kubectlPath      string
	kustomizePath    string
	distributionPath string
}

func (k *KubeProvision) Init() error {
	kustomizePath, err := kust.Download(k.KustomizeVersion, k.BinDirectory)
	if err != nil {
		log.WithError(err).Errorf("Can not init Kustomize")
		return err
	}
	k.kustomizePath = kustomizePath
	kubectlPath, err := kube.Download(k.KubectlVersion, k.BinDirectory)
	if err != nil {
		log.WithError(err).Errorf("Can not init Kubectl")
		return err
	}
	k.kubectlPath = kubectlPath
	k.distributionPath = filepath.Join(k.ManifestDirectory, "distribution")
	err = os.MkdirAll(k.distributionPath, 0755)
	if err != nil {
		log.WithError(err).Error("Can not create distribution directory")
	}
	err = distribution.Download(k.FuryVersion, k.distributionPath)
	if err != nil {
		log.WithError(err).Error("Can not init Fury")
		return err
	}
	err = k.createBaseKustomizationFile()
	if err != nil {
		log.WithError(err).Error("can not create base kustomization.yaml file")
		return err
	}
	err = k.vendor()
	if err != nil {
		log.WithError(err).Error("can not vendor dependencies")
		return err
	}
	return nil
}

func (k *KubeProvision) createBaseKustomizationFile() error {
	baseKustomizationFilePath := filepath.Join(k.ManifestDirectory, "kustomization.yaml")
	f, err := os.Create(baseKustomizationFilePath)
	if err != nil {
		return err
	}
	defer f.Close()
	baseKustomization := `---
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
  - distribution
`
	_, err = f.WriteString(baseKustomization)
	if err != nil {
		return err
	}
	return nil
}

func (k *KubeProvision) vendor() error {
	conf := configuration.Furyconf{}
	furyFilePath := filepath.Join(k.distributionPath, "Furyfile.yml")
	furyFileContent, err := ioutil.ReadFile(furyFilePath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(furyFileContent, &conf)
	if err != nil {
		return err
	}
	err = conf.Validate()
	if err != nil {
		return err
	}
	err = vendoring.Download(k.distributionPath, conf, true)
	if err != nil {
		return err
	}
	return nil
}

func (k *KubeProvision) Build() error {
	cmd := exec.Command(k.kustomizePath, "build", ".")
	cmd.Dir = k.ManifestDirectory
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		log.WithError(err).Fatal("error executing kustomize build")
		return err
	}
	in := bufio.NewScanner(stdout)
	for in.Scan() {
		log.Debug(in.Text())
	}
	if in.Err() != nil {
		log.WithError(err).Error("output error")
		return in.Err()
	}
	return nil
}

func (k *KubeProvision) Deploy(kubeconfigPath string) error {

	cmd := exec.Command(k.kustomizePath, "build", ".")
	cmd.Dir = k.ManifestDirectory
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		log.WithError(err).Fatal("error executing kustomize build")
		return err
	}
	in := bufio.NewScanner(stdout)

	// TEMPORAL FILE TO USE WITH KUBECTL
	file, err := ioutil.TempFile(k.ManifestDirectory, "rendered-")
	if err != nil {
		log.WithError(err).Fatal(err)
		return err
	}
	defer os.Remove(file.Name())

	for in.Scan() {
		d := in.Text()
		log.Debug(d)
		file.WriteString(fmt.Sprintln(d))
	}
	if in.Err() != nil {
		log.WithError(err).Error("output error")
		return in.Err()
	}
	// kubectl
	cmd = exec.Command(k.kubectlPath, "--kubeconfig", kubeconfigPath, "apply", "-f", file.Name())
	cmd.Dir = k.ManifestDirectory
	stdout, err = cmd.StdoutPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		log.WithError(err).Fatal("error executing kubectl apply")
		return err
	}
	in = bufio.NewScanner(stdout)

	for in.Scan() {
		log.Debug(in.Text())
	}
	if in.Err() != nil {
		log.WithError(err).Error("output error")
		return in.Err()
	}
	return nil
}
