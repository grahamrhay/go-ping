# -*- mode: ruby -*-
# vi: set ft=ruby :

VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.vm.box = "debian7"
  config.vm.synced_folder ".", "/go/src/github.com/grahamrhay/go-ping"
  config.vm.provision "shell", path: "provision.sh"
  config.vm.network "forwarded_port", guest: 8080, host: 8080
end
