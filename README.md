# hydrophone

### Build

```
$ make build
go build -o bin/hydrophone main.go
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

### Run

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