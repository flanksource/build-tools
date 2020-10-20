### Github Workflows

It is recommended that forks of this repo disable github actions if they do not
wish to also publish build artefacts

#### release.yaml

Creates a new release when a tag is pushed to the repository

#### binary.yaml

The binary workflow is triggered on a release.  It generates the build-tools
binary and publishes them to the release.

#### docker.yaml

The docker workflow is also triggered by a release, and generates a docker
image containing the build-tools utility as well as a number of other useful
executables.  The image is published to the [Flanksource
Dockerhub](https://hub.docker.com/r/flanksource/build-tools). The generated
image can be used as a [k8s self-hosted github
runner](https://github.com/summerwind/actions-runner-controller)

#### dockertest.yaml

The dockertest workflow runs on pull request, building the docker image and
verifying that is is created correctly using
[goss](https://github.com/aelsabbahy/goss)
