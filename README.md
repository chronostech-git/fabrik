# Fabrik (https://chronostech-git.github.io/fabrik)

Fabrik is an experimental blockchain project written in Go. It includes foundational components for wallets, genesis generation, chain persistence, account modeling, and an early virtual machine (FVM).

> **Status:** early-stage / experimental. This repository is not production-ready.

## What this project currently includes

- Wallet generation and key storage (`cmd/wallet`)
- Genesis file creation (`cmd/genesis`)
- A chain runner that loads genesis and writes blocks (`cmd/chain`)
- Basic account creation flow for external and contract account types (`cmd/account`)
- A prototype VM pipeline (parse → compile → execute) via `cmd/fvm`
- Internal packages for blockchain, state, storage, serialization, cryptography, networking scaffolding, and types

## Repository layout

```text
cmd/
  account/   # account CLI
  chain/     # chain execution CLI (loads genesis, appends test block)
  fvm/       # virtual machine runner CLI
  node/      # placeholder CLI (currently no logic)

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

- `github.com/alexflint/go-arg`
- `github.com/holiman/uint256`
- `github.com/syndtr/goleveldb`
- `golang.org/x/crypto`

## Build

Build all CLIs:

```bash
make
```

This creates binaries in `./cli`:

- `cli/genesis`
- `cli/node`
- `cli/fvm`
- `cli/account`

Other available make targets:

```bash
make node
make chain
make fvm
make account
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

Start a node server
```bash
./cli/node --boot <port> 
```

Connect to server as a peer
```bash
./cli/node --peer
```

## Typical local flow

```bash
make
./cli/wallet --datadir ./data
./cli/chain --datadir ./data --new [--dump] [--memory]
./cli/account --datadir ./data --type external
```

## Development notes

- Data is datadir-centric (`keystore`, `genesis`, and chain database paths)
- `scripts/clean.sh` exists in addition to `make clean`
- Project module path: `github.com/chronostech-git/fabrik`

# Please see todos.txt for latest and up-to-date TODO items.
