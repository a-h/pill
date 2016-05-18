#!/bin/bash
curl -O https://storage.googleapis.com/golang/go1.6.2.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.6.2.linux-amd64.tar.gz
sudo bash -c 'echo "export PATH=$PATH:/usr/local/go/bin" >> /etc/profile.d/go-tools.sh'
source /etc/profile.d/go-tools.sh

# Setup workspace.
mkdir -p work/src/github.com/a-h/
export GOPATH=$HOME/work
cd $HOME/work/src/github.com/a-h/
git clone https://github.com/a-h/pill
# Get dependencies.
cd $HOME/work/src/github.com/a-h/pill
go get ./...
cd $HOME/work/src/github.com/a-h/pill/httpservice/main
go build

# Run.
#TODO: Get the username setup properly, so we don't use the root mongodb user.
nohup ./main --connectionString mongodb://admin:123456@10.0.1.250:27017 &
