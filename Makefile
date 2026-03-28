BINDIR := fabrik-cli

GENESIS := $(BINDIR)/genesis
WALLET  := $(BINDIR)/wallet
NODE 	:= $(BINDIR)/node
FVM     := $(BINDIR)/fvm

.PHONY: all genesis wallet clean

all: genesis wallet node fvm

genesis:
	@mkdir -p $(BINDIR)
	go build -o $(GENESIS) ./cmd/genesis

wallet:
	@mkdir -p $(BINDIR)
	go build -o $(WALLET) ./cmd/wallet

node:
	@mkdir -p $(BINDIR)
	go build -o $(NODE) ./cmd/node

fvm:
	@mkdir -p $(BINDIR)
	go build -o $(FVM) ./cmd/fvm

clean:
	rm -rf $(BINDIR)
	