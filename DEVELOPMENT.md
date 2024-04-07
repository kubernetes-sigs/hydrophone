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
Usage of bin/hydrophone:
  -busybox-image string
        specify an alternate busybox container image. (default "registry.k8s.io/e2e-test-images/busybox:1.36.1-1")
  -cleanup
        cleanup resources (pods, namespaces etc).
  -conformance
        run conformance tests.
  -conformance-image string
        specify a conformance container image of your choice. (default "registry.k8s.io/conformance:v1.29.0")
  -dry-run
        run in dry run mode.
  -focus string
        focus runs a specific e2e test. e.g. - sig-auth. allows regular expressions.
  -kubeconfig string
        path to the kubeconfig file.
  -list-images
        list all images that will be used during conformance tests.
  -output-dir string
        directory for logs. (defaults to current directory)
  -parallel int
        number of parallel threads in test framework. (default 1)
  -skip string
        skip specific tests. allows regular expressions.
  -test-repo string
        alternate registry for test images
  -test-repo-list string
        yaml file to override registries for test images
  -verbosity int
        verbosity of test framework. (default 4)
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
