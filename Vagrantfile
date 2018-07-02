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
  config.vm.synced_folder "build/linux64", "/opt/stockpile"
end
