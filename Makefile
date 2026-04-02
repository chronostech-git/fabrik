BINDIR := cli

NODE 	:= $(BINDIR)/node
CHAIN   := $(BINDIR)/chain
FVM     := $(BINDIR)/fvm
ACCOUNT := $(BINDIR)/account
FABNET  := $(BINDIR)/fabnet

.PHONY: all clean

all: node fvm account chain fabnet

node:
	@mkdir -p $(BINDIR)
	go build -o $(NODE) ./cmd/node
	chmod +x $(NODE)

chain:
	@mkdir -p $(BINDIR)
	go build -o $(CHAIN) ./cmd/chain
	chmod +x $(CHAIN)

fvm:
	@mkdir -p $(BINDIR)
	go build -o $(FVM) ./cmd/fvm
	chmod +x $(FVM)

account:
	@mkdir -p $(BINDIR)
	go build -o $(ACCOUNT) ./cmd/account
	chmod +x $(ACCOUNT)

fabnet:
	@mkdir -p $(BINDIR)
	go build -o $(FABNET) ./cmd/fabnet
	chmod +x $(FABNET)

clean:
	rm -rf $(BINDIR)
	