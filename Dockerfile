FROM golang:1.16 as builder
# upx 3.95 has issues compressing darwin binaries - https://github.com/upx/upx/issues/301
RUN  apt-get update && apt-get install -y xz-utils && \
  wget -nv -O upx.tar.xz https://github.com/upx/upx/releases/download/v3.96/upx-3.96-amd64_linux.tar.xz; tar xf upx.tar.xz; mv upx-3.96-amd64_linux/upx /usr/bin
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ./ ./
ARG NAME
ARG VERSION
RUN GOOS=linux GOARCH=amd64 make linux compress

FROM ubuntu:20.10
ARG TARGETPLATFORM=amd64
ARG RUNNER_VERSION=2.274.2
ARG DOCKER_VERSION=19.03.12
ENV ACTIONS_RUNNER_VERSION=actions-runner-controller-0.9.0
ENV RUNNER_ASSETS_DIR=/runnertmp

USER root
RUN apt-get update  &&  apt-get install -y  software-properties-common gnupg-agent curl apt-transport-https && \
  add-apt-repository universe && DEBIAN_FRONTEND=noninteractive apt-get install -y  \
  genisoimage  wget jq git sudo npm python-setuptools python3-pip python3-dev build-essential xz-utils ca-certificates supervisor \
  unzip zip software-properties-common sshuttle tzdata  openssh-client rsync shellcheck libunwind8 libyaml-dev libkrb5-3 zlib1g libicu67 liblttng-ust0 && \
  rm -Rf /var/lib/apt/lists/*  && \
  rm -Rf /usr/share/doc && rm -Rf /usr/share/man  && \
  apt-get clean

RUN wget -nv -O go.tar.gz https://golang.org/dl/go1.16.2.linux-amd64.tar.gz && \
  tar -C /usr/local -xzf go.tar.gz  && \
  rm go.tar.gz

ENV  PATH=$PATH:/usr/local/go/bin
# upx 3.95 has issues compressing darwin binaries - https://github.com/upx/upx/issues/301
RUN wget -nv -O upx.tar.xz https://github.com/upx/upx/releases/download/v3.96/upx-3.96-amd64_linux.tar.xz && \
  tar xf upx.tar.xz && \
  mv upx-3.96-amd64_linux/upx /usr/bin && \
  rm -rf upx-3.96-amd64_linux upx.tar.xz

RUN wget -nv --no-check-certificate https://github.com/moshloop/systools/releases/download/3.6/systools.deb && dpkg -i systools.deb
ARG SOPS_VERSION=3.5.0
RUN install_deb https://github.com/mozilla/sops/releases/download/v${SOPS_VERSION}/sops_${SOPS_VERSION}_amd64.deb
RUN install_bin https://github.com/CrunchyData/postgres-operator/releases/download/v4.1.0/expenv
RUN install_bin https://github.com/hongkailiu/gojsontoyaml/releases/download/e8bd32d/gojsontoyaml
RUN install_bin https://github.com/atkrad/wait4x/releases/download/v0.3.0/wait4x-linux-amd64
RUN pip3 install awscli mkdocs mkdocs-material markdown==3.2.1 mkdocs-same-dir mkdocs-autolinks-plugin mkdocs-material-extensions mkdocs-markdownextradata-plugin
RUN wget -nv https://github.com/meterup/github-release/releases/download/v0.7.5/linux-amd64-github-release.bz2 &&  \
  bzip2 -d linux-amd64-github-release.bz2 && \
  chmod +x linux-amd64-github-release && \
  mv linux-amd64-github-release /usr/local/bin/github-release
RUN npm install -g npm@latest
RUN npm install -g pnpm
RUN pnpm install -g netlify-cli gh
RUN go get github.com/mjibson/esc
RUN mv /root/go/bin/esc /usr/local/bin/
RUN curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.10.0/kind-linux-amd64 && \
  chmod +x ./kind && \
  mv ./kind /usr/local/bin/
RUN wget -nv  -O kubectl  https://dl.k8s.io/release/v1.20.0/bin/linux/amd64/kubectl && \
  chmod +x ./kubectl && \
  mv ./kubectl /usr/local/bin
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b /usr/local/bin v1.36.0
ENV LC_ALL=C.UTF-8
ENV LANG=C.UTF-8

RUN wget -nv -O govc.gz https://github.com/vmware/govmomi/releases/download/v0.23.0/govc_linux_amd64.gz && \
  gunzip govc.gz && \
  chmod +x govc && \
  mv govc /usr/local/bin/

RUN export ARCH=$(echo ${TARGETPLATFORM} | cut -d / -f2) \
  && curl -L -o /usr/local/bin/dumb-init https://github.com/Yelp/dumb-init/releases/download/v1.2.2/dumb-init_1.2.2_x86_64 \
  && chmod +x /usr/local/bin/dumb-init



# Docker installation
RUN  adduser --disabled-password --gecos "" --uid 1000 runner \
  && groupadd docker \
  && usermod -aG sudo runner \
  && usermod -aG docker runner \
  && echo "%sudo   ALL=(ALL:ALL) NOPASSWD:ALL" >> /etc/sudoers

RUN  curl -L -o docker.tgz https://download.docker.com/linux/static/stable/x86_64/docker-${DOCKER_VERSION}.tgz && \
  tar --extract \
  --file docker.tgz \
  --strip-components 1 \
  --directory /usr/local/bin/ && \
  rm -rf docker.tgz

RUN mkdir -p "$RUNNER_ASSETS_DIR" \
  && cd "$RUNNER_ASSETS_DIR" \
  && curl -L -o runner.tar.gz https://github.com/actions/runner/releases/download/v${RUNNER_VERSION}/actions-runner-linux-x64-${RUNNER_VERSION}.tar.gz \
  && tar xzf ./runner.tar.gz \
  && rm runner.tar.gz \
  && mv ./externals ./externalstmp

RUN echo AGENT_TOOLSDIRECTORY=/opt/hostedtoolcache > .env \
  && mkdir /opt/hostedtoolcache \
  && chgrp docker /opt/hostedtoolcache \
  && chmod g+rwx /opt/hostedtoolcache

RUN wget -nv -O /entrypoint.sh https://raw.githubusercontent.com/summerwind/actions-runner-controller/${ACTIONS_RUNNER_VERSION}/runner/entrypoint.sh && \
  chmod +x /entrypoint.sh

RUN wget -nv -O /usr/local/bin/modprobe https://raw.githubusercontent.com/summerwind/actions-runner-controller/${ACTIONS_RUNNER_VERSION}/runner/modprobe && \
  chmod +x /usr/local/bin/modprobe

RUN wget -nv -O /etc/supervisor/conf.d/dockerd.conf https://raw.githubusercontent.com/summerwind/actions-runner-controller/${ACTIONS_RUNNER_VERSION}/runner/supervisor/dockerd.conf

RUN mkdir -p /opt/bash-utils/ && wget -nv -O /opt/bash-utils/logger.sh  https://raw.githubusercontent.com/summerwind/actions-runner-controller/${ACTIONS_RUNNER_VERSION}/runner/logger.sh && \
  chmod +x  /opt/bash-utils/logger.sh

RUN wget -nv -O  ${RUNNER_ASSETS_DIR}/RunnerService.js https://raw.githubusercontent.com/summerwind/actions-runner-controller/${ACTIONS_RUNNER_VERSION}/runner/patched/RunnerService.js && \
  chmod +x ${RUNNER_ASSETS_DIR}/RunnerService.js

RUN wget -nv -O  ${RUNNER_ASSETS_DIR}/runsvc.sh  https://raw.githubusercontent.com/summerwind/actions-runner-controller/${ACTIONS_RUNNER_VERSION}/runner/patched/runsvc.sh && \
  chmod +x  ${RUNNER_ASSETS_DIR}/runsvc.sh

COPY startup.sh /startup.sh

# COPY --from=builder /app/.bin/build-tools /bin/
# COPY ./ ./
ENTRYPOINT ["/usr/local/bin/dumb-init", "--"]
CMD ["/startup.sh"]


