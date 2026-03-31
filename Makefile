BINDIR := cli

GENESIS := $(BINDIR)/genesis
WALLET  := $(BINDIR)/wallet
NODE 	:= $(BINDIR)/node
CHAIN   := $(BINDIR)/chain
FVM     := $(BINDIR)/fvm
ACCOUNT := $(BINDIR)/account

.PHONY: all genesis wallet clean

all: genesis wallet node fvm account

genesis:
	@mkdir -p $(BINDIR)
	go build -o $(GENESIS) ./cmd/genesis

wallet:
	@mkdir -p $(BINDIR)
	go build -o $(WALLET) ./cmd/wallet

node:
	@mkdir -p $(BINDIR)
	go build -o $(NODE) ./cmd/node

chain:
	@mkdir -p $(BINDIR)
	go build -o $(CHAIN) ./cmd/chain

fvm:
	@mkdir -p $(BINDIR)
	go build -o $(FVM) ./cmd/fvm

account:
	@mkdir -p $(BINDIR)
	go build -o $(ACCOUNT) ./cmd/account

clean:
	rm -rf $(BINDIR)
	