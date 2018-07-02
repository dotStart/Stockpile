# -*- mode: ruby -*-
# vi: set ft=ruby :

# --------------------------------------------------- #
# DO NOT DEPLOY THIS IMAGE IN PRODUCTION ENVIRONMENTS #
# --------------------------------------------------- #
# This image is designed specifically for development #
# purposes (specifically to develop and test plugins  #
# in environments where they are not natively         #
# supported (specifically Windows machines)           #
# --------------------------------------------------- #
Vagrant.configure("2") do |config|
  config.vm.box = "hashicorp/precise64"
  config.vm.network "forwarded_port", guest: 36623, host: 36623, host_ip: "127.0.0.1"
  config.vm.synced_folder ".", "/usr/src/go/src/github.com/dotStart/Stockpile"

  config.vm.provision "shell", inline: <<-SHELL
    echo "==> Initializing basic build environment"
    sudo apt-get update
    sudo apt-get install -y curl git make
    if [ ! -d "/usr/local/go" ]; then \
      echo "Fetching go SDK ..."; \
      curl --silent -Lo /usr/src/go.tar.gz https://dl.google.com/go/go1.10.3.linux-amd64.tar.gz; \
      echo "Extracting go ..."; \
      tar -C /usr/local -xzf /usr/src/go.tar.gz; \
    fi
    export PATH=$PATH:/usr/local/go/bin:/usr/src/go/bin
    export GOPATH=/usr/src/go
    mkdir -p /usr/src/go/bin
    if [ ! -f "/usr/src/go/bin/dep" ]; then \
      echo "==> Building go dep"; \
      go get -d -u github.com/golang/dep; \
      cd $(go env GOPATH)/src/github.com/golang/dep; \
      DEP_LATEST=$(git describe --abbrev=0 --tags); \
      git checkout $DEP_LATEST; \
      go install -ldflags="-X main.version=$DEP_LATEST" ./cmd/dep; \
      git checkout master; \
    fi
    echo "==> Building initial binaries"
    cd $(go env GOPATH)/src/github.com/dotStart/Stockpile
    make
    cp -R /usr/src/go/src/github.com/dotStart/Stockpile/build/linux-amd64 /opt/stockpile
  SHELL
end
