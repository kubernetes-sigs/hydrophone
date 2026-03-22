/*
Copyright 2024 The Kubernetes Authors.

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

package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestParseTestProgress tests the parseTestProgress function's ability to correctly count
// completed tests from log output
func TestParseTestProgress(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedTotal  int
		expectedDone   int
		expectError    bool
		errorSubstring string
	}{
		{
			name: "valid input with completed tests",
			input: `
I0402 12:34:56.789012      12 e2e.go:123] Starting e2e run "01234"
I0402 12:34:57.789012      12 e2e.go:456] Will run 5 of 200 specs

Kubernetes e2e suite
=====================

  SS•SS•SSS•SSSSSSSSSSSSSSS•SSSS
`,
			expectedTotal: 5,
			expectedDone:  4,
			expectError:   false,
		},
		{
			name: "no completed tests yet",
			input: `
I0402 12:34:56.789012      12 e2e.go:123] Starting e2e run "01234"
I0402 12:34:57.789012      12 e2e.go:456] Will run 10 of 200 specs

Kubernetes e2e suite
=====================

  

  Running tests...
`,
			expectedTotal: 10,
			expectedDone:  0,
			expectError:   false,
		},
		{
			name: "all tests completed",
			input: `
I0402 12:34:56.789012      12 e2e.go:123] Starting e2e run "01234"
I0402 12:34:57.789012      12 e2e.go:456] Will run 3 of 7 specs

Kubernetes e2e suite
=====================

  S•S•S•S

  Ran 3 of 7 specs
`,
			expectedTotal: 3,
			expectedDone:  3,
			expectError:   false,
		},
		{
			name: "no test count line",
			input: `
I0402 12:34:56.789012      12 e2e.go:123] Starting e2e run "01234"

Kubernetes e2e suite
=====================

  SS•SS•SSS•SSSSSSSSSSSSSSS•SSS

  Ran 5 of 200 specs
`,
			expectedTotal:  0,
			expectedDone:   0,
			expectError:    true,
			errorSubstring: "could not find test spec count",
		},
		{
			name: "large number of tests",
			input: `
I0402 12:34:56.789012      12 e2e.go:123] Starting e2e run "01234"
I0402 12:34:57.789012      12 e2e.go:456] Will run 1000 of 5000 specs

Kubernetes e2e suite
=====================

  ` + generateDots(500) + `
`,
			expectedTotal: 1000,
			expectedDone:  500,
			expectError:   false,
		},
		{
			name: "mixed success and skip markers",
			input: `
I0402 12:34:56.789012      12 e2e.go:123] Starting e2e run "01234"
I0402 12:34:57.789012      12 e2e.go:456] Will run 15 of 50 specs

Kubernetes e2e suite
=====================

  SSSS•SSSS•SSSS•SSSS•S
`,
			expectedTotal: 15,
			expectedDone:  4,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			total, done, err := parseTestProgress(tt.input)

			if tt.expectError {
				require.Error(t, err)
				if tt.errorSubstring != "" {
					assert.Contains(t, err.Error(), tt.errorSubstring)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedTotal, total, "Total test count does not match")
				assert.Equal(t, tt.expectedDone, done, "Completed test count does not match")
			}
		})
	}
}

// TestParseSpecificTestProgress tests the exact pattern requested in the original issue
func TestParseSpecificTestProgress(t *testing.T) {
	// Test specifically for the format provided in the request
	input := `Will run 5 of 200 specs

SS•SS•SSS•SSSSSSSSSSSSSSS•SSS

`
	total, done, err := parseTestProgress(input)
	require.NoError(t, err)
	assert.Equal(t, 5, total, "Total test count does not match")
	assert.Equal(t, 4, done, "Completed test count (dot characters) does not match")
}

// TestParseTestProgressEdgeCases tests edge cases and boundary conditions
func TestParseTestProgressEdgeCases(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedTotal int
		expectedDone  int
		expectError   bool
	}{
		{
			name: "zero tests to run",
			input: `
I0402 12:34:57.789012      12 e2e.go:456] Will run 0 of 200 specs
`,
			expectedTotal: 0,
			expectedDone:  0,
			expectError:   false,
		},
		{
			name: "maximum boundary test count",
			input: `
I0402 12:34:57.789012      12 e2e.go:456] Will run 10000 of 10000 specs
`,
			expectedTotal: 10000,
			expectedDone:  0,
			expectError:   false,
		},
		{
			name: "malformed spec count",
			input: `
I0402 12:34:57.789012      12 e2e.go:456] Will run abc of 200 specs
`,
			expectedTotal: 0,
			expectedDone:  0,
			expectError:   true,
		},
		{
			name: "dots before spec count line",
			input: `
••••••••••
I0402 12:34:57.789012      12 e2e.go:456] Will run 5 of 200 specs

SS•SS•
`,
			expectedTotal: 5,
			expectedDone:  2,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			total, done, err := parseTestProgress(tt.input)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedTotal, total, "Total test count does not match")
				assert.Equal(t, tt.expectedDone, done, "Completed test count does not match")
			}
		})
	}
}

// generateDots is a helper function to generate a string with n dot characters
func generateDots(n int) string {
	dots := make([]rune, n)
	for i := range dots {
		dots[i] = '•'
	}
	return string(dots)
}

func TestParseTestProgressEmpty(t *testing.T) {
	_, _, err := parseTestProgress("")
	assert.Error(t, err)
}
