// Copyright The Cryostat Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controllers_test

import (
	"github.com/cryostatio/cryostat-operator/internal/controllers"
	. "github.com/onsi/ginkgo/v2"
)

// TODO get multi-namespace tests from ClusterCryostat controller, or move to reconciler_tests.go
// TODO add conversion webhook tests (maybe try enumerate all test.NewCryostat* methods and try converting them)
var _ = Describe("CryostatController", func() {
	c := &controllerTest{
		clusterScoped:   false,
		constructorFunc: newCryostatController,
	}

	c.commonTests()
})

func newCryostatController(config *controllers.ReconcilerConfig) (controllers.CommonReconciler, error) {
	return controllers.NewCryostatReconciler(config)
}
