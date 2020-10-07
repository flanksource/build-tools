### Github Workflows

It is recommended that forks of this repo disable github actions if they do not
wish to also publish build artefacts

#### bump.yaml

The bump workflow automatically tags new commits to master with a monotincally
increasing version number

#### binary.yaml

The binary workflow is triggered on a tag push (usually as a result of the bump
workflow).  It generates the build-tools binary and publishes it as a GitHub
release.

#### docker.yaml

The docker workflow is also triggered by a tag push, and generates a docker
image containing the build-tools utility as well as a number of other useful
executables.  The image is published to the [Flanksource
Dockerhub](https://hub.docker.com/r/flanksource/build-tools)

#### dockertest.yaml

The dockertest workflow runs on pull request, building the docker image and
verifying that is is created correctly using
[goss](https://github.com/aelsabbahy/goss)
