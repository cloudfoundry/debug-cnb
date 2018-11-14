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

package debug_buildpack

import (
	"fmt"
	"reflect"

	"github.com/cloudfoundry/libjavabuildpack"
	"github.com/fatih/color"
)

// Debug indicates that a JVM application should be run with debugging enabled.
const DebugDependency = "debug"

// Debug represents the debug configuration for a JVM application
type Debug struct {
	layer  libjavabuildpack.LaunchLayer
	logger libjavabuildpack.Logger
}

// Contribute makes the contribution to launch
func (d Debug) Contribute() error {
	expected := marker{true}

	var m marker
	if err := d.layer.ReadMetadata(&m); err != nil {
		d.logger.Debug("Marker is not structured correctly")
		return err
	}

	if reflect.DeepEqual(expected, m) {
		d.logger.FirstLine("%s cached launch layer", color.GreenString("Reusing"))
		return nil
	}

	d.logger.Debug("Marker %s does not match expected %s", m, expected)

	d.logger.FirstLine("%s to launch", color.YellowString("Contributing"))

	d.layer.WriteProfile("debug", `PORT=${BPL_DEBUG_PORT:=8080}
SUSPEND=${BPL_DEBUG_SUSPEND:=n}

printf "Debugging enabled on port ${PORT}"

if [[ "${SUSPEND}" = "y" ]]; then
  printf ", suspended on start\n"
else
  printf "\n"
fi

export JAVA_OPTS="${JAVA_OPTS} -agentlib:jdwp=transport=dt_socket,server=y,address=${PORT},suspend=${SUSPEND}"
`)

	return d.layer.WriteMetadata(expected)
}

// String makes Debug satisfy the Stringer interface.
func (d Debug) String() string {
	return fmt.Sprintf("Debug{ layer: %s, logger: %s }", d.layer, d.logger)
}

type marker struct {
	Debug bool `toml:"debug"`
}

// NewDebug creates a new Debug instance. OK is true if build plan contains "debug" dependency, otherwise false.
func NewDebug(build libjavabuildpack.Build) (Debug, bool) {
	_, ok := build.BuildPlan[DebugDependency]
	if !ok {
		return Debug{}, false
	}

	return Debug{build.Launch.Layer(DebugDependency), build.Logger}, true
}
