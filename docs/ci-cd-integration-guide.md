# CI/CD Integration Guide for Hydrophone

This guide provides examples for running **Hydrophone** in CI/CD pipelines, including GitHub Actions, GitLab CI, Jenkins, and Prow. It helps users integrate conformance testing into automated workflows.

---

## 1. GitHub Actions Example

Save this workflow as `.github/workflows/hydrophone-smoke.yml` in your repository:

```yaml
name: Hydrophone Smoke Test

on:
  push:
  pull_request:

jobs:
  hydrophone:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'

      - name: Install Hydrophone
        run: go install sigs.k8s.io/hydrophone@latest

      - name: Decode kubeconfig
        run: echo "${KUBE_CONFIG_DATA}" | base64 -d > $HOME/.kube/config
        env:
          KUBE_CONFIG_DATA: ${{ secrets.KUBE_CONFIG_DATA }}

      - name: Run Hydrophone (dry run)
        run: hydrophone --conformance --conformance-image registry.k8s.io/conformance:v1.34.0 --timeout 5m --dry-run
```

**Notes:**

* `$KUBE_CONFIG_DATA` should be stored as a GitHub secret (base64-encoded kubeconfig).
* `--dry-run` ensures a fast test; remove it for a full conformance run.
* Artifacts like `e2e.log` and `junit_01.xml` will be generated during a full run.

---

## 2. GitLab CI Example

Example `.gitlab-ci.yml` snippet:

```yaml
stages:
  - conformance

hydrophone-conformance:
  image: golang:1.22
  stage: conformance
  script:
    - go install sigs.k8s.io/hydrophone@latest
    - echo "$KUBE_CONFIG_DATA" | base64 -d > ~/.kube/config
    - hydrophone --conformance --conformance-image registry.k8s.io/conformance:v1.34.0 --timeout 10m
  only:
    - main
```

**Notes:**

* Set `$KUBE_CONFIG_DATA` as a GitLab CI/CD variable (masked, protected).
* Adjust `timeout` as needed.

---

## 3. Jenkins Pipeline Example

Declarative pipeline snippet:

```groovy
pipeline {
    agent any
    environment {
        KUBE_CONFIG_DATA = credentials('kubeconfig-base64')
    }
    stages {
        stage('Run Hydrophone') {
            steps {
                sh 'go install sigs.k8s.io/hydrophone@latest'
                sh 'echo $KUBE_CONFIG_DATA | base64 -d > $HOME/.kube/config'
                sh 'hydrophone --conformance --conformance-image registry.k8s.io/conformance:v1.34.0 --timeout 10m'
            }
        }
    }
}
```

**Notes:**

* `credentials('kubeconfig-base64')` should be a Jenkins secret containing the base64-encoded kubeconfig.
* You can customize stages, nodes, or parallelization as needed.

---

## 4. Prow Example

A Prow job snippet (`prow/config.yaml`):

```yaml
periodics:
- name: hydrophone-conformance
  interval: 24h
  agent: kubernetes
  cluster: default
  spec:
    containers:
    - image: golang:1.22
      command:
      - /bin/sh
      - -c
      - |
        go install sigs.k8s.io/hydrophone@latest
        echo "$KUBE_CONFIG_DATA" | base64 -d > /root/.kube/config
        hydrophone --conformance --conformance-image registry.k8s.io/conformance:v1.34.0 --timeout 10m
      env:
      - name: KUBE_CONFIG_DATA
        valueFrom:
          secretKeyRef:
            name: kubeconfig-secret
            key: kubeconfig
```

**Notes:**

* Prow jobs may require a Kubernetes cluster and a secret containing the kubeconfig.
* You can schedule periodic runs or trigger on pull requests.

---

## 5. Tips and Best Practices

* **Namespaces:** Hydrophone uses the `conformance` namespace by default. Use `--cleanup` if re-running tests.
* **Timeouts:** Adjust `--timeout` depending on cluster size.
* **Artifacts:** Logs (`e2e.log`) and JUnit XML (`junit_01.xml`) can be uploaded for CI/CD reporting.
* **Dry Run:** Use `--dry-run` to quickly verify your setup without executing full conformance tests.

---

## 6. Need Help?

If your CI system isn't listed here, or you run into problems, please open an issue on GitHub or reach out to the community on the [Kubernetes Slack](https://slack.k8s.io/) in the `#sig-testing` or `#sig-architecture` channels.

---

