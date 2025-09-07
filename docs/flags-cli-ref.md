# Hydrophone Flags and Options

This document provides a comprehensive guide to all available flags and options for the Hydrophone CLI tool, a lightweight runner for Kubernetes tests.

## Overview

Hydrophone supports multiple execution modes and configuration options to run Kubernetes conformance tests and related operations. All configuration options can be provided via command-line flags or through a YAML configuration file.

## Command Line Flags

### Execution Mode Flags

These flags determine the primary operation mode of Hydrophone. They are mutually exclusive.

#### `--conformance`
- **Type**: Boolean (flag)
- **Default**: `false`
- **Description**: Run conformance tests. This is equivalent to running `--focus '\[Conformance\]'`.
- **Example**: 
  ```bash
  hydrophone --conformance
  ```

#### `--focus`
- **Type**: String
- **Default**: `""`
- **Description**: Focus runs a specific e2e test. Supports regular expressions for pattern matching.
- **Examples**:
  ```bash
  # Run all sig-auth tests
  hydrophone --focus "sig-auth"
  
  # Run specific test pattern
  hydrophone --focus "should.*create.*pod"
  
  # Run conformance tests (equivalent to --conformance)
  hydrophone --focus '\[Conformance\]'
  ```

#### `--cleanup`
- **Type**: Boolean (flag)
- **Default**: `false`
- **Description**: Cleanup resources (pods, namespaces, etc.) left by previous test runs.
- **Example**:
  ```bash
  hydrophone --cleanup
  ```

#### `--list-images`
- **Type**: Boolean (flag)
- **Default**: `false`
- **Description**: List all images that will be used during conformance tests without running the tests. Useful for air-gapped environments to pre-pull required images.
- **Example**:
  ```bash
  hydrophone --list-images
  
  # Output to file for air-gapped environments
  hydrophone --list-images > required-images.txt
  ```

### Configuration Flags

#### `--config`, `-c`
- **Type**: String
- **Default**: `""`
- **Description**: Path to an optional base configuration file in YAML format.
- **Example**:
  ```bash
  hydrophone --config /path/to/config.yaml --conformance
  ```

#### `--kubeconfig`
- **Type**: String
- **Default**: `$KUBECONFIG` environment variable or `~/.kube/config`
- **Description**: Path to the kubeconfig file for cluster authentication.
- **Example**:
  ```bash
  hydrophone --kubeconfig /path/to/kubeconfig --conformance
  ```

#### `--namespace`, `-n`
- **Type**: String
- **Default**: `"conformance"`
- **Description**: The namespace where the conformance pod is created.
- **Example**:
  ```bash
  hydrophone --namespace my-test-namespace --conformance
  ```

#### `--output-dir`, `-o`
- **Type**: String
- **Default**: `"."`
- **Description**: Directory where test logs and results will be stored.
- **Example**:
  ```bash
  hydrophone --output-dir ./test-results --conformance
  ```

### Test Execution Flags

#### `--parallel`, `-p`
- **Type**: Integer
- **Default**: `1`
- **Description**: Number of parallel threads in test framework. Automatically sets the `--nodes` Ginkgo flag.
- **Example**:
  ```bash
  hydrophone --parallel 4 --conformance
  ```

#### `--verbosity`, `-v`
- **Type**: Integer
- **Default**: `4`
- **Description**: Verbosity level of test framework. Values >= 6 automatically set the `-v` Ginkgo flag.
- **Example**:
  ```bash
  hydrophone --verbosity 6 --conformance
  ```

#### `--skip`
- **Type**: String
- **Default**: `""`
- **Description**: Skip specific tests. Supports regular expressions for pattern matching.
- **Example**:
  ```bash
  # Skip all networking tests
  hydrophone --skip "Networking" --conformance
  
  # Skip multiple test patterns
  hydrophone --skip "Networking|Storage" --conformance
  ```

#### `--dry-run`
- **Type**: Boolean (flag)
- **Default**: `false`
- **Description**: Run in dry-run mode without executing actual tests.
- **Example**:
  ```bash
  hydrophone --dry-run --conformance
  ```

### Container Image Flags

#### `--conformance-image`
- **Type**: String
- **Default**: Auto-detected based on cluster version (e.g., `registry.k8s.io/conformance:v1.28.0`)
- **Description**: Specify a conformance container image of your choice.
- **Example**:
  ```bash
  hydrophone --conformance-image registry.k8s.io/conformance:v1.29.0 --conformance
  ```

#### `--busybox-image`
- **Type**: String
- **Default**: `"registry.k8s.io/e2e-test-images/busybox:1.36.1-1"`
- **Description**: Specify an alternate busybox container image used for log extraction.
- **Example**:
  ```bash
  hydrophone --busybox-image my-registry.com/busybox:latest --conformance
  ```

#### `--test-repo`
- **Type**: String
- **Default**: `""`
- **Description**: Registry for pulling Kubernetes test images.
- **Example**:
  ```bash
  hydrophone --test-repo my-private-registry.com/k8s-test --conformance
  ```

#### `--test-repo-list`
- **Type**: String
- **Default**: `""`
- **Description**: YAML file to override registries for test images.
- **Example**:
  ```bash
  hydrophone --test-repo-list /path/to/repo-list.yaml --conformance
  ```

### Advanced Execution Flags

#### `--continue`
- **Type**: Boolean (flag)
- **Default**: `false`
- **Description**: Connect to an already running conformance test pod instead of starting a new one.
- **Example**:
  ```bash
  hydrophone --continue
  ```

#### `--skip-preflight`
- **Type**: String
- **Default**: `""`
- **Description**: Skip namespace check and use the specified namespace directly.
- **Example**:
  ```bash
  hydrophone --skip-preflight my-namespace --conformance
  ```

#### `--startup-timeout`
- **Type**: Duration
- **Default**: `5m`
- **Description**: Maximum time to wait for the conformance test pod to start up.
- **Example**:
  ```bash
  hydrophone --startup-timeout 10m --conformance
  ```

### Progress Status Flags

#### `--disable-progress-status`
- **Type**: Boolean (flag)
- **Default**: `false`
- **Description**: Disable the periodic progress status updates during test execution.
- **Example**:
  ```bash
  hydrophone --disable-progress-status --conformance
  ```

#### `--progress-status-interval`
- **Type**: Duration
- **Default**: `30s`
- **Description**: Interval duration for progress status updates. Cannot be used with `--disable-progress-status`.
- **Example**:
  ```bash
  hydrophone --progress-status-interval 1m --conformance
  ```

### Extra Arguments Flags

#### `--extra-args`
- **Type**: String slice
- **Default**: `[]`
- **Description**: Additional parameters to be provided to the conformance container. Parameters should be specified as key-value pairs in `--key=value` format.
- **Example**:
  ```bash
  hydrophone --extra-args "--clean-start=true,--allowed-not-ready-nodes=2" --conformance
  ```

#### `--extra-ginkgo-args`
- **Type**: String slice
- **Default**: `[]`
- **Description**: Additional parameters to be provided to Ginkgo runner. Same format as `--extra-args`.
- **Example**:
  ```bash
  hydrophone --extra-ginkgo-args "--timeout=2h,--flake-attempts=3" --conformance
  ```

## Configuration File

Instead of passing all options via command line, you can use a YAML configuration file:

```yaml
# config.yaml
kubeconfig: "/path/to/kubeconfig"
parallel: 4
verbosity: 5
outputDir: "./test-results"
skip: "Networking|Storage"
conformanceImage: "registry.k8s.io/conformance:v1.29.0"
busyboxImage: "registry.k8s.io/e2e-test-images/busybox:1.36.1-1"
namespace: "my-conformance"
dryRun: false
testRepo: "my-registry.com/k8s-test"
extraArgs:
  - "--clean-start=true"
  - "--allowed-not-ready-nodes=2"
extraGinkgoArgs:
  - "--timeout=2h"
  - "--flake-attempts=3"
startupTimeout: "10m"
disableProgressStatus: false
progressStatusInterval: "1m"
```

Use the configuration file:
```bash
hydrophone --config config.yaml --conformance
```

## Common Usage Examples

### Basic conformance test run
```bash
hydrophone --conformance
```

### Run conformance tests with custom output directory
```bash
hydrophone --conformance --output-dir ./my-test-results
```

### Run conformance tests with parallel execution
```bash
hydrophone --conformance --parallel 4 --verbosity 6
```

### Run specific test focus with custom image
```bash
hydrophone --focus "should.*create.*pod" \
  --conformance-image registry.k8s.io/conformance:v1.29.0 \
  --output-dir ./pod-tests
```

### Skip certain tests during conformance run
```bash
hydrophone --conformance \
  --skip "Networking.*LoadBalancer|Storage.*CSI" \
  --parallel 2
```

### Run in air-gapped environment with custom registries
```bash
hydrophone --conformance \
  --conformance-image my-registry.com/conformance:v1.28.0 \
  --busybox-image my-registry.com/busybox:1.36.1-1 \
  --test-repo my-registry.com/k8s-test
```

### Cleanup resources after test run
```bash
hydrophone --cleanup
```

### List images without running tests
```bash
hydrophone --list-images
```

### Connect to existing test pod
```bash
hydrophone --continue
```

### Dry run to validate configuration
```bash
hydrophone --conformance --dry-run --verbosity 6
```

### Use configuration file with command line overrides
```bash
hydrophone --config base-config.yaml \
  --conformance \
  --parallel 8 \
  --output-dir ./override-results
```

## Flag Precedence

When using both configuration files and command line flags:
1. Command line flags take precedence over configuration file values
2. Configuration file values take precedence over default values
3. Environment variables (like `$KUBECONFIG`) are used when neither file nor flags specify values

## Validation Rules

- `--parallel` must be greater than 0
- `--verbosity` must be greater than 0
- `--extra-args` and `--extra-ginkgo-args` must follow `--key=value` format
- `--nodes` or `--procs` cannot be used in `--extra-ginkgo-args` when `--parallel` > 1
- `--progress-status-interval` cannot be used with `--disable-progress-status`
- Execution mode flags (`--conformance`, `--focus`, `--cleanup`, `--list-images`) are mutually exclusive