
# extensions/helmfile

## Introduction

In order to install or update multiple Helm charts at once, the [helmfile](https://github.com/roboll/helmfile) tool can be used. This extension makes it easy to test _helmfile_ manifests before merging and applying them.

## Usage

### Lint

To check if all to be installed Helm charts pass linting add the following snippet to your `.estafette.yaml` manifest:

```yaml
  lint:
    image: extensions/helmfile:stable
    action: lint
    file: hemlfile.yaml
```

### Diff

To check the diff of the _helmfile_ manifest against an in-pipeline _kind_ container add:

```yaml
  lint:
    services:
    - name: kubernetes
      image: bsycorp/kind:latest-1.17
      readiness:
        path: /kubernetes-ready
        port: 10080

    image: extensions/helmfile:stable
    action: diff
    file: hemlfile.yaml
```

### Diff

To applythe _helmfile_ manifest against an in-pipeline _kind_ container add:

```yaml
  lint:
    services:
    - name: kubernetes
      image: bsycorp/kind:latest-1.17
      readiness:
        path: /kubernetes-ready
        port: 10080

    image: extensions/helmfile:stable
    action: apply
    file: hemlfile.yaml
```
