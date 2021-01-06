FROM golang:1.13.6 as builder
WORKDIR /app
COPY ./ ./
ARG NAME
ARG VERSION
# upx 3.95 has issues compressing darwin binaries - https://github.com/upx/upx/issues/301
RUN  apt-get update && apt-get install -y xz-utils && \
    wget -nv -O upx.tar.xz https://github.com/upx/upx/releases/download/v3.96/upx-3.96-amd64_linux.tar.xz; tar xf upx.tar.xz; mv upx-3.96-amd64_linux/upx /usr/bin
RUN GOOS=linux GOARCH=amd64 make setup linux compress

FROM summerwind/actions-runner-dind:v2.275.1
USER root
COPY --from=builder /app/.bin/build-tools /bin/
ARG SYSTOOLS_VERSION=3.6
COPY ./ ./
RUN apt-get update && \
  apt-get install -y  genisoimage gnupg-agent curl apt-transport-https wget jq git sudo npm python-setuptools python3-pip python3-dev build-essential xz-utils ca-certificates unzip zip software-properties-common && \
  add-apt-repository ppa:longsleep/golang-backports && \
  apt update && \
  apt install -y golang-go && \
  rm -Rf /var/lib/apt/lists/*  && \
  rm -Rf /usr/share/doc && rm -Rf /usr/share/man  && \
  apt-get clean

# upx 3.95 has issues compressing darwin binaries - https://github.com/upx/upx/issues/301
RUN wget -nv -O upx.tar.xz https://github.com/upx/upx/releases/download/v3.96/upx-3.96-amd64_linux.tar.xz && \
   tar xf upx.tar.xz && \
   mv upx-3.96-amd64_linux/upx /usr/bin && \
   rm -rf upx-3.96-amd64_linux upx.tar.xz


RUN wget -nv --no-check-certificate https://github.com/moshloop/systools/releases/download/${SYSTOOLS_VERSION}/systools.deb && dpkg -i systools.deb
ARG SOPS_VERSION=3.5.0
RUN install_deb https://github.com/mozilla/sops/releases/download/v${SOPS_VERSION}/sops_${SOPS_VERSION}_amd64.deb
RUN install_bin https://github.com/CrunchyData/postgres-operator/releases/download/v4.1.0/expenv
RUN install_bin https://github.com/hongkailiu/gojsontoyaml/releases/download/e8bd32d/gojsontoyaml
RUN pip3 install awscli mkdocs mkdocs-material markdown==3.2.1
RUN wget -nv https://github.com/meterup/github-release/releases/download/v0.7.5/linux-amd64-github-release.bz2 &&  \
  bzip2 -d linux-amd64-github-release.bz2 && \
  chmod +x linux-amd64-github-release && \
  mv linux-amd64-github-release /usr/local/bin/github-release
RUN npm install -g npm@latest && hash -r && npm install node --reinstall-packages-from=node
RUN npm install -g netlify-cli gh
RUN go get github.com/mjibson/esc
RUN mv /root/go/bin/esc /usr/local/bin/
RUN curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.9.0/kind-linux-amd64 && \
    chmod +x ./kind && \
    mv ./kind /usr/local/bin/
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b /usr/local/bin v1.31.0
ENV LC_ALL=C.UTF-8
ENV LANG=C.UTF-8

RUN wget -nv -O govc.gz https://github.com/vmware/govmomi/releases/download/v0.23.0/govc_linux_amd64.gz && \
    gunzip govc.gz && \
    chmod +x govc && \
    mv govc /usr/local/bin/
USER runner
# Do not override entrypoint, the one specified in the summerwind image is required


