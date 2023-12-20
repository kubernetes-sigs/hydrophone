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
  -focus string
        focus runs a specific e2e test. e.g. - sig-auth
  -image string
        image let's you select your conformance container image of your choice. for example, for v1.28.0 version of tests, use - 'registry.k8s.io/conformance-amd64:v1.28.0' (default "registry.k8s.io/conformance:v1.28.0")
```

## Run

Ensure there is a `KUBECONFIG` environment variable specified or `$HOME/.kube/config` file present before running `hydrophone`

To run conformance tests use:
```
bin/hydrophone --focus '[Conformance]'
```

To run a specific test use:
```
bin/hydrophone --focus 'Simple pod should support exec'
```

To specify a version of conformance image use:
```
bin/hydrophone --image 'registry.k8s.io/conformance-amd64:v1.28.0'
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