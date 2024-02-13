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

- **Replacing Kubernetes Testing Frameworks**: Not intended to replace existing frameworks, but to complement them.
- **Extensive Test Development**: Focus is on running existing tests, not developing new ones.
- **Broad Tool Integration**: Limited integration with third-party tools; maintains simplicity.


## Community

Please reach out for bugs, feature requests, and other issues!
The maintainers of this project are reachable via:

- [Kubernetes Slack] in the [#hydrophone], [#sig-testing] and [#k8s-conformance] channels
- [Filing an issue] against this repo
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
[Kubernetes Slack]: http://slack.k8s.io/C06E3NPR4A3
[SIG-Testing Mailing List]: https://groups.google.com/forum/#!forum/kubernetes-sig-testing
[SIG-Release Mailing List]: https://groups.google.com/forum/#!forum/kubernetes-sig-release
[Conformance suite]: https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/conformance-tests.md
