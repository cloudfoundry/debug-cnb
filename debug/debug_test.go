/*
 * Copyright 2018 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package debug_test

import (
	"testing"

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/debug-buildpack/debug"
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestDebug(t *testing.T) {
	spec.Run(t, "Debug", testDebug, spec.Report(report.Terminal{}))
}

func testDebug(t *testing.T, when spec.G, it spec.S) {

	it("returns true if build plan does exist", func() {
		f := test.NewBuildFactory(t)
		f.AddBuildPlan(t, debug.Dependency, buildplan.Dependency{})

		_, ok := debug.NewDebug(f.Build)
		if !ok {
			t.Errorf("NewDebug = %t, expected true", ok)
		}
	})

	it("returns false if build plan does not exist", func() {
		f := test.NewBuildFactory(t)

		_, ok := debug.NewDebug(f.Build)
		if ok {
			t.Errorf("NewDebug = %t, expected false", ok)
		}
	})

	it("contributes debug configuration", func() {
		f := test.NewBuildFactory(t)
		f.AddBuildPlan(t, debug.Dependency, buildplan.Dependency{})

		d, _ := debug.NewDebug(f.Build)
		if err := d.Contribute(); err != nil {
			t.Fatal(err)
		}

		layer := f.Build.Layers.Layer("debug")
		test.BeLayerLike(t, layer, false, true, true)
		test.BeProfileLike(t, layer, "debug", `PORT=${BPL_DEBUG_PORT:=8080}
SUSPEND=${BPL_DEBUG_SUSPEND:=n}

printf "Debugging enabled on port ${PORT}"

if [[ "${SUSPEND}" = "y" ]]; then
  printf ", suspended on start\n"
else
  printf "\n"
fi

export JAVA_OPTS="${JAVA_OPTS} -agentlib:jdwp=transport=dt_socket,server=y,address=${PORT},suspend=${SUSPEND}"
`)
	})
}
