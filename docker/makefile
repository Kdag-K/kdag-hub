KNODE_VERSION=0.1.1
DOGYE_VERSION=0.1.1

all: solidity knode-image dogye-image

solidity:
	wget -O solc https://github.com/ethereum/solidity/releases/download/v0.5.10/solc-static-linux
	chmod +x solc

knode-image:
	go build \
		--ldflags '-extldflags "-static"' \
		-o knode/knode ../cmd/knode/
	cp solc knode/solc
	docker build --no-cache=true -t Kdag-K/knode:$(KNODE_VERSION) -t Kdag-K/knode:latest ./knode/

dogye-image:
	go build \
		--ldflags '-extldflags "-static"' \
		-o dogye/dogye ../cmd/dogye/
	cp solc dogye/solc
	docker build --no-cache=true -t Kdag-K/dogye:$(DOGYE_VERSION) -t Kdag-K/dogye:latest ./dogye/

.PHONY: solidity knode-image dogye-image
