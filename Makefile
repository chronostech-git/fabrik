BINDIR := cli

CHAIN   := $(BINDIR)\chain.exe
FVM     := $(BINDIR)\fvm.exe
ACCOUNT := $(BINDIR)\account.exe
FABNET  := $(BINDIR)\fabnet.exe
WALLET  := $(BINDIR)\wallet.exe

.PHONY: all clean deps

all: deps fvm account chain fabnet wallet

deps:
	@echo Ensuring Go modules are installed...
	go mod tidy
	go mod download

chain:
	go build -o $(CHAIN) .\cmd\chain

fvm:
	go build -o $(FVM) .\cmd\fvm

account:
	go build -o $(ACCOUNT) .\cmd\account

fabnet:
	go build -o $(FABNET) .\cmd\fabnet

wallet:
	go build -o $(WALLET) .\cmd\wallet

clean:
	if exist $(BINDIR) rmdir /s /q $(BINDIR)