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

package debug

import (
	"fmt"

	"github.com/cloudfoundry/libcfbuildpack/build"
	"github.com/cloudfoundry/libcfbuildpack/layers"
	"github.com/cloudfoundry/libcfbuildpack/logger"
)

// Dependency indicates that a JVM application should be run with debugging enabled.
const Dependency = "debug"

// Debug represents the debug configuration for a JVM application
type Debug struct {
	layer  layers.Layer
	logger logger.Logger
}

// Contribute makes the contribution to launch
func (d Debug) Contribute() error {
	return d.layer.Contribute(marker{true}, func(layer layers.Layer) error {
		return layer.WriteProfile("debug", `PORT=${BPL_DEBUG_PORT:=8080}
SUSPEND=${BPL_DEBUG_SUSPEND:=n}

printf "Debugging enabled on port ${PORT}"

if [[ "${SUSPEND}" = "y" ]]; then
  printf ", suspended on start\n"
else
  printf "\n"
fi

export JAVA_OPTS="${JAVA_OPTS} -agentlib:jdwp=transport=dt_socket,server=y,address=${PORT},suspend=${SUSPEND}"
`)
	}, layers.Cache, layers.Launch)
}

// String makes Debug satisfy the Stringer interface.
func (d Debug) String() string {
	return fmt.Sprintf("Debug{ layer: %s, logger: %s }", d.layer, d.logger)
}

// NewDebug creates a new Debug instance. OK is true if build plan contains "debug" dependency, otherwise false.
func NewDebug(build build.Build) (Debug, bool) {
	_, ok := build.BuildPlan[Dependency]
	if !ok {
		return Debug{}, false
	}

	return Debug{build.Layers.Layer(Dependency), build.Logger}, true
}

type marker struct {
	Debug bool `toml:"debug"`
}

func (m marker) Identity() (string, string) {
	return "Debug", ""
}
