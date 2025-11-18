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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyDash(t *testing.T) {
	tests := []struct {
		name   string
		in     string
		expect string
	}{
		{
			name:   "empty string becomes dash",
			in:     "",
			expect: "-",
		},
		{
			name:   "non-empty string unchanged",
			in:     "hydrophone",
			expect: "hydrophone",
		},
		{
			name:   "dash string unchanged",
			in:     "hydrophone-",
			expect: "hydrophone-",
		},
		{
			name:   "whitespace string unchanged",
			in:     " ",
			expect: " ",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := emptyDash(tt.in)
			assert.Equal(t, tt.expect, got)
		})
	}
}
