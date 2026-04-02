# Fabrik (https://chronostech-git.github.io/fabrik)

Fabrik is an experimental blockchain project written in Go. It includes foundational components for wallets, genesis generation, chain persistence, account modeling, and an early virtual machine (FVM).

> **Status:** early-stage / experimental. This repository is not production-ready. 

## Repository layout

```text
cmd/
  account/   # account CLI
  chain/     # chain execution CLI (loads genesis, appends test block)
  fvm/       # virtual machine runner CLI
  fabnet/    # The server for p2p communications
  node/      # very early stages of development 

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
- `cli/node`

Other available make targets:

```bash
make node
make chain
make fvm
make account
make fabnet
make clean
```

## CLI usage

### 1) Run chain

```bash
./cli/chain <datadir> --new [--memory] [--dump]
```

Behavior:
- Initializes chain storage (currently opens LevelDB path `<datadir>/manifest`)
- Creates a chain with the genesis block & coinbase transaction
- Flushes cached blocks to disk
- Pretty-prints chain output when `--dump` is set

***Example output:***

```bash
> cli/chain --datadir data --new --dump

2026/03/31 15:40:48 Coinbase transaction created, signed, and verified:
         hash=0xeec2f312a02fe7534e488f2287ed6cc47cbc58194ebbe133ea89e0d514d777a6
         sig=0x4ba52ca90590dcfb040d478d34b0b1f80d983ac141e511eca258162ca214c97aeeec3be74f422ffaef42447bb6322eb576c66034f9b1cf707d92139d594110d0
         valid=true

Genesis Data
        hash: 0x3e7923130e75257ef709fb8d8c8bffa4050808a7d725dcb90f1f9adedf847b19
        value: 1000000000000000

Current Block Data
        hash: 0xb577e981ac2b9242de2798b863f5897c8cba147185b82fc6ac9831be85362bf7
        time: 1774986048
        txroot: 0x225cc6da0db382fbfa94c37ee083e91e152ab1b0de41adec164f465e41cfdf21
        height: 0

State Balance Data
        Account #1
                addr: 0xca6a7ba3a1a2d4e4f1dc8339e1d8675a2ffef401
                balance: 2000000000000000

Finished.
```

---

### 2) Create an account record

```bash
./cli/account <datadir> --type <external|contract>
```

Behavior:
- Loads wallet key from datadir keystore
- Creates either an external or contract account with that address
- Adds it to in-memory account state

---

### 3) Run FVM (prototype)

```bash
./cli/fvm --file <contract.fab> [--gas 100000] [--debug]
```

Behavior:
- Parses `.fab` input
- Compiles to bytecode
- Executes bytecode on FVM with configured gas limit
- Optional debug output for stack and remaining gas

---

### 4) `node` command status

Start a fabnet node server
```bash
./cli/fabnet --ipaddr <ip> --port <port>
```

Output
```bash
2026/04/02 11:57:27 [FABNET] Server starting... You may now connect using cli/node!
```

Connect to server as a peer
```bash
./cli/node --connect <ip>:<port>
```

Output
```bash
2026/04/02 11:58:03 [NODE] Dialing peer 127.0.0.1:5000... 
```

## Typical local flow

```bash
make
./cli/wallet --datadir ./data
./cli/chain --datadir ./data --new [--dump] [--memory]
./cli/account --datadir ./data --type external
./cli/node --connect <ip>:<port>
./cli/fabnet --ipaddr <ip> --port <port>
```

## Development notes

- Data is datadir-centric (`keystore`, `genesis`, and chain database paths)
- `scripts/clean.sh` exists in addition to `make clean`
- Project module path: `github.com/chronostech-git/fabrik`

# Please see todos.txt for latest and up-to-date TODO items.
