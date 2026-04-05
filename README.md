# Fabrik (https://chronostech-git.github.io/fabrik)

Fabrik is an experimental blockchain project written in Go. It includes foundational components for wallets, genesis generation, chain persistence, account modeling, and an early virtual machine (FVM).

> **Status:** early-stage / experimental. This repository is not production-ready. 

**note: to view changes, todos, and updates--please visit todos.txt**

## Repository layout

```text
cmd/
  account/   # account CLI
  chain/     # chain execution CLI 
  fvm/       # virtual machine runner CLI
  fabnet/    # The server for p2p communications

internal/
  accounts/
  blockchain/
  crypto/
  fvm/
  p2p/
  serialize/rlp/
  state/
  storage/
  types/

scripts/
  clean.sh
```

## Requirements

- Go **1.25.8** (as declared in `go.mod`)
- `make`

## Dependencies

Primary module dependencies:

github.com/alexflint/go-arg v1.6.1
github.com/holiman/uint256 v1.3.2
github.com/syndtr/goleveldb v1.0.0
golang.org/x/crypto v0.49.0

## Build

Build all CLIs:

```bash
make
```

This creates binaries in `./cli`:

- `cli/account`
- `cli/chain`
- `cli/fabnet`
- `cli/fvm`

Other available make targets:

```bash
make chain
make fvm
make account
make fabnet
make clean
```

## CLI usage

### 1) Run chain - Major update @ 1:04am on 04/05/2026

```bash
.\cli\chain --datadir <datadir> --gaslimit <gaslimit> [--dump]
```

***Example output:***

```bash
> .\cli\chain --datadir data --gaslimit 1000 

2026/04/05 01:39:52 Coinbase transaction signed and verified
        signer: 0x013bbe45ba4506d5886601d47db697fb823edaf2
        sig: 0x4c2f44dcd24eeb991af6f379c02eae0c84c2e9532fa71a08a5efc91e523b5e272a19246d717cb0bb83f26141705e87eddd1934cce1abf57e1008c66459d22ba7
        hash: 0xac7327f8b30d32e223276f853899fa6df7b2b84eea92c8b45e96f7082e937ac7
        valid: true

2026/04/05 01:39:52 Genesis block created at 1775367592 unix time with hash 0x5adb616b34ff5d0c1b96c6a817ea56e8f1ef2e8bbb4b722ccfd28f014aacdad5

```

---

### 2) Create an account record

```bash
./cli/account <datadir> --type <external|contract> [--wallet] --stake 0 --gas 0 [--debug]
```

### 3) Run FVM 

```bash
./cli/fvm --file <contract.fab> [--gas 100000] [--debug]
```

### 4) `fabnet` command

Terminal (A)
  Start FABNET server using the following command
```bash
./cli/fabnet --ipaddr <server-ip> --port <server-port>
```

Output
```bash
2026/04/02 16:27:33 [FABNET] Server started on <server-ip>:<server-port>
```

Terminal (B)
  Test the connection as a simulated peer.
  NOTE: Later, we will have peer discovery that downloads a list of 
        known peers on the network.
```bash
./cli/fabnet --ipaddr <local-ip> --port <different-port> --connect <server-ip>:<server-port>
```

Output
```bash
2026/04/02 16:27:39 [FABNET] Peer connected with address 127.0.0.1:33848
```

You can further testing by writing "PING"
```bash
> PING
```

Server output
```bash
2026/04/02 16:27:41 PING from 127.0.0.1:33848
```

The server should reply with
```bash
> [From <server-ip>:<server-port>] PONG
```

## Development notes

- Data is datadir-centric (`keystore`, `genesis`, and chain database paths)
- `scripts/clean.sh` exists in addition to `make clean`
- Project module path: `github.com/chronostech-git/fabrik`

