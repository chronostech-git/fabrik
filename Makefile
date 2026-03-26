BINDIR := tools

GENESIS := $(BINDIR)/genesis
WALLET  := $(BINDIR)/wallet

.PHONY: all genesis wallet clean

all: genesis wallet

genesis:
	@mkdir -p $(BINDIR)
	go build -o $(GENESIS) ./cmd/genesis

wallet:
	@mkdir -p $(BINDIR)
	go build -o $(WALLET) ./cmd/wallet

clean:
	rm -rf $(BINDIR)