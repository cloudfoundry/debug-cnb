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

package debug_buildpack_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/buildpack/libbuildpack"
	"github.com/cloudfoundry/debug-buildpack"
	"github.com/cloudfoundry/libjavabuildpack"
	"github.com/cloudfoundry/libjavabuildpack/test"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestDebug(t *testing.T) {
	spec.Run(t, "Debug", testDebug, spec.Report(report.Terminal{}))
}

func testDebug(t *testing.T, when spec.G, it spec.S) {

	it("returns true if build plan does exist", func() {
		f := test.NewBuildFactory(t)
		f.AddBuildPlan(t, debug_buildpack.DebugDependency, libbuildpack.BuildPlanDependency{})

		_, ok := debug_buildpack.NewDebug(f.Build)
		if !ok {
			t.Errorf("NewDebug = %t, expected true", ok)
		}
	})

	it("returns false if build plan does not exist", func() {
		f := test.NewBuildFactory(t)

		_, ok := debug_buildpack.NewDebug(f.Build)
		if ok {
			t.Errorf("NewDebug = %t, expected false", ok)
		}
	})

	it("contributes debug configuration", func() {
		f := test.NewBuildFactory(t)
		f.AddBuildPlan(t, debug_buildpack.DebugDependency, libbuildpack.BuildPlanDependency{})

		d, _ := debug_buildpack.NewDebug(f.Build)
		if err := d.Contribute(); err != nil {
			t.Fatal(err)
		}

		layerRoot := filepath.Join(f.Build.Launch.Root, "debug")
		test.BeFileLike(t, filepath.Join(layerRoot, "profile.d", "debug"), 0644,
			`PORT=${BPL_DEBUG_PORT:=8080}
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

	it("reuses debug configuration", func() {
		f := test.NewBuildFactory(t)
		f.AddBuildPlan(t, debug_buildpack.DebugDependency, libbuildpack.BuildPlanDependency{})

		libjavabuildpack.WriteToFile(strings.NewReader("debug = true"), filepath.Join(f.Build.Launch.Root, "debug.toml"), 0644)

		d, _ := debug_buildpack.NewDebug(f.Build)
		if err := d.Contribute(); err != nil {
			t.Fatal(err)
		}


		if exist, err := libjavabuildpack.FileExists(filepath.Join(f.Build.Launch.Root, "debug", "profile.d", "debug")) ; err != nil {
			t.Fatal(err)
		} else if exist {
			t.Errorf("Contribute created profile.d/debug, expected not to")
		}
	})

}
