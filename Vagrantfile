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
  config.vm.box = "ubuntu/bionic64"
  config.vm.network "forwarded_port", guest: 36623, host: 36623, host_ip: "127.0.0.1"
  config.vm.synced_folder ".", "/opt/go/src/github.com/dotStart/Stockpile"

  config.vm.provision "shell", inline: <<-SHELL
    echo "==> Initializing basic build environment"
    sudo apt-get update
    sudo apt-get install -y curl git golang make nodejs npm redis-server

    export GOPATH=/opt/go
    export PATH="$PATH:/opt/go/bin"
    if [ ! -f "$(go env GOPATH)/bin/dep" ]; then \
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
    cp -R $(go env GOPATH)/src/github.com/dotStart/Stockpile/build/linux-amd64 /opt/stockpile
  SHELL
end
