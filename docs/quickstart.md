# Quickstart Guide

Welcome to Hydrophone! This guide will help you get up and running with Hydrophone in minutes. Hydrophone is a lightweight runner for Kubernetes conformance tests that uses official conformance images from the Kubernetes Release Team.

## Prerequisites

Before you begin, ensure you have:

- **Go 1.19+** installed on your system with `GOPATH` properly set
- Access to a Kubernetes cluster (local or cloud-based)

## Installation

### Option 1: Install via Go (Recommended)

```bash
go install sigs.k8s.io/hydrophone@latest
```

### Option 2: Download from Releases

Visit the [releases page](https://github.com/kubernetes-sigs/hydrophone/releases) and download the latest binary for your platform.

### Option 3: Build from Source

```bash
git clone https://github.com/kubernetes-sigs/hydrophone.git
cd hydrophone
make build
```

## Quick Start with Different Cluster Types

### Local Development with KIND

[KIND](https://kind.sigs.k8s.io/) is perfect for local development and testing.

1. **Install KIND** (if not already installed):

   ```bash
   go install sigs.k8s.io/kind@latest
   ```

2. **Create a cluster**:

   ```bash
   kind create cluster --name hydrophone-test
   ```

3. **Run your first test**:
   ```bash
   hydrophone --focus 'Simple pod should contain last line of the log'
   ```

### Google Kubernetes Engine (GKE)

1. **Create a GKE cluster**:

   ```bash
   gcloud container clusters create hydrophone-cluster \
     --zone us-central1-a \
     --num-nodes 3 \
     --machine-type e2-standard-2
   ```

2. **Get credentials**:

   ```bash
   gcloud container clusters get-credentials hydrophone-cluster --zone us-central1-a
   ```

3. **Run conformance tests**:
   ```bash
   hydrophone --conformance
   ```

### Amazon EKS

1. **Create an EKS cluster** using eksctl or AWS Console
2. **Configure kubectl**:

   ```bash
   aws eks update-kubeconfig --region us-west-2 --name hydrophone-cluster
   ```

3. **Run tests**:
   ```bash
   hydrophone --focus 'sig-storage.*should.*mount'
   ```

### Azure Kubernetes Service (AKS)

1. **Create an AKS cluster**:

   ```bash
   az aks create \
     --resource-group hydrophone-rg \
     --name hydrophone-cluster \
     --node-count 3 \
     --enable-addons monitoring
   ```

2. **Get credentials**:

   ```bash
   az aks get-credentials --resource-group hydrophone-rg --name hydrophone-cluster
   ```

3. **Run tests**:
   ```bash
   hydrophone --conformance
   ```

## Running Tests

### Basic Commands

Hydrophone provides several ways to run tests:

#### Run the Entire Conformance Suite

```bash
hydrophone --conformance
```

This runs all Kubernetes conformance tests and is the most comprehensive validation of your cluster.

#### Run a Specific Test

```bash
hydrophone --focus 'Simple pod should contain last line of the log'
```

The `--focus` flag allows you to run specific tests using regular expressions. This is perfect for:

- Testing specific functionality
- Debugging issues
- Quick validation of changes

#### Run Tests by SIG (Special Interest Group)

```bash
# Run all storage-related tests
hydrophone --focus 'sig-storage.*'

# Run all authentication tests
hydrophone --focus 'sig-auth.*'

# Run all networking tests
hydrophone --focus 'sig-network.*'
```

### Advanced Configuration

#### Custom Conformance Image

```bash
hydrophone --conformance \
  --conformance-image 'registry.k8s.io/conformance:v1.33.2'
```

#### Parallel Test Execution

```bash
hydrophone --conformance --parallel 4
```

#### Verbose Output

```bash
hydrophone --conformance --verbosity 6
```

#### Custom Namespace

```bash
hydrophone --conformance --namespace my-test-namespace
```

#### Skip Specific Tests

```bash
hydrophone --conformance --skip '.*should.*fail'
```

## Example Test Runs

### Example 1: Simple Pod Test

Let's run a basic test to verify pod functionality:

```bash
hydrophone --focus 'Simple pod should contain last line of the log'
```

**Expected Output:**

```
[INFO] Starting test execution...
[INFO] Creating conformance pod in namespace: conformance
[INFO] Waiting for pod to be ready...
[INFO] Test pod is ready, starting test execution...
[INFO] Running test: Simple pod should contain last line of the log
[PASS] Simple pod should contain last line of the log
[INFO] Test completed successfully
[INFO] Cleaning up resources...
```

### Example 2: Storage Test

Test persistent volume functionality:

```bash
hydrophone --focus 'sig-storage.*PersistentVolumes.*should.*mount'
```

**Expected Output:**

```
[INFO] Starting test execution...
[INFO] Creating conformance pod in namespace: conformance
[INFO] Waiting for pod to be ready...
[INFO] Test pod is ready, starting test execution...
[INFO] Running test: [sig-storage] PersistentVolumes-local [Volume type: local] should mount
[PASS] [sig-storage] PersistentVolumes-local [Volume type: local] should mount
[INFO] Test completed successfully
[INFO] Cleaning up resources...
```

### Example 3: Full Conformance Suite

Run the complete conformance test suite:

```bash
hydrophone --conformance --parallel 2 --verbosity 4
```

**Expected Output:**

```
[INFO] Starting conformance test suite...
[INFO] Creating conformance pod in namespace: conformance
[INFO] Waiting for pod to be ready...
[INFO] Test pod is ready, starting conformance suite...
[INFO] Running conformance tests with 2 parallel nodes...
[INFO] Test progress: 50/500 tests completed
[INFO] Test progress: 100/500 tests completed
...
[INFO] Conformance suite completed
[INFO] Results saved to: /tmp/hydrophone-output/
[INFO] Cleaning up resources...
```

## Viewing Test Results

### Real-time Output

Hydrophone provides real-time output during test execution. You'll see:

- Test progress updates
- Pass/fail status for each test
- Error messages and stack traces
- Resource cleanup information

### Output Directory

By default, Hydrophone saves detailed results to `/tmp/hydrophone-output/`. You can customize this with the `--output-dir` flag:

```bash
hydrophone --conformance --output-dir ./my-test-results
```

### Log Files

The output directory contains:

- `conformance.log` - Complete test execution log
- `junit.xml` - JUnit format results (if available)
- `test-results/` - Individual test result files

### Viewing Logs

```bash
# View the main log file
cat /tmp/hydrophone-output/conformance.log

# Follow logs in real-time
tail -f /tmp/hydrophone-output/conformance.log

# Search for specific test results
grep "PASS\|FAIL" /tmp/hydrophone-output/conformance.log
```

## Troubleshooting

### Common Issues

#### 1. Cluster Access Issues

**Problem**: `error loading kubeconfig`
**Solution**: Ensure your `KUBECONFIG` environment variable is set or use `--kubeconfig` flag:

```bash
export KUBECONFIG=/path/to/your/kubeconfig
# or
hydrophone --kubeconfig /path/to/your/kubeconfig --conformance
```

#### 2. Image Pull Issues

**Problem**: `Failed to pull image`
**Solution**: Check if your cluster can access `registry.k8s.io` or use a custom image:

```bash
hydrophone --conformance \
  --conformance-image 'your-registry.com/conformance:v1.33.2'
```

#### 3. Resource Constraints

**Problem**: Tests fail due to insufficient resources
**Solution**: Ensure your cluster has adequate CPU and memory:

```bash
# For KIND clusters, create with more resources
kind create cluster --name hydrophone-test \
  --config - <<EOF
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        system-reserved: memory=512Mi
        kube-reserved: memory=512Mi
- role: worker
  kubeadmConfigPatches:
  - |
    kind: JoinConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        system-reserved: memory=512Mi
        kube-reserved: memory=512Mi
EOF
```

### Getting Help

- **Slack**: Join [#hydrophone], [#sig-testing], or [#k8s-conformance] on [Kubernetes Slack](http://slack.k8s.io/)
- **Issues**: File issues on [GitHub](https://github.com/kubernetes-sigs/hydrophone/issues)
- **Mailing Lists**: [SIG-Testing](https://groups.google.com/forum/#!forum/kubernetes-sig-testing) and [SIG-Release](https://groups.google.com/forum/#!forum/kubernetes-sig-release)

## Next Steps

Now that you have Hydrophone running, you can:

1. **Explore more test patterns** using the `--focus` flag
2. **Customize your test environment** with configuration files
3. **Integrate Hydrophone** into your CI/CD pipelines
4. **Contribute** to the project by filing issues or submitting PRs

For more advanced usage, check out:

- [Air-gapped Environments](air-gapped.md) - Running Hydrophone without internet access
- [Contributing Guide](../CONTRIBUTING.md) - How to contribute to Hydrophone
- [Development Guide](../DEVELOPMENT.md) - Building and developing Hydrophone

Happy testing! 
