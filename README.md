# hydrophone

Lightweight runner for kubernetes tests. Uses the conformance image(s) released by
the kubernetes release team to either run individual tests or the entire Conformance suite.
Design is pretty simple, it starts the conformance image as a pod in the `conformance`
namespace, waits for it to finish and then prints out the results.

## Build

```
$ make build
go build -o bin/hydrophone main.go
```

## Install

```
go install github.com/dims/hydrophone@latest
```

## Command line options

```
$ bin/hydrophone --help
Usage of bin/hydrophone:
  -busybox-image string
        specify an alternate busybox container image. (default "registry.k8s.io/e2e-test-images/busybox:1.36.1-1")
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
  -output-dir string
        directory for logs. (defaults to current directory)
  -parallel int
        number of parallel threads in test framework. (default 1)
  -skip string
        skip specific tests. allows regular expressions.
  -verbosity int
        verbosity of test framework. (default 4)
```

## Run

Ensure there is a `KUBECONFIG` environment variable specified or `$HOME/.kube/config` file present before running `hydrophone` Alternatively, you can specify the path to the kubeconfig file with the `--kubeconfig` option.

To run conformance tests use:
```
bin/hydrophone --conformance
```

To run a specific test use:
```
bin/hydrophone --focus 'Simple pod should contain last line of the log'
```

To specify a version of conformance image use:
```
bin/hydrophone --image 'registry.k8s.io/conformance-amd64:v1.29.0'
```

## Troubleshooting

Check if the pod is running:
```
kubectl get pods -n conformance
```

use `kubectl logs` or `kubectl exec` to see what is happening in the pod.

## Cleanup

Delete the pod
```
kubectl delete -n conformance pods/e2e-conformance-test
```

Delete the namespace
```
kubectl delete -n conformance pods/e2e-conformance-test && kubectl delete ns conformance
```