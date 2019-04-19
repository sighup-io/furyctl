// Copyright © 2018 Sighup SRL support@sighup.io
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package furyconfig

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

const (
	configFileName          = "Furyfile"
	repoPrefix              = "git@github.com:sighup-io/fury-kubernetes"
	defaultVendorFolderName = "vendor"
)

// Furyconf is reponsible for the structure of the Furyfile
type Furyconf struct {
	VendorFolderName string    `yaml:"vendorFolderName"`
	Roles            []Package `yaml:"roles"`
	Modules          []Package `yaml:"modules"`
	Bases            []Package `yaml:"bases"`
}

// Package is the type to contain the definition of a single package
type Package struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	URL     string
	Dir     string
	Kind    string
}

// Validate is used for validation of configuration and initization of default paramethers
func (f *Furyconf) Validate() error {
	if f.VendorFolderName == "" {
		f.VendorFolderName = defaultVendorFolderName
	}
	return nil
}

// GetPackages parse the packages from FuryFile and return them in a list
func (f *Furyconf) GetPackages() []Package {
	list, err := f.Parse()
	if err != nil {
		log.Println("ERROR PARSING Furyfile: ", err)
	}

	return list
}

// Parse reads the furyconf structs and created a list of packaged to be downloaded
func (f *Furyconf) Parse() ([]Package, error) {
	pkgs := make([]Package, 0, 0)
	// First we aggreggate all packages in one single list
	for _, v := range f.Roles {
		v.Kind = "roles"
		pkgs = append(pkgs, v)
	}
	for _, v := range f.Modules {
		v.Kind = "modules"
		pkgs = append(pkgs, v)
	}
	for _, v := range f.Bases {
		v.Kind = "katalog"
		pkgs = append(pkgs, v)
	}

	// Now we generate the dowload url and local dir
	for i := 0; i < len(pkgs); i++ {
		block := strings.Split(pkgs[i].Name, "/")
		if len(block) == 2 {
			// the double slash is required to separate the repository name from a folder inside the repository
			// example: github.com/my-org/my-repo//foldername
			pkgs[i].URL = fmt.Sprintf("%s-%s//%s/%s?ref=%s", repoPrefix, block[0], pkgs[i].Kind, block[1], pkgs[i].Version)
		} else if len(block) == 1 {
			pkgs[i].URL = fmt.Sprintf("%s-%s//%s?ref=%s", repoPrefix, block[0], pkgs[i].Kind, pkgs[i].Version)
		}
		pkgs[i].Dir = fmt.Sprintf("%s/%s/%s", f.VendorFolderName, pkgs[i].Kind, pkgs[i].Name)
	}
	return pkgs, nil
}

// MustReadFuryFile load Furyfile from current path and validates it
func MustReadFuryFile(configFolderPath string) *Furyconf {
	viper.SetConfigType("yml")
	viper.AddConfigPath(configFolderPath)
	viper.SetConfigName(configFileName)

	c := new(Furyconf)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(c)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	err = c.Validate()
	if err != nil {
		log.Println("ERROR VALIDATING Furyfile: ", err)
	}

	return c
}
