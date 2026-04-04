BINDIR := cli

NODE    := $(BINDIR)\node.exe
CHAIN   := $(BINDIR)\chain.exe
FVM     := $(BINDIR)\fvm.exe
ACCOUNT := $(BINDIR)\account.exe
FABNET  := $(BINDIR)\fabnet.exe

.PHONY: all clean deps

all: deps fvm account chain fabnet

deps:
	@echo Ensuring Go modules are installed...
	go mod tidy
	go mod download

chain:
	@echo Building chain...
	go build -o $(CHAIN) .\cmd\chain

fvm:
	@echo Building fvm...
	go build -o $(FVM) .\cmd\fvm

account:
	@echo Building account...
	go build -o $(ACCOUNT) .\cmd\account

fabnet:
	@echo Building fabnet...
	go build -o $(FABNET) .\cmd\fabnet

clean:
	@echo Cleaning binaries...
	if exist $(BINDIR) rmdir /s /q $(BINDIR)