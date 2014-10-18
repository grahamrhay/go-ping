VERSION=1.3.3
OS=linux
ARCH=amd64

curl -OL https://storage.googleapis.com/golang/go$VERSION.$OS-$ARCH.tar.gz
tar -C /usr/local -xzf go$VERSION.$OS-$ARCH.tar.gz
echo "export PATH=$PATH:/usr/local/go/bin" >> /home/vagrant/.profile
echo "export GOPATH=/go" >> /home/vagrant/.profile
