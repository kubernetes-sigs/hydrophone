# Hydrophone Documentation

Welcome to the Hydrophone documentation! This directory contains comprehensive guides and documentation for using Hydrophone.

## Available Guides

### [Quickstart Guide](quickstart.md) 

**Get started with Hydrophone in minutes!** This comprehensive guide covers:

- Installation instructions (Go install and downloading releases)
- Running against different cluster types (KIND, GKE, EKS, AKS)
- Minimal end-to-end examples with `--focus` and `--conformance` flags
- How to view and interpret test results
- Troubleshooting common issues

### [Air-gapped Environments](air-gapped.md)

Learn how to run Hydrophone in environments without internet access:

- Setting up internal registries
- Identifying required images
- Configuring Hydrophone for offline use

### [CI/CD Integration Guide](ci-cd-integration-guide.md)

Learn how to integrate Hydrophone into your automated workflows:

- GitHub Actions examples
- GitLab CI, Jenkins, and Prow configurations
- Automating conformance testing in CI/CD pipelines

### [Flags and CLI Reference](flags-cli-ref.md)

Complete reference for all Hydrophone command-line flags and options:

- Execution mode flags (`--conformance`, `--focus`, `--skip`)
- Output and configuration options
- YAML configuration file format

### [Understanding E2E Results](understanding-e2e-results.md)

Learn how to read and interpret test results:

- Test output symbols and their meanings
- Understanding test status indicators
- Analyzing failed test results

## Additional Resources

- [Main README](../README.md) - Project overview and basic usage
- [Contributing Guide](../CONTRIBUTING.md) - How to contribute to Hydrophone
- [Development Guide](../DEVELOPMENT.md) - Building and developing Hydrophone
- [Code of Conduct](../code-of-conduct.md) - Community guidelines

## Getting Help

- **Slack**: Join [#hydrophone], [#sig-testing], or [#k8s-conformance] on [Kubernetes Slack](http://slack.k8s.io/)
- **Issues**: File issues on [GitHub](https://github.com/kubernetes-sigs/hydrophone/issues)
- **Mailing Lists**: [SIG-Testing](https://groups.google.com/forum/#!forum/kubernetes-sig-testing) and [SIG-Release](https://groups.google.com/forum/#!forum/kubernetes-sig-release)

---

_Start with the [Quickstart Guide](quickstart.md) if you're new to Hydrophone!_
