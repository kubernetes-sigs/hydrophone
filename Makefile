# Copyright 2023 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

GOTOOLCHAIN=auto

build:
	go build -o bin/hydrophone main.go

run:
	go run main.go

fmt:
	go fmt ./...

test-unit:
	go test -v ./...

test-race:
	go test --race -v ./...

test: test-unit test-race

verify:
	@hack/verify-all.sh -v

test-e2e: build
	@hack/run-e2e.sh
