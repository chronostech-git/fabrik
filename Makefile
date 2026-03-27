BINDIR := fabrik-cli

GENESIS := $(BINDIR)/genesis
WALLET  := $(BINDIR)/wallet
NODE 	:= $(BINDIR)/node

.PHONY: all genesis wallet clean

all: genesis wallet node

genesis:
	@mkdir -p $(BINDIR)
	go build -o $(GENESIS) ./cmd/genesis

wallet:
	@mkdir -p $(BINDIR)
	go build -o $(WALLET) ./cmd/wallet

node:
	@mkdir -p $(BINDIR)
	go build -o $(NODE) ./cmd/node

clean:
	rm -rf $(BINDIR)
	