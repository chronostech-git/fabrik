BINDIR := cli

NODE 	:= $(BINDIR)/node
CHAIN   := $(BINDIR)/chain
FVM     := $(BINDIR)/fvm
ACCOUNT := $(BINDIR)/account
FABNET  := $(BINDIR)/fabnet

.PHONY: all clean

all: fvm account chain fabnet

chain:
	@mkdir -p $(BINDIR)
	go build -o $(CHAIN) ./cmd/chain

fvm:
	@mkdir -p $(BINDIR)
	go build -o $(FVM) ./cmd/fvm

account:
	@mkdir -p $(BINDIR)
	go build -o $(ACCOUNT) ./cmd/account

fabnet:
	@mkdir -p $(BINDIR)
	go build -o $(FABNET) ./cmd/fabnet

clean:
	rm -rf $(BINDIR)
	