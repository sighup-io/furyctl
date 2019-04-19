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
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	gp := os.Getenv("GOPATH")
	ap := filepath.Join(gp, "src/github.com/sighup-io/furyctl", "fixtures")
	c := MustReadFuryFile(ap)
	packages, err := c.Parse()

	got := packages
	want := []Package{
		Package{Name: "kube-node-common", Version: "master", URL: "git@github.com:sighup-io/fury-kubernetes-kube-node-common//roles?ref=master", Dir: "vendor/roles/kube-node-common", Kind: "roles"},
		Package{Name: "aws-ark", Version: "master", URL: "git@github.com:sighup-io/fury-kubernetes-aws-ark//modules?ref=master", Dir: "vendor/modules/aws-ark", Kind: "modules"},
		Package{Name: "monitoring/prometheus-operator", Version: "master", URL: "git@github.com:sighup-io/fury-kubernetes-monitoring//katalog/prometheus-operator?ref=master", Dir: "vendor/katalog/monitoring/prometheus-operator", Kind: "katalog"},
	}

	if err != nil {
		panic(err)
	}

	if !reflect.DeepEqual(packages, want) {
		t.Errorf("got %v want %v\n", got, want)
	}

}
