# hydrophone

This document describes the process for running `hydrophone` on your local machine.

## Getting Started

### Prerequisites

- [Go](https://go.dev/dl/)
- [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
- [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)

### Build

```bash
$ make build
go build -o bin/hydrophone main.go
```

### Install

```bash
$ go install sigs.k8s.io/hydrophone@latest
```

### Command line options

```bash
$ bin/hydrophone --help
Hydrophone is a lightweight runner for Kubernetes tests.

Usage:
  hydrophone [flags]

Flags:
      --busybox-image string        specify an alternate busybox container image. (default "registry.k8s.io/e2e-test-images/busybox:1.36.1-1")
      --cleanup                     cleanup resources (pods, namespaces etc).
  -c, --config string               path to an optional base configuration file.
      --conformance                 run conformance tests.
      --conformance-image string    specify a conformance container image of your choice.
      --dry-run                     run in dry run mode.
      --extra-args strings          Additional parameters to be provided to the conformance container. These parameters should be specified as key-value pairs, separated by commas. Each parameter should start with -- (e.g., --clean-start=true,--allowed-not-ready-nodes=2)
      --extra-ginkgo-args strings   Additional parameters to be provided to Ginkgo runner. This flag has the same format as --extra-args.
      --focus string                focus runs a specific e2e test. e.g. - sig-auth. allows regular expressions.
  -h, --help                        help for hydrophone
      --kubeconfig string           path to the kubeconfig file.
      --list-images                 list all images that will be used during conformance tests.
  -n, --namespace string            the namespace where the conformance pod is created. (default "conformance")
  -o, --output-dir string           directory for logs. (default ".")
  -p, --parallel int                number of parallel threads in test framework (automatically sets the --nodes Ginkgo flag). (default 1)
      --skip string                 skip specific tests. allows regular expressions.
      --startup-timeout duration    max time to wait for the conformance test pod to start up. (default 5m0s)
      --test-repo string            skip specific tests. allows regular expressions.
      --test-repo-list string       yaml file to override registries for test images.
  -v, --verbosity int               verbosity of test framework (values >= 6 automatically sets the -v Ginkgo flag). (default 4)
```

## Run

Ensure there is a `KUBECONFIG` environment variable specified or `$HOME/.kube/config` file present before running `hydrophone` Alternatively, you can specify the path to the kubeconfig file with the `--kubeconfig` option.

To run conformance tests use:

```bash
$ bin/hydrophone --conformance
```

To run a specific test use:

```bash
$ bin/hydrophone --focus 'Simple pod should contain last line of the log'
```

To specify a version of conformance image use:

```bash
$ bin/hydrophone --conformance-image 'registry.k8s.io/conformance:v1.29.0'
```

## Cleanup

Use `hydrophone`'s `--cleanup` flag to remove the tests from your cluster again:

```bash
$ hydrophone --cleanup
```

Note that interrupted tests might leave behind namespaces, check with
`kubectl get namespaces` for those ending with a numeric suffix like `-NNNN` and
delete those as necessary.

## Troubleshooting

Check if the Pod is running:

```bash
$ kubectl --namespace conformance get pods
```

Use `kubectl logs` or `kubectl exec` to see what is happening in the Pod.
