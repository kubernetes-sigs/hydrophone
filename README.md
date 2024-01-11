# hydrophone

![Kubernetes Logo](https://raw.githubusercontent.com/kubernetes-sigs/kubespray/master/docs/img/kubernetes-logo.png)

Lightweight runner for kubernetes tests. Uses the conformance image(s) released by
the kubernetes release team to either run individual tests or the entire [Conformance suite].
Design is pretty simple, it starts the conformance image as a pod in the `conformance`
namespace, waits for it to finish and then prints out the results.

### Project Goals

- **Simplified Kubernetes Testing**: Easy-to-use tool for running Kubernetes conformance tests.
- **Official Conformance Images**: Utilize official conformance images from the Kubernetes Release Team.
- **Flexible Test Execution**: Ability to run individual test, the entire Conformance Test Suite, or anything in between.

### Project Non-Goals

- **Replacing Kubernetes Testing Frameworks**: Not intended to replace existimg frameworks, but to complement them.
- **Extensive Test Development**: Focus is on running existing tests, not developing new ones.
- **Broad Tool Integration**: Limited integration with third-party tools; maintains simplicity.

## Build

```
$ make build
go build -o bin/hydrophone main.go
```

## Install

```
go install sigs.k8s.io/hydrophone@latest
```

## Command line options

```
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

```
bin/hydrophone --conformance
```

To run a specific test use:

```
bin/hydrophone --focus 'Simple pod should contain last line of the log'
```

To specify a version of conformance image use:

```
bin/hydrophone --conformance-image 'registry.k8s.io/conformance:v1.29.0'
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

## Community

Please reach out for bugs, feature requests, and other issues!
The maintainers of this project are reachable via:

- [Kubernetes Slack] in the [#sig-testing] and [#k8s-conformance] channels
- [filing an issue] against this repo
- The Kubernetes [SIG-Testing Mailing List] and [SIG-Release Mailing List]

Current maintainers are [@dims] and [@rjsadow] - feel free to
reach out if you have any questions!

Pull Requests are very welcome!
If you're planning a new feature, please file an issue to discuss first.

Check the [issue tracker] for `help wanted` issues if you're unsure where to
start, or feel free to reach out to discuss. 🙂

See also: our own [contributor guide] and the Kubernetes [community page].

### Code of conduct

Participation in the Kubernetes community is governed by the [Kubernetes Code of Conduct].

<!--links-->

[Kubernetes Code of Conduct]: code-of-conduct.md
[community page]: https://kubernetes.io/community/
[contributor guide]: https://sigs.k8s.io/hydrophone/blob/main/CONTRIBUTING.md
[issue tracker]: https://github.com/kubernetes-sigs/hydrophone/issues
[@dims]: https://github.com/dims
[@rjsadow]: https://github.com/rjsadow
[filing an issue]: https://sigs.k8s.io/hydrophone/issues/new
[Kubernetes Slack]: http://slack.k8s.io/
[SIG-Testing Mailing List]: https://groups.google.com/forum/#!forum/kubernetes-sig-testing
[SIG-Release Mailing List]: https://groups.google.com/forum/#!forum/kubernetes-sig-release
[Conformance suite]: https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/conformance-tests.md
