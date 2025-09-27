# CI/CD Integration Guide for Hydrophone

This guide provides examples for running **Hydrophone** in CI/CD pipelines, including GitHub Actions, GitLab CI, Jenkins, and Prow. It helps users integrate conformance testing into automated workflows.

---

## 1. Introduction

This guide covers GitHub Actions, GitLab CI, Jenkins Pipeline, Prow jobs, tips, and beginner-friendly examples to help you integrate Hydrophone into your CI/CD system.

---

## 2. Beginner-Friendly Examples

These snippets let you run Hydrophone in CI/CD pipelines **without needing a full Kubernetes setup**.  
You can copy–paste them directly.

---

## 3. GitHub Actions

Create `.github/workflows/hydrophone.yml`:

```yaml
name: Hydrophone Quick Test
on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  hydrophone:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: 1.21
      - name: Run Hydrophone (dry run)
        run: |
          go run main.go --dry-run
```

---

## 4. GitLab CI

Add this to `.gitlab-ci.yml`:

```yaml
stages:
  - test

hydrophone:
  stage: test
  image: golang:1.21
  script:
    - go run main.go --dry-run
```

---

## 5. Jenkins Pipeline

In a `Jenkinsfile`:

```groovy
pipeline {
  agent any
  stages {
    stage('Run Hydrophone') {
      steps {
        sh 'go run main.go --dry-run'
      }
    }
  }
}
```

---
## 6. Why these help beginners?

- **Dry Run Mode**: Hydrophone runs without requiring a full Kubernetes cluster. You can test your setup safely.  
- **Official Images**: Uses Ubuntu/Golang images only, no complex dependencies.  
- **Build On Top**: Later you can add kubeconfig, secrets, logs, and artifact uploads easily.  
- **Copy–Paste Ready**: Beginners can directly use these snippets without deep CI/CD knowledge.

---

## 7. GitHub Actions Example

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
          go-version: '1.25'

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

## 8. GitLab CI Example

Example `.gitlab-ci.yml` snippet:

```yaml
stages:
  - conformance

hydrophone-conformance:
  image: golang:1.25.0
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

## 9. Jenkins Pipeline Example

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

## 10. Prow Example

A Prow job snippet (`prow/config.yaml`):

```yaml
periodics:
- name: hydrophone-conformance
  interval: 24h
  agent: kubernetes
  cluster: default
  spec:
    containers:
    - image: golang:1.24.0
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

## 11. Tips and Best Practices

* **Namespaces:** Hydrophone uses the `conformance` namespace by default. Use `--cleanup` if re-running tests.
* **Timeouts:** Adjust `--timeout` depending on cluster size.
* **Artifacts:** Logs (`e2e.log`) and JUnit XML (`junit_01.xml`) can be uploaded for CI/CD reporting.
* **Dry Run:** Use `--dry-run` to quickly verify your setup without executing full conformance tests.

---

## 12. Need Help?

If your CI system isn't listed here, or you run into problems, please open an issue on GitHub or reach out to the community on the [Kubernetes Slack](https://slack.k8s.io/) in the `#sig-testing` or `#hydrophone` channels.

---