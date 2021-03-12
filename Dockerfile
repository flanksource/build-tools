ARG BASEIMAGE=ubuntu:20.10
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

ARG BASEIMAGE=ubuntu:20.10
FROM $BASEIMAGE
USER root
RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y  \
  genisoimage gnupg-agent curl apt-transport-https wget jq git sudo npm python-setuptools python3-pip python3-dev build-essential xz-utils ca-certificates \
  unzip zip software-properties-common sshuttle podman buildah tzdata runc && \
  rm -Rf /var/lib/apt/lists/*  && \
  rm -Rf /usr/share/doc && rm -Rf /usr/share/man  && \
  apt-get clean

RUN  wget https://golang.org/dl/go1.16.2.linux-amd64.tar.gz && \
  tar -C /usr/local -xzf go1.16.2.linux-amd64.tar.gz && \
  rm go1.16.2.linux-amd64.tar.gz

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
RUN alias docker=podman
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b /usr/local/bin v1.36.0
ENV LC_ALL=C.UTF-8
ENV LANG=C.UTF-8

RUN wget -nv -O govc.gz https://github.com/vmware/govmomi/releases/download/v0.23.0/govc_linux_amd64.gz && \
  gunzip govc.gz && \
  chmod +x govc && \
  mv govc /usr/local/bin/
COPY --from=builder /app/.bin/build-tools /bin/
COPY ./ ./
ARG USER=root
USER $USER
# Do not override entrypoint, the one specified in the summerwind image is required


