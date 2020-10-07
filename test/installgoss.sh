#!/bin/bash

set -x
mkdir -p ~/bin
export GOSS_DST=~/bin
curl -fsSL https://goss.rocks/install | sh
goss -version

