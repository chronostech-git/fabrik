BINDIR := cli

NODE    := $(BINDIR)/node
CHAIN   := $(BINDIR)/chain
FVM     := $(BINDIR)/fvm
ACCOUNT := $(BINDIR)/account
FABNET  := $(BINDIR)/fabnet

.PHONY: all clean deps

all: deps fvm account chain fabnet

# Install missing dependencies
deps:
	@echo "Ensuring Go modules are installed..."
	go mod tidy
	go mod download

chain:
	@mkdir -p $(BINDIR)
	@echo "Building chain..."
	go build -o $(CHAIN) ./cmd/chain

fvm:
	@mkdir -p $(BINDIR)
	@echo "Building fvm..."
	go build -o $(FVM) ./cmd/fvm

account:
	@mkdir -p $(BINDIR)
	@echo "Building account..."
	go build -o $(ACCOUNT) ./cmd/account

fabnet:
	@mkdir -p $(BINDIR)
	@echo "Building fabnet..."
	go build -o $(FABNET) ./cmd/fabnet

clean:
	@echo "Cleaning binaries..."
	rm -rf $(BINDIR)