export GOLANG_VERSION="1.18.4"

# Install wget
apt update && apt install -y build-essential wget
# Install Go
wget https://golang.org/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz
tar -C /usr/local -xzf go${GOLANG_VERSION}.linux-amd64.tar.gz
rm -f go${GOLANG_VERSION}.linux-amd64.tar.gz

export PATH="$PATH:/usr/local/go/bin"

go build .
