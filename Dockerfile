FROM ubuntu:bionic

ARG SYSTOOLS_VERSION=3.6

RUN apt-get update && \
  apt-get install -y  genisoimage gnupg-agent curl apt-transport-https wget jq git sudo npm python-setuptools python-pip python-dev build-essential xz-utils upx-ucl ca-certificates unzip zip software-properties-common && \
  rm -Rf /var/lib/apt/lists/*  && \
  rm -Rf /usr/share/doc && rm -Rf /usr/share/man  && \
  apt-get clean


RUN add-apt-repository ppa:longsleep/golang-backports && \
  apt update && \
  apt install -y golang-go

RUN wget -nv --no-check-certificate https://github.com/moshloop/systools/releases/download/${SYSTOOLS_VERSION}/systools.deb && dpkg -i systools.deb
ARG SOPS_VERSION=3.5.0
RUN install_deb https://github.com/mozilla/sops/releases/download/v${SOPS_VERSION}/sops_${SOPS_VERSION}_amd64.deb
RUN install_bin https://github.com/CrunchyData/postgres-operator/releases/download/v4.1.0/expenv
RUN install_bin https://github.com/hongkailiu/gojsontoyaml/releases/download/e8bd32d/gojsontoyaml
RUN pip install awscli mkdocs mkdocs-material
RUN wget -nv https://github.com/meterup/github-release/releases/download/v0.7.5/linux-amd64-github-release.bz2 &&  \
  bzip2 -d linux-amd64-github-release.bz2 && \
  mv linux-amd64-github-release /usr/local/bin
RUN npm install -g netlify-cli now gh
RUN install_certs google.com:443
RUN go get github.com/mjibson/esc
RUN mv /root/go/bin/esc /usr/bin/

