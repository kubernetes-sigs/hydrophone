/*
Copyright 2025 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"
)

func buildVersionString() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "hydrophone (version: unknown; build info not available)"
	}

	var (
		modulePath = info.Main.Path
		moduleVer  = info.Main.Version
		goVer      = info.GoVersion

		vcs       string
		revision  string
		vcsTime   string
		modified  *bool
		buildMode string
	)

	for _, s := range info.Settings {
		switch s.Key {
		case "vcs":
			vcs = s.Value
		case "vcs.revision":
			revision = s.Value
		case "vcs.time":
			vcsTime = s.Value
		case "vcs.modified":
			if b, err := strconv.ParseBool(s.Value); err == nil {
				modified = &b
			}
		case "-buildmode":
			buildMode = s.Value
		}
	}

	var b strings.Builder
	fmt.Fprintf(&b, "  module:   %s\n", emptyDash(modulePath))
	fmt.Fprintf(&b, "  version:  %s\n", emptyDash(moduleVer))
	fmt.Fprintf(&b, "  go:       %s\n", emptyDash(goVer))

	if vcs != "" || revision != "" || vcsTime != "" || modified != nil {
		fmt.Fprintf(&b, "  vcs:      %s\n", emptyDash(vcs))
		fmt.Fprintf(&b, "  revision: %s\n", emptyDash(revision))
		fmt.Fprintf(&b, "  time:     %s\n", emptyDash(vcsTime))
		if modified != nil {
			fmt.Fprintf(&b, "  dirty:    %t\n", *modified)
		}
	}

	if buildMode != "" {
		fmt.Fprintf(&b, "  mode:     %s\n", buildMode)
	}

	return b.String()
}

func emptyDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}
