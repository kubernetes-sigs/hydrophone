# Air Gapped Environments

In this guide we will show how to run Hydrophone in an air-gapped environment with all the required images being pulled from an internal registry. We will also show how to identify the images required to run Hydrophone. For testing, we'll also show how to populate the internal registry with the required images. This guide will assume that a Kubernetes cluster is already running in the air-gapped environment. If you need to create a cluster for testing, we recommend looking at [KIND](https://kind.sigs.k8s.io/docs/user/working-offline/).

## Requirements

Following are the requirements for this guide:

- Docker 18.03 or newer
- Hydrophone installed on the host machine

## Identifying Images

On an air-gapped environment, access to the public Internet is restricted, so Hydrophone can't pull images from public Container registries (`registry.k8s.io`, `gcr.io`, `docker.io`, etc.)
We need to identify the images required to run Hydrophone.

Hydrophone provides a `--list-images` flag. This flag prints a list of images required to run the conformance image.

To print the list of images, run:

```bash
hydrophone --list-images
```

This list contains images required to run the tests inside the conformance image. However, it does not include the conformance image itself or the busybox image that hydrophone uses to pull test results. The images identified will need to be pulled from the public registry and transfered to the internal registry.

## Preparing the Internal Registry

As access to the public registries is restricted, we have to run an internal container registry.
In this guide, we will launch the registry on the same machine using Docker:

```bash
$ docker run -d -p 5001:5000 --restart always --name registry-aigrapped registry:2
```

This registry will be accepting connections on port 5001 on the host IPs.
Once the image is created we will need to populate it with the images required for the conformance tests.

We can do this by pulling the images, re-tagging them, and pushing them to the internal registry.

```bash
$ for image in `hydrophone --list-images && echo registry.k8s.io/conformance:v1.29.0 && echo registry.k8s.io/e2e-test-images/busybox:1.36.1-1`; do
    docker pull $image;
    docker tag $image `echo $image | sed -E 's#^[^/]+/#127.0.0.1:5001/#'`;
    docker push `echo $image | sed -E 's#^[^/]+/#127.0.0.1:5001/#'`;
done
```

We can now verify that the images are pushed to the registry:

```bash
$ curl  http://127.0.0.1:5001/v2/_catalog
{"repositories":["alpine/socat","build-image/distroless-iptables","cloud-provider-gcp/gcp-compute-persistent-disk-csi-driver","e2e-test-images/agnhost","e2e-test-images/apparmor-loader","e2e-test-images/busybox","e2e-test-images/cuda-vector-add","e2e-test-images/httpd","e2e-test-images/ipc-utils","e2e-test-images/jessie-dnsutils", ...]}
```

## Launching Hydrophone in an Air-gapped Environment

For Hydrophone to use the internal registry, we'll need to set up an image list to map the images to the internal registry. This list can be passed to Hydrophone using the `--test-repo-list` flag.

```bash
# make a tempfile for the registry config
export REG_CONFIG=$(mktemp)
cat <<EOF >> $REG_CONFIG
gcAuthenticatedRegistry: 127.0.0.1:5001/authenticated-image-pulling
invalidRegistry: invalid.registry.k8s.io/invalid
privateRegistry: 127.0.0.1:5001/k8s-authenticated-test
microsoftRegistry: 127.0.0.1:5001
dockerLibraryRegistry: 127.0.0.1:5001/library
promoterE2eRegistry: 127.0.0.1:5001/e2e-test-images
buildImageRegistry: 127.0.0.1:5001/build-image
gcEtcdRegistry: 127.0.0.1:5001
gcRegistry: 127.0.0.1:5001
sigStorageRegistry: 127.0.0.1:5001/sig-storage
cloudProviderGcpRegistry: 127.0.0.1:5001/cloud-provider-gcp
EOF
```

Now we can use Hydrophone to run the entire conformance suite with the following command:

```bash
$ hydrophone \
    --conformance \
    --conformance-image localhost:5001/conformance:v1.29.0 \
    --busybox-image localhost:5001/e2e-test-images/busybox:1.36.1-1 \
    --test-repo-list $REG_CONFIG
```

> Note: `--conformance-image` and `--busybox-image` are required to be set to the internal registry as they are not included in the list of images returned by `--list-images`

## Closing Notes

This guide shows how to run Hydrophone in an air-gapped environment. However, it is not fully comprehensive, running in an air-gapped environment might require additional configuration changes, for example using custom settings for DNS, container registry, and your kubernetes cluster.
