#!/usr/bin/env bash

cd /root

# Install go
PKG=go1.15.15.linux-amd64.tar.gz
curl -L -o $PKG "https://golang.org/dl/${PKG}"
rm -rf /usr/local/go && tar -C /usr/local -xzf $PKG
rm $PKG
fgrep 'export PATH=$PATH:/usr/local/go/bin' /root/.bashrc || echo 'export PATH=$PATH:/usr/local/go/bin' >> /root/.bashrc

# Install Delve for debugging
apt-get update
apt-get -y install git
rm -rf delve
git clone https://github.com/go-delve/delve
cd delve
/usr/local/go/bin/go install github.com/go-delve/delve/cmd/dlv
fgrep 'export PATH=$PATH:/root/go/bin' /root/.bashrc || echo 'export PATH=$PATH:/root/go/bin' >> /root/.bashrc


